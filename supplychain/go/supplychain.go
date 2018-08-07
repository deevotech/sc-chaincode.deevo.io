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

type farmerTree struct {
    ObjectType string `json:"docType"`
    TreeId int `json:"treeId"`
    Name string `json:"name"`
    Qty int `json:"qty"`
    StartTime string `json:"startTime"`
    EndTime string `json:"endTime"`
    LiveTime string `json:"liveTime"`
    Location string `json:"location"`
    Owner int `json:"owner"`
    RateHarvest int `json:"rateharvest"`
}

type agriProduct struct {
    Objectype string `json:"docType"`
    AProductBatchCode int `json:"aProductBatchCode"`
    Name string `json:"name"`
    TreeId int `json:"treeId"`
    Qty int `json:"qty"`
    Owner int `json:"owner"`
}
type product struct {
    Objectype string `json:"docType"`
    AProductBatchCode int `json:"aProductBatchCode"`
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
    } else if function == "changeOwnerMaterial" {
        return t.changeOwnerMaterial(stub, args)
    } else if function == "getHistoryForMaterial" {
        return t.getHistoryForMaterial(stub, args)
    } else if function == "queryMaterialsByOwner" {
        return t.queryMaterialsByOwner(stub, args)
    } else if function == "initFarmerTree" {
        return t.initFarmerTree(stub, args)
    } else if function == "harvestAgriProduct" {
        return t.harvestAgriProduct(stub, args)
    } else if function == "changeOwnerAgriProduct" {
        return t.changeOwnerAgriProduct(stub, args)
    } else if function == "makeProduct" {
        return t.makeProduct(stub, args)
    } else if function == "changeOwnerProduct" {
        return t.changeOwnerProduct(stub, args)
    } else if function == "queryAgriProductByOwner" {
        return t.queryAgriProductByOwner(stub, args)
    } else if function == "queryProductByOwner" {
        return t.queryProductByOwner(stub, args)
    } else if function == "getHistoryForAgriProduct" {
        return t.getHistoryForAgriProduct(stub, args)
    } else if function == "getHistoryForProduct" {
        return t.getHistoryForProduct(stub, args)
    }
    // getHistory AgriProduct, get HistoryProduct
    fmt.Println("invoke did not find func: " + function) //error
    return shim.Error("Received unknown function invocation")
}
// ============================================================
// initSupplierMaterial - create a new material, store into chaincode state
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
func (t *SimpleChaincode) changeOwnerMaterial(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

func (t *SimpleChaincode) getHistoryForMaterial(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    name := strings.ToLower(args[0])
    fmt.Printf("- start getHistoryForMaterial: %s\n", name)

    resultsIterator, err := stub.GetHistoryForKey(name)
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

    fmt.Printf("- getHistoryForMaterial returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) queryMaterialsByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    //   0
    // "bob"
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    owner := strings.ToLower(args[0])

    queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"supplierMaterial\",\"owner\":\"%s\"}}", owner)

    queryResults, err := getQueryResultForQueryString(stub, queryString)
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(queryResults)
}

// ============================================================
// initSupplierMaterial - create a new material, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initFarmerTree(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var err error
    /*ObjectType string `json:"docType"`
    TreeId int `json:"treeId"` 0
    Name string `json:"name"` 1
    Qty int `json:"qty"` 2
    StartTime string `json:"startTime"` 3
    EndTime string `json:"endTime"` 4
    LiveTime int `json:"liveTime"` 5
    Location string `json:"location"` 6
    Owner int `json:"owner"` 7
    RateHarvest int `json:"rateharvest"` 8
    */
    if len(args) != 9 {
        return shim.Error("Incorrect number of arguments. Expecting 9")
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
    if len(args[4]) <= 0 {
        return shim.Error("5th argument must be a non-empty string")
    }
    if len(args[5]) <= 0 {
        return shim.Error("6th argument must be a non-empty string")
    }
    if len(args[6]) <= 0 {
        return shim.Error("7th argument must be a non-empty string")
    }
    if len(args[7]) <= 0 {
        return shim.Error("8th argument must be a non-empty string")
    }
    if len(args[8]) <= 0 {
        return shim.Error("9th argument must be a non-empty string")
    }

    treeId, err := strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    name := strings.ToLower(args[1])
    qty, err := strconv.Atoi(args[2])
    if err != nil {
        return shim.Error("3rd argument must be a numeric string")
    }
    starttime := strings.ToLower(args[3])
    endtime := strings.ToLower(args[4])
    livetime := strings.ToLower(args[5])
    location := strings.ToLower(args[6])
    owner, err := strconv.Atoi(args[7])
    if err != nil {
        return shim.Error("8rd argument must be a numeric string")
    }
    harvestrate, err := strconv.Atoi(args[8])
    if err != nil {
        return shim.Error("9rd argument must be a numeric string")
    }

    // ==== Check if farmerTree already exists ====
    farmerTreeAsBytes, err := stub.GetState(name)
    if err != nil {
        return shim.Error("Failed to get tree: " + err.Error())
    } else if farmerTreeAsBytes != nil {
        fmt.Println("This stree already exists: " + name)
        return shim.Error("This tree exists: " + name)
    }

    // ==== Create farmerTreeAs object and marshal to JSON ====
    objectType := "farmerTree"
    farmerTree := &farmerTree{objectType, treeId, name, qty, starttime, endtime, livetime, location, owner, harvestrate}
    farmerTreeJSONasBytes, err := json.Marshal(farmerTree)
    if err != nil {
        return shim.Error(err.Error())
    }

    // === Save org to state ===
    err = stub.PutState(name, farmerTreeJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    indexName := "owner-name"
    ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(farmerTree.Owner), farmerTree.Name})
    if err != nil {
        return shim.Error(err.Error())
    }

    value := []byte{0x00}
    stub.PutState(ownerNameIndexKey, value)
    
    orgAsBytes, err := stub.GetState(strconv.Itoa(treeId))
    if orgAsBytes == nil {
        objectType := "org"
        org := &org{objectType, treeId, name, "6", "67.0006, -70.5476"}
        orgJSONasBytes, err := json.Marshal(org)
        if err != nil {
            return shim.Error(err.Error())
        }
    
        // === Save org to state ===
        err = stub.PutState(strconv.Itoa(treeId), orgJSONasBytes)
        if err != nil {
            return shim.Error(err.Error())
        }
        indexName := "orgType~name"
        orgTypeNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{org.OrgType, org.Name})
        if err != nil {
            return shim.Error(err.Error())
        }

        value := []byte{0x00}
        stub.PutState(orgTypeNameIndexKey, value)

        fmt.Println("- end init org")
    }

    fmt.Println("- end init supplier farmerTree")
    return shim.Success(nil)
}

