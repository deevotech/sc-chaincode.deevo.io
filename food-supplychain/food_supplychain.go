package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// FoodChaincode demo chaincode
type FoodChaincode struct {
}

func main() {
	err := shim.Start(new(FoodChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *FoodChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *FoodChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initOrgData" {
		return t.initOrgData(stub, args)
	} else if function == "createLog" { //create a new org
		return t.createLog(stub, args)
	} else if function == "updateLog" {
		return t.updateLog(stub, args)
	} else if function == "createAuditor" { //create a new org
		return t.createAuditor(stub, args)
	} else if function == "updateAuditor" {
		return t.updateAuditor(stub, args)
	} else if function == "createAuditAction" { //create a new org
		return t.createAuditAction(stub, args)
	} else if function == "updateAuditAction" {
		return t.updateAuditAction(stub, args)
	} else if function == "getObject" {
		return t.getObject(stub, args)
	} else if function == "getAuditOfObject" {
		return t.getAuditOfObject(stub, args)
	} else if function == "getAuditsOfAuditor" {
		return t.getAuditsOfAuditor(stub, args)
	} else if function == "createTraceable" {
		return t.createTraceable(stub, args)
	} else if function == "updateTraceable" {
		return t.updateTraceable(stub, args)
	} else if function == "getLogsOfSupplychain" {
		return t.getLogsOfSupplychain(stub, args)
	} else if function == "getLogsOfProduct" {
		return t.getLogsOfProduct(stub, args)
	} else if function == "getQueryResultForQueryString" {
		return t.getQueryResultForQueryString(stub, args)
	} else if function == "getHistoryOfObject" {
		return t.getHistoryOfObject(stub, args)
	}
	// getHistory AgriProduct, get HistoryProduct
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// Init all data of an ORGs
// ========================================
func (t *FoodChaincode) initOrgData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start initOrgData", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	newData := InitData{}
	err := json.Unmarshal([]byte(args[0]), &newData)
	if err != nil {
		return shim.Error("Failed to decode json of ORG: " + err.Error())
	}

	for _, traceable := range newData.Traceable {
		traceableAsBytes, err := json.Marshal(traceable)
		if err != nil {
			return shim.Error("Failed to encode json of object " + traceable.ID)
		}
		result := t.createObject(stub, traceableAsBytes, traceable.ID)
		if result.Status != shim.OK {
			return result
		}
	}

	for _, auditor := range newData.Auditors {
		auditorAsBytes, err := json.Marshal(auditor)
		if err != nil {
			return shim.Error("Failed to encode json of Auditor " + auditor.ID)
		}
		if auditor.ObjectType != TYPE_AUDITOR {
			return shim.Error("Expexted objectType " + TYPE_AUDITOR + " for Auditor")
		}
		result := t.createObject(stub, auditorAsBytes, auditor.ID)
		if result.Status != shim.OK {
			return result
		}
	}

	fmt.Println("- end initOrgData (success)")
	return shim.Success(nil)
}

// Methods of Traceable data
func (t *FoodChaincode) createTraceable(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createTraceable", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newTraceable := Traceable{}
	err := json.Unmarshal(jsonBytes, &newTraceable)
	if err != nil {
		return shim.Error("Failed to decode json: " + err.Error())
	}

	result := t.createObject(stub, jsonBytes, newTraceable.ID)

	if result.Status == shim.OK {
		fmt.Println("- end createTraceable (success)")
	}
	return result
}

func (t *FoodChaincode) updateTraceable(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateTraceable", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newTraceable := Traceable{}
	err := json.Unmarshal(jsonBytes, &newTraceable)
	if err != nil {
		return shim.Error("Failed to decode json: " + err.Error())
	}

	result := t.updateObject(stub, jsonBytes, newTraceable.ID)

	if result.Status == shim.OK {
		fmt.Println("- end updateTraceable (success)")
	}
	return result
}

// Methods on Log
// ========================================
func (t *FoodChaincode) createLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createLog", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newLog := Log{}
	err := json.Unmarshal(jsonBytes, &newLog)
	if err != nil {
		return shim.Error("Failed to decode json of Log: " + err.Error())
	}
	if newLog.ObjectType != TYPE_LOG {
		return shim.Error("Expexted objectType " + TYPE_LOG + " for Log")
	}

	result := t.createObject(stub, jsonBytes, newLog.ID)

	if len(newLog.Supplychain) > 0 {
		result = t.putCompositeKey(stub, CK_SC_LOG, []string{newLog.Supplychain, newLog.ID})
		if result.Status != shim.OK {
			fmt.Println("- end createLog (failed)")
			return result
		}
	}

	if len(newLog.Product) > 0 {
		result = t.putCompositeKey(stub, CK_PRODUCT_LOG, []string{newLog.Product, newLog.ID})
		if result.Status != shim.OK {
			fmt.Println("- end createLog (failed)")
			return result
		}
	}

	if result.Status == shim.OK {
		fmt.Println("- end createLog (success)")
	}
	return result
}

