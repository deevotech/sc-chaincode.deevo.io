/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &AccountChaincode{}
}

// AccountChaincode example simple Chaincode implementation
type AccountChaincode struct {
}
type account struct {
	ObjectType string  `json:"docType"`
	Publickey  string  `json:"publickey"`
	Address    string  `json:"address"`
	Balance    float64 `json:"balance"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
}

// Init initializes chaincode
// ===========================
func (t *AccountChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//err := cid.AssertAttributeValue(stub, "account.init", "true")
	/*if err != nil {
		return shim.Error(err.Error())
	}*/
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *AccountChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initAcc" { //create a new accounts
		return t.initAcc(stub, args)
	} else if function == "getBalance" { //change name of a specific account
		return t.getBalance(stub, args)
	} else if function == "transfer" { //transfer
		return t.transfer(stub, args)
	} else if function == "getHistoryForAccs" {
		return t.getHistoryForAccs(stub, args)
	}
	// getHistory AgriProduct, get HistoryProduct
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initorg - create a new org, store into chaincode state
// ============================================================
func (t *AccountChaincode) initAcc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//   0       1       2     3
	// "1", "adf", "1", "67.0006, -70.5476"
	if len(args) != 3 {
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
		return shim.Error("3rd argument must be a non-empty string")
	}
	balance, err := strconv.ParseFloat(args[2], 6)
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	publickey := args[0]
	address := args[1]

	// ==== Check if org already exists ====
	accAsBytes, err := stub.GetState(address)
	if err != nil {
		return shim.Error("Failed to get account: " + err.Error())
	} else if accAsBytes != nil {
		fmt.Println("This account already exists: " + address)
		return shim.Error("This account already exists: " + address)
	}

	// ==== Create org object and marshal to JSON ====
	objectType := "account"
	account := &account{objectType, publickey, address, balance}
	accountJSONasBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save org to state ===
	err = stub.PutState(address, accountJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "address"
	addressIndexKey, err := stub.CreateCompositeKey(indexName, []string{address})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the org.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(addressIndexKey, value)

	// ==== org saved and indexed. Return success ====
	fmt.Println("- end init account")
	return shim.Success(nil)
}

// ===============================================
// get Balance
// ===============================================
func (t *AccountChaincode) getBalance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var address, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the org to query")
	}

	address = args[0]
	valAsbytes, err := stub.GetState(address) //get the address from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + address + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"org does not exist: " + address + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===========================================================
// transfer a value from to
// ===========================================================
func (t *AccountChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 6 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	strValue := args[0]
	hash := args[1]
	r := new(big.Int)
	s := new(big.Int)
	_, err := fmt.Sscan(args[2], r)
	if err != nil {
		return shim.Error("format r")
	}
	_, err = fmt.Sscan(args[3], s)
	if err != nil {
		return shim.Error("format s")
	}
	fromAddress := args[4]
	toAddress := args[5]
	numValue, err := strconv.ParseFloat(args[0], 6)

	if err != nil {
		return shim.Error("5rd argument must be a numeric string")
	}
	fmt.Println("- start transfer ", fromAddress, toAddress, strValue)

	accAsBytesFrom, err := stub.GetState(fromAddress)
	if err != nil {
		return shim.Error("Failed to get from account:" + err.Error())
	} else if accAsBytesFrom == nil {
		return shim.Error("from account does not exist")
	}

	accAsBytesTo, err := stub.GetState(toAddress)
	if err != nil {
		return shim.Error("Failed to get to account:" + err.Error())
	} else if accAsBytesTo == nil {
		return shim.Error("to account does not exist")
	}

	accountFromTransfer := account{}
	err = json.Unmarshal(accAsBytesFrom, &accountFromTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	accountToTransfer := account{}
	err = json.Unmarshal(accAsBytesTo, &accountToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	//publickeyBytes, err := x509.MarshalPKIXPublicKey(&accountFromTransfer.Publickey)
	fmt.Println(accountFromTransfer.Publickey)
	pub, _ := pem.Decode([]byte(accountFromTransfer.Publickey))
	if pub == nil {
		panic("failed to parse PEM block containing the public key")
	}
	prepublickey, err := x509.ParsePKIXPublicKey(pub.Bytes)
	fmt.Println(prepublickey)
	if err != nil {
		return shim.Error("Pre publickey not true")
	}
	publicKey := prepublickey.(*ecdsa.PublicKey)
	fmt.Println("hash: ", hash)
	fmt.Println("r: ", r)
	fmt.Println("s: ", s)

	if !ecdsa.Verify(publicKey, []byte(hash), r, s) {
		return shim.Error("Verity not true")
	}
	if accountFromTransfer.Balance < numValue {
		return shim.Error("balance too small")
	}
	accountFromTransfer.Balance = accountFromTransfer.Balance - numValue
	accountToTransfer.Balance = accountToTransfer.Balance + numValue

	accJSONasBytesFrom, _ := json.Marshal(accountFromTransfer)
	err = stub.PutState(fromAddress, accJSONasBytesFrom)
	if err != nil {
		return shim.Error(err.Error())
	}
	accJSONasBytesTo, _ := json.Marshal(accountToTransfer)
	err = stub.PutState(toAddress, accJSONasBytesTo)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transfer (success)")
	return shim.Success(nil)
}

// ===== Example: Ad hoc rich query ========================================================
// queryorgs uses a query string to perform a query for orgs.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryorgsForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *AccountChaincode) queryAccs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *AccountChaincode) getHistoryForAccs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	publickey := args[0]

	fmt.Printf("- start getHistoryForAcc: %s\n", publickey)

	resultsIterator, err := stub.GetHistoryForKey(publickey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the org
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON org)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAccount returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