func (t *SimpleChaincode) harvestAgriProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    /*
    Objectype string `json:"docType"` 
    AProductBatchCode int `json:"aProductBatchCode"` 0
    Name string `json:"name"` 1
    TreeId int `json:"treeId"` 2
    Qty int `json:"qty"` 3
    Owner int `json:"owner"` 4
    */
    var err error
 
    if len(args) != 6 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }
    // ==== Input sanitation ====
    fmt.Println("- start argiProductHarvest")
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
    if len(args[4]) <= 0 {
        return shim.Error("5th argument must be a non-empty string")
    }
    batchcode, err := strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    name := strings.ToLower(args[1])
    treeId, err := strconv.Atoi(args[2])
    if err != nil {
        return shim.Error("3rd argument must be a numeric string")
    }
    qty, err := strconv.Atoi(args[3])
    if err != nil {
        return shim.Error("4rd argument must be a numeric string")
    }
    owner, err := strconv.Atoi(args[4])
    if err != nil {
        return shim.Error("5rd argument must be a numeric string")
    }

    // ==== Check if agriProduct already exists ====
    supplierMaterialAsBytes, err := stub.GetState(name)
    if err != nil {
        return shim.Error("Failed to get supplier material: " + err.Error())
    } else if supplierMaterialAsBytes != nil {
        fmt.Println("This supplier material already exists: " + name)
        return shim.Error("This supplier material already exists: " + name)
    }

    // ==== Create agriproduct object and marshal to JSON ====
    objectType := "agriProduct"
    agriProduct := &agriProduct{objectType, batchcode, name, treeId, qty, owner}
    agriProductJSONasBytes, err := json.Marshal(agriProduct)
    if err != nil {
        return shim.Error(err.Error())
    }

    // === Save org to state ===
    err = stub.PutState(name, agriProductJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    indexName := "owner-name"
    ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(agriProduct.Owner), agriProduct.Name})
    if err != nil {
        return shim.Error(err.Error())
    }

    value := []byte{0x00}
    stub.PutState(ownerNameIndexKey, value)

    // ==== agriProduct saved and indexed. Return success ====
    fmt.Println("- end harvestAgriProduct")
    return shim.Success(nil)
}