func (t *FoodChaincode) updateLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateLog", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newLog := Log{}
	err := json.Unmarshal(jsonBytes, &newLog)
	if err != nil {
		return shim.Error("Failed to decode json of Log: " + err.Error())
	}
	if newLog.ObjectType != TYPE_LOG {
		return shim.Error("Expexted objectType " + TYPE_LOG + " for Log")
	}

	result := t.updateLogHandler(stub, jsonBytes, newLog)

	if result.Status == shim.OK {
		fmt.Println("- end updateLog (success)")
	}
	return shim.Success(nil)
}

// Methods on Auditor
// ========================================
func (t *FoodChaincode) createAuditor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createAuditor", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newAuditor := Auditor{}
	err := json.Unmarshal(jsonBytes, &newAuditor)
	if err != nil {
		return shim.Error("Failed to decode json of Auditor: " + err.Error())
	}
	if newAuditor.ObjectType != TYPE_AUDITOR {
		return shim.Error("Expexted objectType " + TYPE_AUDITOR + " for Auditor")
	}

	result := t.createObject(stub, jsonBytes, newAuditor.ID)

	if result.Status == shim.OK {
		fmt.Println("- end createAuditor (success)")
	}
	return result
}

func (t *FoodChaincode) updateAuditor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateAuditor", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newAuditor := Auditor{}
	err := json.Unmarshal(jsonBytes, &newAuditor)
	if err != nil {
		return shim.Error("Failed to decode json of Auditor: " + err.Error())
	}
	if newAuditor.ObjectType != TYPE_AUDITOR {
		return shim.Error("Expexted objectType " + TYPE_AUDITOR + " for Auditor")
	}

	result := t.updateObject(stub, jsonBytes, newAuditor.ID)

	if result.Status == shim.OK {
		fmt.Println("- end updateAuditor (success)")
	}
	return shim.Success(nil)
}

// Methods on AuditActions
// ========================================
func (t *FoodChaincode) createAuditAction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createAuditAction", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newAuditAction := AuditAction{}
	err := json.Unmarshal(jsonBytes, &newAuditAction)
	if err != nil {
		return shim.Error("Failed to decode json of AuditAction: " + err.Error())
	}

	if len(newAuditAction.ObjectID) < 1 {
		return shim.Error("ObjectID can not by empty")
	}
	if len(newAuditAction.Auditor) < 1 {
		return shim.Error("AuditorID can not by empty")
	}
	if newAuditAction.ObjectType != TYPE_AUDITACTION {
		return shim.Error("Expexted objectType " + TYPE_AUDITACTION + " for AuditAction")
	}

	result := t.createObject(stub, jsonBytes, newAuditAction.ID)

	if result.Status != shim.OK {
		fmt.Println("- end createAuditAction (failed)")
		return result
	}

	result = t.putCompositeKey(stub, CK_AUDITOR_AUDIT, []string{newAuditAction.Auditor, newAuditAction.ID})
	if result.Status != shim.OK {
		fmt.Println("- end createAuditAction (failed)")
		return result
	}

	result = t.putCompositeKey(stub, CK_AUDIT_OBJ, []string{newAuditAction.ObjectID, newAuditAction.ID})
	if result.Status != shim.OK {
		fmt.Println("- end createAuditAction (failed)")
		return result
	}

	fmt.Println("- end createAuditAction (success)")

	return result
}

