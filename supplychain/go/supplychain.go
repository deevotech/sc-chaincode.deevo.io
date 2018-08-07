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
    "time"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type org struct {
    ObjectType string `json:"docType"`
    Id int `json:"id"`
    Name string `json:"name"`
    OrgType string `json:"orgType"`
    Location string `json:"location"`
}
type supplierMaterial struct {
    ObjectType string `json:"docType"`
    BatchCode int `json:"batchCode"`
    Name string `json:"name"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
}
type farmerMaterial struct {
    ObjectType string `json:"docType"`
    BatchCode int `json:"batchCode"`
    Name string `json:"name"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
}
type farmerTree struct {
    ObjectType string `json:"docType"`
    BatchCode int `json:"batchCode"`
    Name string `json:"name"`
    Qty int `json:"qty"`
    StartTime string `json:"startTime"`
    EndTime string `json:"endTime"`
    LiveTime int `json:"liveTime"`
    Location string `json:"location"`
    Owner int `json:"owner"`
    RateHarvest int `json:"rateharvest"`
}
type farmerMaterialTree struct {
    ObjecType string `json:"docType"`
    MaterialBatchCode int `json:"materialBatchCode"`
    TreeBatchCode int `json:"treeBatchCode"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
}
type agriProduct struct {
    Objectype string `json:"docType"`
    AProductBatchCode int `json:"aProductBatchCode"`
    Timestamp string `json:"timestamp"`
    Name string `json:"name"`
    TreeBatchCode int `json:"treeBatchCode"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
}
type product struct {
    Objectype string `json:"docType"`
    AProductBatchCode string `json:"aProductBatchCode"`
    productBatchCode int `json:"productBatchCode"`
    Name string `json:"name"`
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
    if function == "initOrg" { //create a new org
        return t.initOrg(stub, args)
    } else if function == "changeOrg" { //change name of a specific org
        return t.changeOrg(stub, args)
    } else if function == "delete" { //delete a org
        return t.delete(stub, args)
    } else if function == "readOrg" { //read a org
        return t.readOrg(stub, args)
    } else if function == "queryOrgsByType" { //find orgs for type X using rich query
        return t.queryOrgsByType(stub, args)
    } else if function == "queryOrgs" { //find orgs based on an ad hoc rich query
        return t.queryOrgs(stub, args)
    } else if function == "getHistoryForOrg" { //get history of values for a org
        return t.getHistoryForOrg(stub, args)
    } else if function == "initSupplierMaterial" {
        return t.initSupplierMaterial(stub, args)
    } else if function == "sellMaterial" {
        return t.sellMaterial(stub, args)
    } 
    /*else if function == "initTree" {
        return t.initTree(stub, args)
    } else if function == "materialToTree" {
        return t.materialToTree(stub, args)
    } else if function == "harvestAgriProduct" {
        return t.harvestAgriProduct(stub, args)
    } else if function == "sellAgriProduct" {
        return t.sellAgriProduct(stub, args)
    } else if function == "makeProduct" {
        return t.makeProduct(stub, args)
    } else if function == "sellProduct" {
        return t.sellProduct(stub, args)
    }*/

    fmt.Println("invoke did not find func: " + function) //error
    return shim.Error("Received unknown function invocation")
}
// ============================================================
// initorg - create a new material, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initSupplierMaterial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
    // 0         1    2    3
    // bachcode name qty owner
    // 1        material1  2  1
    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments. Expecting 4")
    }
    // ==== Input sanitation ====
    fmt.Println("- start init material")
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
    batchcode, err := strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    name := strings.ToLower(args[1])
    qty, err := strconv.Atoi(args[2])
    if err != nil {
        return shim.Error("3rd argument must be a numeric string")
    }
    owner, err := strconv.Atoi(args[3])
    if err != nil {
        return shim.Error("4rd argument must be a numeric string")
    }

    // ==== Check if org already exists ====
    supplierMaterialAsBytes, err := stub.GetState(name)
    if err != nil {
        return shim.Error("Failed to get supplier material: " + err.Error())
    } else if supplierMaterialAsBytes != nil {
        fmt.Println("This supplier material already exists: " + name)
        return shim.Error("This supplier material already exists: " + name)
    }

    // ==== Create supplierMaterial object and marshal to JSON ====
    objectType := "supplierMaterial"
    supplierMaterial := &supplierMaterial{objectType, batchcode, name, qty, owner}
    supplierMaterialJSONasBytes, err := json.Marshal(supplierMaterial)
    if err != nil {
        return shim.Error(err.Error())
    }

    // === Save org to state ===
    err = stub.PutState(name, supplierMaterialJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    indexName := "materialBatchcode-name"
    batchcodeNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(supplierMaterial.BatchCode), supplierMaterial.Name})
    if err != nil {
        return shim.Error(err.Error())
    }

    value := []byte{0x00}
    stub.PutState(batchcodeNameIndexKey, value)

    // ==== org saved and indexed. Return success ====
    fmt.Println("- end init supplier material")
    return shim.Success(nil)
}
func (t *SimpleChaincode) sellMaterial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    //   0       1
    // "name", "owner"
    // "material1" "1"
    if len(args) < 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    owner, err:= strconv.Atoi(args[1])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    name := strings.ToLower(args[0])
    fmt.Println("- start transferorg ", owner, name)

    supplierMaterialAsBytes, err := stub.GetState(name)
    if err != nil {
        return shim.Error("Failed to get supplier Material:" + err.Error())
    } else if supplierMaterialAsBytes == nil {
        return shim.Error("supplier material does not exist")
    }

    suppplierMaterialToTransfer := supplierMaterial{}
    err = json.Unmarshal(supplierMaterialAsBytes, &suppplierMaterialToTransfer) //unmarshal it aka JSON.parse()
    if err != nil {
        return shim.Error(err.Error())
    }
    suppplierMaterialToTransfer.Owner = owner //change the name

    supplierMaterialJSONasBytes, _ := json.Marshal(suppplierMaterialToTransfer)
    err = stub.PutState(name, supplierMaterialJSONasBytes) //rewrite the supplier Material
    if err != nil {
        return shim.Error(err.Error())
    }

    fmt.Println("- end transfeMaterial (success)")
    return shim.Success(nil)
}
// ============================================================
// initorg - create a new org, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error

    //   0       1       2     3
    // "1", "adf", "1", "67.0006, -70.5476"
    if len(args) != 4 {
        return shim.Error("Incorrect number of arguments. Expecting 4")
    }

    // ==== Input sanitation ====
    fmt.Println("- start init org")
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
    orgId, err := strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("3rd argument must be a numeric string")
    }
    orgName := args[1]
    orgType := strings.ToLower(args[2])
    orgLocation := strings.ToLower(args[3])

    // ==== Check if org already exists ====
    orgAsBytes, err := stub.GetState(strconv.Itoa(orgId))
    if err != nil {
        return shim.Error("Failed to get org: " + err.Error())
    } else if orgAsBytes != nil {
        fmt.Println("This org already exists: " + strconv.Itoa(orgId))
        return shim.Error("This org already exists: " + strconv.Itoa(orgId))
    }

    // ==== Create org object and marshal to JSON ====
    objectType := "org"
    org := &org{objectType, orgId, orgName, orgType, orgLocation}
    orgJSONasBytes, err := json.Marshal(org)
    if err != nil {
        return shim.Error(err.Error())
    }

    // === Save org to state ===
    err = stub.PutState(strconv.Itoa(orgId), orgJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    //  ==== Index the org to enable color-based range queries, e.g. return all blue orgs ====
    //  An 'index' is a normal key/value entry in state.
    //  The key is a composite key, with the elements that you want to range query on listed first.
    //  In our case, the composite key is based on indexName~color~name.
    //  This will enable very efficient state range queries based on composite keys matching indexName~color~*
    indexName := "orgType~name"
    orgTypeNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{org.OrgType, org.Name})
    if err != nil {
        return shim.Error(err.Error())
    }
    //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the org.
    //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
    value := []byte{0x00}
    stub.PutState(orgTypeNameIndexKey, value)

    // ==== org saved and indexed. Return success ====
    fmt.Println("- end init org")
    return shim.Success(nil)
}