func (t *SimpleChaincode) changeOwnerAgriProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    //   0       1
    // "name", "owner"
    // "agriproduct1" "1"
    if len(args) < 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    name := strings.ToLower(args[0])
    owner, err:= strconv.Atoi(args[1])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    fmt.Println("- start transAgriproduct ", owner, name)

    agriProductAsBytes, err := stub.GetState(name)
    if err != nil {
        return shim.Error("Failed to get agri product:" + err.Error())
    } else if agriProductAsBytes == nil {
        return shim.Error("supplier material does not exist")
    }

    agriProductToTransfer := agriProduct{}
    err = json.Unmarshal(agriProductAsBytes, &agriProductToTransfer) //unmarshal it aka JSON.parse()
    if err != nil {
        return shim.Error(err.Error())
    }
    agriProductToTransfer.Owner = owner //change the name

    agriProductJSONasBytes, _ := json.Marshal(agriProductToTransfer)
    err = stub.PutState(name, agriProductJSONasBytes) //rewrite the supplier Material
    if err != nil {
        return shim.Error(err.Error())
    }

    fmt.Println("- end transfeAgriProduct (success)")
    return shim.Success(nil)
}
func (t *SimpleChaincode) makeProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    /*
    Objectype string `json:"docType"`
    AProductBatchCode string `json:"aProductBatchCode"` 0
    productBatchCode int `json:"productBatchCode"` 1
    Name string `json:"name"` 2
    Qty int `json:"qty"` 3
    Owner int `json:"owner"` 4
    */
    var err error
 
    if len(args) != 5 {
        return shim.Error("Incorrect number of arguments. Expecting 5")
    }
    // ==== Input sanitation ====
    fmt.Println("- start make product")
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
    if len(args[4]) <= 0 {
        return shim.Error("5th argument must be a non-empty string")
    }
    aProductBatchCode, err := strconv.Atoi(args[0])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    productBatchCode, err := strconv.Atoi(args[1])
    if err != nil {
        return shim.Error("2rd argument must be a numeric string")
    }
    name := strings.ToLower(args[2])
    qty, err := strconv.Atoi(args[3])
    if err != nil {
        return shim.Error("4rd argument must be a numeric string")
    }
    owner, err := strconv.Atoi(args[4])
    if err != nil {
        return shim.Error("5rd argument must be a numeric string")
    }

    // ==== Check if Product already exists ====
    productAsBytes, err := stub.GetState(name)
    if err != nil {
        return shim.Error("Failed to get supplier material: " + err.Error())
    } else if productAsBytes != nil {
        fmt.Println("This supplier material already exists: " + name)
        return shim.Error("This supplier material already exists: " + name)
    }

    // ==== Create agriproduct object and marshal to JSON ====
    objectType := "product"
    product := &product{objectType, aProductBatchCode, productBatchCode, name, qty, owner}
    productJSONasBytes, err := json.Marshal(product)
    if err != nil {
        return shim.Error(err.Error())
    }

    // === Save product to state ===
    err = stub.PutState(name, productJSONasBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    indexName := "owner-name"
    ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{strconv.Itoa(product.Owner), product.Name})
    if err != nil {
        return shim.Error(err.Error())
    }

    value := []byte{0x00}
    stub.PutState(ownerNameIndexKey, value)

    // ==== product saved and indexed. Return success ====
    fmt.Println("- end make product")
    return shim.Success(nil)
}
 func (t *SimpleChaincode) changeOwnerProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    //   0       1
    // "name", "owner"
    // "product1" "1"
    if len(args) < 2 {
        return shim.Error("Incorrect number of arguments. Expecting 2")
    }

    name := strings.ToLower(args[0])
    owner, err:= strconv.Atoi(args[1])
    if err != nil {
        return shim.Error("1rd argument must be a numeric string")
    }
    fmt.Println("- start transferProduct ", owner, name)

    productAsBytes, err := stub.GetState(name)
    if err != nil {
        return shim.Error("Failed to get product:" + err.Error())
    } else if productAsBytes == nil {
        return shim.Error("supplier material does not exist")
    }

    productToTransfer := agriProduct{}
    err = json.Unmarshal(productAsBytes, &productToTransfer) //unmarshal it aka JSON.parse()
    if err != nil {
        return shim.Error(err.Error())
    }
    productToTransfer.Owner = owner //change the name

    productJSONasBytes, _ := json.Marshal(productToTransfer)
    err = stub.PutState(name, productJSONasBytes) //rewrite the product
    if err != nil {
        return shim.Error(err.Error())
    }

    fmt.Println("- end transferProduct (success)")
    return shim.Success(nil)
}

func (t *SimpleChaincode) queryAgriProductByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    //   0
    // "bob"
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    owner := strings.ToLower(args[0])

    queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"agriProduct\",\"owner\":\"%s\"}}", owner)

    queryResults, err := getQueryResultForQueryString(stub, queryString)
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(queryResults)
}
func (t *SimpleChaincode) queryProductByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    //   0
    // "bob"
    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    owner := strings.ToLower(args[0])

    queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"product\",\"owner\":\"%s\"}}", owner)

    queryResults, err := getQueryResultForQueryString(stub, queryString)
    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(queryResults)
}
func (t *SimpleChaincode) getHistoryForAgriProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    name := strings.ToLower(args[0])
    fmt.Printf("- start getHistoryForAgriProduct: %s\n", name)

    resultsIterator, err := stub.GetHistoryForKey(name)
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

    fmt.Printf("- getHistoryForAgriProduct returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}

func (t *SimpleChaincode) getHistoryForProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) < 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    name := strings.ToLower(args[0])
    fmt.Printf("- start getHistoryForProduct: %s\n", name)

    resultsIterator, err := stub.GetHistoryForKey(name)
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

    fmt.Printf("- getHistoryForProduct returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())
}
