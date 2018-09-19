package main

import (
	"encoding/json"
	"fmt"

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

	result := t.updateObject(stub, jsonBytes, newLog.ID)

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
			return shim.Error("Failed to get existed Audit with ID: " + ID + ", error: " + err.Error())
		} else if auditAsBytes == nil {
			return shim.Error("Audit with ID " + ID + " does not exist")
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
			return shim.Error("Failed to get existed Audit with ID: " + ID + ", error: " + err.Error())
		} else if auditAsBytes == nil {
			return shim.Error("Audit with ID " + ID + " does not exist")
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
