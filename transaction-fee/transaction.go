package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type account struct {
	ObjectType string  `json:"docType"`
	Publickey  string  `json:"publickey"`
	Address    string  `json:"address"`
	Balance    float64 `json:"balance"`
}

// TxFeeChaincode example simple Chaincode implementation
type TxFeeChaincode struct {
}
type MyObject struct {
	ObjectType string `json:"docType"`
	Id         string `json:"id"`
	Value      string `json:"value"`
}

func main() {
	err := shim.Start(new(TxFeeChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *TxFeeChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}
func (t *TxFeeChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "ReadObject" { //change name of a specific account
		return t.readObject(stub, args)
	} else if function == "CreateObject" { //transfer
		return t.createObject(stub, args)
	} else if function == "UpdateObject" { //transfer
		return t.updateObject(stub, args)
	}
	// getHistory AgriProduct, get HistoryProduct
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}
func (t *TxFeeChaincode) readObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var id, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the org to query")
	}

	id = args[0]
	valAsbytes, err := stub.GetState(id) //get the id
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"org does not exist: " + id + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}
func (t *TxFeeChaincode) createObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var id, value, signature, address, fee string
	var err error
	//   0       1
	// "abcd", "1234"
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start init account")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3nd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4nd argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5nd argument must be a non-empty string")
	}

	id = args[0]
	value = args[1]
	signature = args[2]
	address = args[3]
	fee = args[4]
	chainCodeArgs := util.ToChaincodeArgs("getBalance", address)
	response := stub.InvokeChaincode("deevo-account", chainCodeArgs, "deevochannel")
	if response.Status != shim.OK {
		return shim.Error(response.Message)
	}
	accountFromTransfer := account{}
	err = json.Unmarshal(response.Payload, &accountFromTransfer)
	if err != nil {
		return shim.Error(err.Error())
	}
	if accountFromTransfer.Balance < 0 {
		return shim.Error("balance < 0 ")
	}
	fmt.Println("balance: %f", accountFromTransfer.Balance)
	fmt.Println("Signature: " + signature)
	fmt.Println("Signature: " + fee)
	//return shim.Success(nil)

	// ==== Check if org already exists ====
	obAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get object: " + err.Error())
	} else if obAsBytes != nil {
		fmt.Println("This object already exists: " + id)
		return shim.Error("This object already exists: " + id)
	}

	// ==== Create org object and marshal to JSON ====
	objectType := "MyObject"
	ob := &MyObject{objectType, id, value}
	obJSONasBytes, err := json.Marshal(ob)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save org to state ===
	err = stub.PutState(id, obJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "id"
	idIndexKey, err := stub.CreateCompositeKey(indexName, []string{id})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the org.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	valueIndex := []byte{0x00}
	stub.PutState(idIndexKey, valueIndex)

	// ==== org saved and indexed. Return success ====
	fmt.Println("- end init object")
	return shim.Success(nil)
}
func (t *TxFeeChaincode) updateObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//   0       1
	// "abcd", "1234"
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	// ==== Input sanitation ====
	fmt.Println("- start init account")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	id := args[0]
	value := args[1]

	// ==== Check if org already exists ====
	obAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get object: " + err.Error())
	} else if obAsBytes == nil {
		fmt.Println("This object not exists: " + id)
		return shim.Error("This object not exists: " + id)
	}

	ob := MyObject{}
	err = json.Unmarshal(obAsBytes, &ob) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	ob.Value = value
	obJSONasBytes, _ := json.Marshal(ob)
	err = stub.PutState(id, obJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- end update (success)")
	return shim.Success(nil)
}
