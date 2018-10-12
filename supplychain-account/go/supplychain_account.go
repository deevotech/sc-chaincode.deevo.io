/*
Copyright Deevo LTD. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
*
struct org {
    ObjectType
    Name: string 
    Type: int (1: material supply, 2: farmer, 3: factory, 4: retailer, 5: consumer, 6: tree)
    Id: int
}
*/

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strconv"
    "time"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type account struct {
    ObjectType string `json:"docType"`
    Publickey string `json:"publickey"`
    OrgType string `json:"orgType"`
	Certificate string `json:"certificate"`
	Role int `json:"role"`
}
// ===================================================================================
// Main
// ===================================================================================
func main() {
    err := shim.Start(new(SimpleChaincode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	err := cid.AssertAttributeValue(stub, "supplychain_account.init", "true")
	if err != nil {
		return shim.Error(err.Error())
	}
    return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "initAcc" { //create a new account
        return t.initAcc(stub, args)
    } else if function == "changeRole" { //change name of a specific account
        return t.changeRole(stub, args)
    } else if function == "delete" { //delete a account
        return t.delete(stub, args)
    } else if function == "readAcc" { //read a account
        return t.readAcc(stub, args)
    } else if function == "queryAccsByRole" {
		return t.queryAccsByRole(stub, args)
	} else if function == "queryAccsByRole" {
		return t.queryAccsByRole(stub, args)
	} else if function == "queryAccs" {
		return t.queryAccs(stub, args)
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
func (t *SimpleChaincode) initAcc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	err = cid.AssertAttributeValue(stub, "supplychain_account.initAcc", "true")
	if err != nil {
		return shim.Error(err.Error())
	}
    //   0       1       2     3
    // "1", "adf", "1", "67.0006, -70.5476"
    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
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
    if len(args[3]) <= 0 {
        return shim.Error("4th argument must be a non-empty string")
	}
    role, err := strconv.Atoi(args[3])
    if err != nil {
        return shim.Error("5rd argument must be a numeric string")
    }
    publickey := args[0]
    certificate := args[1]
	orgType := args[2]

    // ==== Check if org already exists ====
    accAsBytes, err := stub.GetState(publickey)
    if err != nil {
        return shim.Error("Failed to get org: " + err.Error())
    } else if accAsBytes != nil {
        fmt.Println("This org already exists: " + publickey)
        return shim.Error("This org already exists: " + publickey)
    }

    // ==== Create org object and marshal to JSON ====
    objectType := "account"
    account := &account{objectType, publickey, certificate, orgType, role}
    accountJSONasBytes, err := json.Marshal(account)
    if err != nil {
        return shim.Error(err.Error())
    }

    // === Save org to state ===
    err = stub.PutState(publickey, accountJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //  ==== Index the org to enable color-based range queries, e.g. return all blue orgs ====
    //  An 'index' is a normal key/value entry in state.
    //  The key is a composite key, with the elements that you want to range query on listed first.
    //  In our case, the composite key is based on indexName~color~name.
    //  This will enable very efficient state range queries based on composite keys matching indexName~color~*
    indexName := "role~publickey"
    rolePublickeyIndexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(account.Role), account.Publickey})
    if err != nil {
        return shim.Error(err.Error())
    }
    //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the org.
    //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
    value := []byte{0x00}
    stub.PutState(rolePublickeyIndexKey, value)

    // ==== org saved and indexed. Return success ====
    fmt.Println("- end init account")
    return shim.Success(nil)
}

// ===============================================
// readorg - read a org from chaincode state
// ===============================================
func (t *SimpleChaincode) readAcc(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var publickey, jsonResp string
	var err error
	err = cid.AssertAttributeValue(stub, "supplychain_account.readAcc", "true")
	if err != nil {
		return shim.Error(err.Error())
	}

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting name of the org to query")
    }

    publickey = args[0]
    valAsbytes, err := stub.GetState(publickey) //get the org from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + publickey + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"Error\":\"org does not exist: " + publickey + "\"}"
        return shim.Error(jsonResp)
    }

    return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a org key/value pair from state
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := cid.AssertAttributeValue(stub, "supplychain_account.delete", "true")
	if err != nil {
		return shim.Error(err.Error())
	}
    var jsonResp string
    var accJSON account
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    publickey := args[0]

    // to maintain the color~name index, we need to read the org first and get its color
    valAsbytes, err := stub.GetState(publickey) //get the org from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " +publickey + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"Error\":\"org does not exist: " + publickey + "\"}"
        return shim.Error(jsonResp)
    }

    err = json.Unmarshal([]byte(valAsbytes), &accJSON)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to decode JSON of: " + publickey + "\"}"
        return shim.Error(jsonResp)
    }

    err = stub.DelState(publickey) //remove the org from chaincode state
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }

    // maintain the index
    indexName := "role~publickey"
    rolePublickeyIndexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(accJSON.Role), accJSON.Publickey})
    if err != nil {
        return shim.Error(err.Error())
    }

    //  Delete index entry to state.
    err = stub.DelState(rolePublickeyIndexKey)
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }
    return shim.Success(nil)
}

// ===========================================================
// transfer a org by setting a new owner name on the org
// ===========================================================
func (t *SimpleChaincode) changeRole(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	err := cid.AssertAttributeValue(stub, "supplychain_account.changeRole", "true")
	if err != nil {
		return shim.Error(err.Error())
	}
    //   0       1
    // "1", "org1"
    if len(args) < 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    publickey:= args[0]
   
	role, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("5rd argument must be a numeric string")
	}
    fmt.Println("- start transferorg ", publickey, role)

    accAsBytes, err := stub.GetState(publickey)
    if err != nil {
        return shim.Error("Failed to get account:" + err.Error())
    } else if accAsBytes == nil {
        return shim.Error("account does not exist")
    }

    accountToTransfer := account{}
    err = json.Unmarshal(accAsBytes, &accountToTransfer) //unmarshal it aka JSON.parse()
    if err != nil {
        return shim.Error(err.Error())
    }
    accountToTransfer.Role = role //change the role

    accJSONasBytes, _ := json.Marshal(accountToTransfer)
    err = stub.PutState(publickey, accJSONasBytes) //rewrite the org
    if err != nil {
        return shim.Error(err.Error())
    }

    fmt.Println("- end transferorg (success)")
    return shim.Success(nil)
}

func (t *SimpleChaincode) queryAccsByRole(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	err := cid.AssertAttributeValue(stub, "supplychain_account.queryAccsByRole", "true")
	if err != nil {
		return shim.Error(err.Error())
	}
    //   0
    // "bob"
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    role, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("5rd argument must be a numeric string")
	}

    queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"account\",\"Role\":\"%d\"}}", role)

    queryResults, err := getQueryResultForQueryString(stub, queryString)
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(queryResults)
}

// ===== Example: Ad hoc rich query ========================================================
// queryorgs uses a query string to perform a query for orgs.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryorgsForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *SimpleChaincode) queryAccs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	err := cid.AssertAttributeValue(stub, "supplychain_account.queryAccs", "true")
	if err != nil {
		return shim.Error(err.Error())
	}
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

func (t *SimpleChaincode) getHistoryForAccs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	err := cid.AssertAttributeValue(stub, "supplychain_account.getHistoryForAccs", "true")
	if err != nil {
		return shim.Error(err.Error())
	}
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    publickey := args[0]

    fmt.Printf("- start getHistoryForOrg: %s\n", publickey)

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