// ===============================================
// readorg - read a org from chaincode state
// ===============================================
func (t *SimpleChaincode) readOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var orgId, jsonResp string
    var err error

    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting name of the org to query")
    }

    orgId = args[0]
    valAsbytes, err := stub.GetState(orgId) //get the org from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + orgId + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"Error\":\"org does not exist: " + orgId + "\"}"
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
    orgId, err:= strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }

    // to maintain the color~name index, we need to read the org first and get its color
    valAsbytes, err := stub.GetState(strconv.Itoa(orgId)) //get the org from chaincode state
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + strconv.Itoa(orgId) + "\"}"
        return shim.Error(jsonResp)
    } else if valAsbytes == nil {
        jsonResp = "{\"Error\":\"org does not exist: " + strconv.Itoa(orgId) + "\"}"
        return shim.Error(jsonResp)
    }

    err = json.Unmarshal([]byte(valAsbytes), &orgJSON)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to decode JSON of: " + strconv.Itoa(orgId) + "\"}"
        return shim.Error(jsonResp)
    }

    err = stub.DelState(strconv.Itoa(orgId)) //remove the org from chaincode state
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }

    // maintain the index
    indexName := "orgType~name"
    orgTypeNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{orgJSON.OrgType, orgJSON.Name})
    if err != nil {
        return shim.Error(err.Error())
    }

    //  Delete index entry to state.
    err = stub.DelState(orgTypeNameIndexKey)
    if err != nil {
        return shim.Error("Failed to delete state:" + err.Error())
    }
    return shim.Success(nil)
}

// ===========================================================
// transfer a org by setting a new owner name on the org
// ===========================================================
func (t *SimpleChaincode) changeOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    //   0       1
    // "1", "org1"
    if len(args) < 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    orgId, err:= strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    newOrgName := strings.ToLower(args[1])
    fmt.Println("- start transferorg ", orgId, newOrgName)

    orgAsBytes, err := stub.GetState(strconv.Itoa(orgId))
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
    orgToTransfer.Name = newOrgName //change the name

    orgJSONasBytes, _ := json.Marshal(orgToTransfer)
    err = stub.PutState(strconv.Itoa(orgId), orgJSONasBytes) //rewrite the org
    if err != nil {
        return shim.Error(err.Error())
    }

    fmt.Println("- end transferorg (success)")
    return shim.Success(nil)
}

func (t *SimpleChaincode) queryOrgsByType(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    //   0
    // "bob"
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    orgType := strings.ToLower(args[0])

    queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"org\",\"orgType\":\"%s\"}}", orgType)

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

func (t *SimpleChaincode) getHistoryForOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    orgId, err:= strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    fmt.Printf("- start getHistoryForOrg: %s\n", strconv.Itoa(orgId))

    resultsIterator, err := stub.GetHistoryForKey(strconv.Itoa(orgId))
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

    fmt.Printf("- getHistoryForOrg returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}
