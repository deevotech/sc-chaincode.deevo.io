/*
Copyright Deevo LTD. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
*
struct org {
    ObjectType
    Name: string 
    Type: int (1: material supply, 2: farmer, 3: factory, 4: retailer, 5: consumer)
    Id: int
}
*/

package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strconv"
    "strings"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type org {
    ObjectType string `json:"docType"`
    Name string `json:"name"`
    OrgType int `json:"orgType"`
    Id int `json:"id"`
    Location string `json:"location"`
}
type materialsupplier {
    ObjectType string `json:"docType"`
    Name string `json:"name"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
    BatchCode int `json:"batchCode"`
    SupplyTime string `json:"supplyTime"`
}
type farmerMaterial {
    ObjectType string `json:"docType"`
    Name string `json:"name"`
    Qty int `json:"qty"`
    BatchCode int `json:"batchCode"`
    Owner int `json:"owner"`
}
type farmerTree {
    ObjectType string `json:"docType"`
    Name string `json:"name"`
    Qty int `json:"qty"`
    BatchCode int `json:"batchCode"`
    StartTime string `json:"startTime"`
    EndTime string `json:"endTime"`
    LiveTime int `json:"liveTime"`
    Location string `json:"location"`
    Owner int `json:"owner"`
}
type farmerMaterialTree {
    ObjecType string `json:"docType"`
    MaterialBatchCode int `json:"materialBatchCode"`
    TreeBatchCode int `json:"treeBatchCode"`
    Qty int `json:"qty"`
    Timestamp string `json:"timestamp"`
    Owner int `json:"owner"`
}
type farmerProduct {
    Objectype string `json:"docType"`
    Timestamp string `json:"timestamp"`
    Name string `json:"name"`
    TreeBatchCode int `json:"treeBatchCode"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
    FProductBatchCode int `json:"fProductBatchCode"`
}
type factoryProduct {
    Objectype string `json:"docType"`
    Timestamp string `json:"timestamp"`
    Name string `json:"name"`
    FProductBatchCode string `json:"fProductBatchCode"`
    productBatchCode int `json:"productBatchCode"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
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
    return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("invoke is running " + function)

    // Handle different functions
    switch function {
    case "initOrg":
        //create a new org
        return t.initOrg(stub, args)
    case "readOrg":
        //read a org
        return t.readOrg(stub, args)
    case "changeOrg":
        //change name
        return t.changeOrg(stub, args)
    case "queryOrgs":
        //find org based on an ad hoc rich query
        return t.queryOrg(stub, args)
    case "getOrgByRangeId":
        //get org based on range query
        return t.getOrgByRangeId(stub, args)
    case "queryOrgByType":
        return t.queryOrgById(stub, args)
    default:
        //error
        fmt.Println("invoke did not find func: " + function)
        return shim.Error("Received unknown function invocation")
    }
}

// ============================================================
// initOrg - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error

    //  0-name  1-type  2-id
    // "asdf",  "1",  "1"
    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments. Expecting 3")
    }

    // ==== Input sanitation ====
    fmt.Println("- start init marble")
    if len(args[0]) == 0 {
        return shim.Error("1st argument must be a non-empty string")
    }
    if len(args[1]) == 0 {
        return shim.Error("1st argument must be a non-empty string")
    }
    if len(args[2]) == 0 {
        return shim.Error("2nd argument must be a non-empty string")
    }
    if len(args[3]) == 0 {
        return shim.Error("3rd argument must be a non-empty string")
    }
    orgName := args[0]
    orgType, err := strconv.Atoi(args[1])
    if err != nil {
        return shim.Error("3rd argument must be a numeric string")
    }
    orgId, err := strconv.Atoi(args[2])
    if err != nil {
        return shim.Error("3rd argument must be a numeric string")
    }

    // ==== Check if org already exists ====
    orgAsBytes, err := stub.GetPrivateData("collectionOrgs", orgId)
    if err != nil {
        return shim.Error("Failed to get orgs: " + err.Error())
    } else if orgAsBytes != nil {
        fmt.Println("This org already exists: " + orgId)
        return shim.Error("This org already exists: " + orgId)
    }

    // ==== Create org object and marshal to JSON ====
    objectType := "org"
    org := &org{objectType, orgName, orgType, orgId}
    orgJSONasBytes, err := json.Marshal(org)
    if err != nil {
        return shim.Error(err.Error())
    }
    // === Save org to state ===
    err = stub.PutPrivateData("collectionOrgs", orgName, orgJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //  ==== Index the org to enable type and name range queries, e.g. return all blue marbles ====
    //  An 'index' is a normal key/value entry in state.
    //  The key is a composite key, with the elements that you want to range query on listed first.
    //  In our case, the composite key is based on indexName~type~name.
    //  This will enable very efficient state range queries based on composite keys matching indexName~type~*
    indexName := "orgtype~name"
    orgtypeNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{org.OrgType, org.Name})
    if err != nil {
        return shim.Error(err.Error())
    }
    //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the org.
    //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
    value := []byte{0x00}
    stub.PutPrivateData("collectionOrgs", orgtypeNameIndexKey, value)

    // ==== org saved and indexed. Return success ====
    fmt.Println("- end init org")
    return shim.Success(nil)
}

// ===============================================
// readOrg - read a org from chaincode state
// ===============================================
func (t *SimpleChaincode) readOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var name, jsonResp string
    var err error

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting id of the org to query")
    }

    id = args[0]
    valAsbytes, err := stub.GetPrivateData("collectionOrgs", id) //get the marble from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"Error\":\"Org does not exist: " + id + "\"}"
        return shim.Error(jsonResp)
    }

    return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a org key/value pair from state
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var jsonResp string
    var orgJSON org
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }
    orgId, err := strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }

    // to maintain the
    valAsbytes, err := stub.GetPrivateData("collectionOrgs", orgId) //get the org from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + orgId + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"Error\":\"Marble does not exist: " + orgId + "\"}"
        return shim.Error(jsonResp)
    }

    err = json.Unmarshal([]byte(valAsbytes), &orgJSON)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to decode JSON of: " + orgId + "\"}"
        return shim.Error(jsonResp)
    }

    err = stub.DelPrivateData("collectionOrgs", orgId) //remove the marble from chaincode state
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }

    // maintain the index
    indexName := "orgtype~name"
    orgtypeNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{orgJSON.OrgType, orgJSON.Name})
    if err != nil {
        return shim.Error(err.Error())
    }

    //  Delete index entry to state.
    err = stub.DelPrivateData("collectionOrgs", orgtypeNameIndexKey)
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }

    return shim.Success(nil)
}

// ===========================================================
// change a org by setting a new name on the org
// ===========================================================
func (t *SimpleChaincode) changeOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    //   0       1
    // "name", "bob"
    if len(args) < 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    orgId := args[0]
    orgName := strings.ToLower(args[1])
    fmt.Println("- start transferMarble ", orgId, orgName)

    orgAsBytes, err := stub.GetPrivateData("collectionOrgs", orgId)
    if err != nil {
        return shim.Error("Failed to get org:" + err.Error())
    } else if orgAsBytes == nil {
        return shim.Error("org does not exist")
    }

    orgToTransfer := org{}
    err = json.Unmarshal(orgAsBytes, &orgToTransfer) //unmarshal it aka JSON.parse()
    if err != nil {
        return shim.Error(err.Error())
    }
    orgToTransfer.Name = orgName //change the org name

    orgJSONasBytes, _ := json.Marshal(orgToTransfer)
    err = stub.PutPrivateData("collectionOrgs", orgName, orgJSONasBytes) //rewrite the org
    if err != nil {
        return shim.Error(err.Error())
    }

    fmt.Println("- end transferOrg (success)")
    return shim.Success(nil)
}


func (t *SimpleChaincode) getOrgByRangeId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) < 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    startKey := args[0]
    endKey := args[1]

    resultsIterator, err := stub.GetPrivateDataByRange("collectionOrgs", startKey, endKey)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    // buffer is a JSON array containing QueryResults
    var buffer bytes.Buffer
    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
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

    fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) queryOrgById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    orgId := strings.ToLower(args[0])

    queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"org\",\"orgId\":\"%s\"}}", orgId)

    queryResults, err := getQueryResultForQueryString(stub, queryString)
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(queryResults)
}

func (t *SimpleChaincode) queryOrgs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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

    resultsIterator, err := stub.GetPrivateDataQueryResult("collectionOrgs", queryString)
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