func (t *FoodChaincode) updateAuditAction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateAuditAction", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newAuditActions := AuditAction{}
	err := json.Unmarshal(jsonBytes, &newAuditActions)
	if err != nil {
		return shim.Error("Failed to decode json of AuditAction: " + err.Error())
	}

	if newAuditActions.ObjectType != TYPE_AUDITACTION {
		return shim.Error("Expexted objectType " + TYPE_AUDITACTION + " for AuditAction")
	}

	result := t.updateObject(stub, jsonBytes, newAuditActions.ID)

	if result.Status == shim.OK {
		fmt.Println("- end updateAuditAction (success)")
	}
	return shim.Success(nil)
}

// Query methods
// ========================================
func (t *FoodChaincode) getObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start getObject", args)
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ID := args[0]
	objectType := args[1]
	existedObjectAsBytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Failed to get existed Object with ID: " + ID + ", error: " + err.Error())
	} else if existedObjectAsBytes == nil {
		return shim.Error("Object with ID " + ID + " does not exist")
	}

	liteModel := LiteModel{}
	err = json.Unmarshal(existedObjectAsBytes, &liteModel)
	if err != nil {
		return shim.Error("Failed to get decode object: " + err.Error())
	}
	if liteModel.ObjectType != objectType {
		return shim.Error("ObjectType does not match")
	}

	fmt.Println("- end getObject (success)")
	return shim.Success(existedObjectAsBytes)
}

func (t *FoodChaincode) getLogsOfSupplychain(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start getLogsOfSupplychain", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]

	existedObjectAsBytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Failed to get existed Supplychain with ID: " + ID + ", error: " + err.Error())
	} else if existedObjectAsBytes == nil {
		return shim.Error("Supplychain with ID " + ID + " does not exist")
	}

	sc := Traceable{}
	err = json.Unmarshal(existedObjectAsBytes, &sc)
	if err != nil {
		return shim.Error("Failed to get decode object: " + err.Error())
	}

	if sc.ObjectType != TYPE_SUPPLYCHAIN {
		return shim.Error("Object with ID: " + ID + "is not a Supplychain")
	}

	resultsIterator, err := stub.GetStateByPartialCompositeKey(CK_SC_LOG, []string{ID})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	result, responseAsBytes := t.getLogsFromIterator(stub, resultsIterator)
	if err != nil {
		return shim.Error("Failed to get encode response: " + err.Error())
	}

	if result.Status != shim.OK {
		fmt.Println("- end getLogsOfSupplychain (failed)")
		return result
	}

	fmt.Println("- end getLogsOfSupplychain (success)")
	return shim.Success(responseAsBytes)
}

func (t *FoodChaincode) getLogsOfProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start getLogsOfProduct", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]

	existedObjectAsBytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Failed to get existed Product with ID: " + ID + ", error: " + err.Error())
	} else if existedObjectAsBytes == nil {
		return shim.Error("Product with ID " + ID + " does not exist")
	}

	product := Traceable{}
	err = json.Unmarshal(existedObjectAsBytes, &product)
	if err != nil {
		return shim.Error("Failed to get decode object: " + err.Error())
	}

	if product.ObjectType != TYPE_PRODUCT {
		return shim.Error("Object with ID: " + ID + "is not a Product")
	}

	resultsIterator, err := stub.GetStateByPartialCompositeKey(CK_PRODUCT_LOG, []string{ID})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	result, responseAsBytes := t.getLogsFromIterator(stub, resultsIterator)
	if err != nil {
		return shim.Error("Failed to get encode response: " + err.Error())
	}

	if result.Status != shim.OK {
		fmt.Println("- end getLogsOfProduct (failed)")
		return result
	}

	fmt.Println("- end getLogsOfProduct (success)")
	return shim.Success(responseAsBytes)
}

func (t *FoodChaincode) getAuditOfObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start getAuditOfObject", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]
	resultsIterator, err := stub.GetStateByPartialCompositeKey(CK_AUDIT_OBJ, []string{ID})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedAuditID := compositeKeyParts[1]
		fmt.Printf("- found an audit:%s\n", returnedAuditID)

		auditAsBytes, err := stub.GetState(returnedAuditID)
		if err != nil {
			return shim.Error("Failed to get existed Audit with ID: " + returnedAuditID + ", error: " + err.Error())
		} else if auditAsBytes == nil {
			return shim.Error("Audit with ID " + returnedAuditID + " does not exist")
		}
		fmt.Println("- end getAuditOfObject (success)")
		return shim.Success(auditAsBytes)
	}

	fmt.Println("- end getAuditOfObject (failed)")
	return shim.Error("No audit was found")
}

func (t *FoodChaincode) getAuditsOfAuditor(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start getAuditsOfAuditor", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]
	resultsIterator, err := stub.GetStateByPartialCompositeKey(CK_AUDITOR_AUDIT, []string{ID})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	response := []AuditAction{}
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedAuditID := compositeKeyParts[1]
		fmt.Printf("- found an audit:%s\n", returnedAuditID)

		auditAsBytes, err := stub.GetState(returnedAuditID)
		if err != nil {
			return shim.Error("Failed to get existed Audit with ID: " + returnedAuditID + ", error: " + err.Error())
		} else if auditAsBytes == nil {
			return shim.Error("Audit with ID " + returnedAuditID + " does not exist")
		}

		audit := AuditAction{}
		err = json.Unmarshal(auditAsBytes, &audit)
		if err != nil {
			return shim.Error("Failed to get decode audit: " + err.Error())
		}
		response = append(response, audit)
	}

	responseAsBytes, err := json.Marshal(response)
	if err != nil {
		return shim.Error("Failed to get encode response: " + err.Error())
	}

	fmt.Println("- end getAuditsOfAuditor (success)")
	return shim.Success(responseAsBytes)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func (t *FoodChaincode) getQueryResultForQueryString(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Printf("- getQueryResultForQueryString args:\n%s\n", args)

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
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

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// =========================================================================================
// getHistoryOfObject
// =========================================================================================
func (t *FoodChaincode) getHistoryOfObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Printf("- getHistoryOfObject args:\n%s\n", args)

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]

	resultsIterator, err := stub.GetHistoryForKey(ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the object
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

	fmt.Printf("- getHistoryOfObject returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// Helper methods
// ========================================
func (t *FoodChaincode) createObject(stub shim.ChaincodeStubInterface, bytes []byte, ID string) pb.Response {
	existedObjectAsBytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Failed to get existed Object with ID: " + ID + ", error: " + err.Error())
	} else if existedObjectAsBytes != nil {
		return shim.Error("Object with ID " + ID + " already existed")
	}

	err = stub.PutState(ID, bytes)
	if err != nil {
		return shim.Error("Failed to create new object with ID: " + ID + ", error: " + err.Error())
	}

	return shim.Success(nil)
}

func (t *FoodChaincode) updateObject(stub shim.ChaincodeStubInterface, bytes []byte, ID string) pb.Response {
	existedObjectAsBytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Failed to get existed Object with ID: " + ID + ", error: " + err.Error())
	} else if existedObjectAsBytes == nil {
		return shim.Error("Object with ID " + ID + " does not exist")
	}

	err = stub.PutState(ID, bytes)
	if err != nil {
		return shim.Error("Failed to update the object with ID: " + ID + ", error: " + err.Error())
	}
	return shim.Success(nil)
}

func (t *FoodChaincode) updateLogHandler(stub shim.ChaincodeStubInterface, bytes []byte, newLog Log) pb.Response {
	existedObjectAsBytes, err := stub.GetState(newLog.ID)
	if err != nil {
		return shim.Error("Failed to get existed Object with ID: " + newLog.ID + ", error: " + err.Error())
	} else if existedObjectAsBytes == nil {
		return shim.Error("Object with ID " + newLog.ID + " does not exist")
	}

	oldLog := Log{}
	err = json.Unmarshal(existedObjectAsBytes, &oldLog)
	if err != nil {
		return shim.Error("Failed to get decode Log: " + err.Error())
	}

	var result pb.Response

	if len(newLog.Supplychain) > 0 {
		result = t.updateCompositeKey(
			stub,
			CK_SC_LOG,
			[]string{oldLog.Supplychain, oldLog.ID},
			[]string{newLog.Supplychain, newLog.ID})
		if result.Status != shim.OK {
			fmt.Println("- end updateLog (failed)")
			return result
		}
	} else if len(oldLog.Supplychain) > 0 {
		result = t.deleteCompositeKey(stub, CK_SC_LOG, []string{oldLog.Supplychain, oldLog.ID})
		if result.Status != shim.OK {
			fmt.Println("- end updateLog (failed)")
			return result
		}
	}

	if len(newLog.Product) > 0 {
		result = t.updateCompositeKey(
			stub,
			CK_PRODUCT_LOG,
			[]string{oldLog.Product, oldLog.ID},
			[]string{newLog.Product, newLog.ID})
		if result.Status != shim.OK {
			fmt.Println("- end updateLog (failed)")
			return result
		}
	} else if len(oldLog.Product) > 0 {
		result = t.deleteCompositeKey(stub, CK_PRODUCT_LOG, []string{oldLog.Product, oldLog.ID})
		if result.Status != shim.OK {
			fmt.Println("- end updateLog (failed)")
			return result
		}
	}

	err = stub.PutState(newLog.ID, bytes)
	if err != nil {
		return shim.Error("Failed to update the object with ID: " + newLog.ID + ", error: " + err.Error())
	}
	return shim.Success(nil)
}

func (t *FoodChaincode) putCompositeKey(stub shim.ChaincodeStubInterface, objectType string, values []string) pb.Response {
	cKey, err := stub.CreateCompositeKey(objectType, values)
	if err != nil {
		return shim.Error("Failed to create composite key: " + err.Error())
	}
	fmt.Println("save new composite key:", cKey)
	err = stub.PutState(cKey, []byte{0x00})
	if err != nil {
		return shim.Error("Failed to save composite key: " + err.Error())
	}
	return shim.Success(nil)
}

func (t *FoodChaincode) updateCompositeKey(stub shim.ChaincodeStubInterface, objectType string, oldValues []string, newValues []string) pb.Response {

	valueChanged := false
	for i, value := range oldValues {
		if newValues[i] != value {
			valueChanged = true
			break
		}
	}

	if !valueChanged {
		fmt.Println("Values don't change. Don't create new composite key")
		return shim.Success(nil)
	}

	result := t.deleteCompositeKey(stub, objectType, oldValues)
	if result.Status != shim.OK {
		return result
	}

	return t.putCompositeKey(stub, objectType, newValues)
}

func (t *FoodChaincode) deleteCompositeKey(stub shim.ChaincodeStubInterface, objectType string, oldValues []string) pb.Response {
	cKey, err := stub.CreateCompositeKey(objectType, oldValues)
	if err != nil {
		return shim.Error("Failed to create composite key: " + err.Error())
	}
	err = stub.DelState(cKey)
	if err != nil {
		return shim.Error("Failed to delete composite key: " + err.Error())
	}
	fmt.Println("Deleted old composite key: ", cKey)

	return shim.Success(nil)
}

func (t *FoodChaincode) getLogsFromIterator(stub shim.ChaincodeStubInterface, resultsIterator shim.StateQueryIteratorInterface) (pb.Response, []byte) {
	response := []Log{}
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error()), nil
		}

		_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error()), nil
		}
		returnedLogID := compositeKeyParts[1]
		fmt.Printf("- found an log:%s\n", returnedLogID)

		logAsBytes, err := stub.GetState(returnedLogID)
		if err != nil {
			return shim.Error("Failed to get existed Log with ID: " + returnedLogID + ", error: " + err.Error()), nil
		} else if logAsBytes == nil {
			return shim.Error("Log with ID " + returnedLogID + " does not exist"), nil
		}
		log := Log{}
		err = json.Unmarshal(logAsBytes, &log)
		if err != nil {
			return shim.Error("Failed to get decode log: " + err.Error()), nil
		}
		response = append(response, log)
	}

	responseAsBytes, err := json.Marshal(response)
	if err != nil {
		return shim.Error("Failed to get encode response: " + err.Error()), nil
	}

	return shim.Success(nil), responseAsBytes
}
