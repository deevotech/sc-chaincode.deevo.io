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

// Org model
type Org struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Address    string `json:"address"`
}

// Party model
type Party struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	ORG        string `json:"org"`
}

// Location model
type Location struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Party      string `json:"party"`
}

// Product model
type Product struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
}

// Log model
type Log struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Content    string `json:"content"`
	Time       int64  `json:"name"`
	Location   string `json:"location"`
	ObjectID   string `json:"objectID"`
}

// InitData model
type InitData struct {
	ORG       Org      `json:"org"`
	Party1    Party    `json:"party_1"`
	Party2    Party    `json:"party_2"`
	Location1 Location `json:"location_1"`
	Location2 Location `json:"location_2"`
	Product1  Product  `json:"product_1"`
	Product2  Product  `json:"product_2"`
}

// Auditor model
type Auditor struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	Name       string `json:"name"`
}

// AuditAction model
type AuditAction struct {
	ObjectType string  `json:"docType"`
	ID         string  `json:"id"`
	Time       int64   `json:"time"`
	Creator    Auditor `json:"auditor"`
	Location   string  `json:"location"`
	ObjectID   string  `json:"objectID"`
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
	} else if function == "createORG" { //create a new org
		return t.createORG(stub, args)
	} else if function == "updateORG" {
		return t.updateORG(stub, args)
	} else if function == "createParty" { //create a new org
		return t.createParty(stub, args)
	} else if function == "updateParty" {
		return t.updateParty(stub, args)
	} else if function == "createLocation" { //create a new org
		return t.createLocation(stub, args)
	} else if function == "updateLocation" {
		return t.updateLocation(stub, args)
	} else if function == "createProduct" { //create a new org
		return t.createProduct(stub, args)
	} else if function == "updateProduct" {
		return t.updateProduct(stub, args)
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

	newOrg := newData.ORG
	newOrgAsBytes, err := json.Marshal(newOrg)
	if err != nil {
		return shim.Error("Failed to encode json of object " + newOrg.ID)
	}
	result := t.createObject(stub, newOrgAsBytes, newOrg.ID)
	if result.Status != 200 {
		return result
	}

	newParty1 := newData.Party1
	newParty1AsBytes, err := json.Marshal(newParty1)
	if err != nil {
		return shim.Error("Failed to encode json of object " + newParty1.ID)
	}
	result = t.createObject(stub, newParty1AsBytes, newParty1.ID)
	if result.Status != 200 {
		return result
	}

	newParty2 := newData.Party2
	newParty2AsBytes, err := json.Marshal(newParty2)
	if err != nil {
		return shim.Error("Failed to encode json of object " + newParty2.ID)
	}
	result = t.createObject(stub, newParty2AsBytes, newParty2.ID)
	if result.Status != 200 {
		return result
	}

	newLocation1 := newData.Location1
	newLocation1AsBytes, err := json.Marshal(newLocation1)
	if err != nil {
		return shim.Error("Failed to encode json of object " + newLocation1.ID)
	}
	result = t.createObject(stub, newLocation1AsBytes, newLocation1.ID)
	if result.Status != 200 {
		return result
	}

	newLocation2 := newData.Location2
	newLocation2AsBytes, err := json.Marshal(newLocation2)
	if err != nil {
		return shim.Error("Failed to encode json of object " + newLocation2.ID)
	}
	result = t.createObject(stub, newLocation2AsBytes, newLocation2.ID)
	if result.Status != 200 {
		return result
	}

	newProduct1 := newData.Product1
	newProduct1AsBytes, err := json.Marshal(newProduct1)
	if err != nil {
		return shim.Error("Failed to encode json of object " + newProduct1.ID)
	}
	result = t.createObject(stub, newProduct1AsBytes, newProduct1.ID)
	if result.Status != 200 {
		return result
	}

	newProduct2 := newData.Product2
	newProduct2AsBytes, err := json.Marshal(newProduct2)
	if err != nil {
		return shim.Error("Failed to encode json of object " + newProduct2.ID)
	}
	result = t.createObject(stub, newProduct2AsBytes, newProduct2.ID)
	if result.Status != 200 {
		return result
	}

	fmt.Println("- end initOrgData (success)")
	return shim.Success(nil)
}

// Methods on ORG
// ========================================
func (t *FoodChaincode) createORG(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createORG", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newOrg := Org{}
	err := json.Unmarshal(jsonBytes, &newOrg)
	if err != nil {
		return shim.Error("Failed to decode json of ORG: " + err.Error())
	}

	result := t.createObject(stub, jsonBytes, newOrg.ID)

	if result.Status == 200 {
		fmt.Println("- end createORG (success)")
	}
	return result
}

func (t *FoodChaincode) updateORG(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateORG", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newOrg := Org{}
	err := json.Unmarshal(jsonBytes, &newOrg)
	if err != nil {
		return shim.Error("Failed to decode json of ORG: " + err.Error())
	}

	result := t.updateObject(stub, jsonBytes, newOrg.ID)

	if result.Status == 200 {
		fmt.Println("- end updateORG (success)")
	}
	return shim.Success(nil)
}

// Methods on Party
// ========================================
func (t *FoodChaincode) createParty(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createParty", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newParty := Party{}
	err := json.Unmarshal(jsonBytes, &newParty)
	if err != nil {
		return shim.Error("Failed to decode json of Party: " + err.Error())
	}

	result := t.createObject(stub, jsonBytes, newParty.ID)

	if result.Status == 200 {
		fmt.Println("- end createParty (success)")
	}
	return result
}

func (t *FoodChaincode) updateParty(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateParty", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newParty := Party{}
	err := json.Unmarshal(jsonBytes, &newParty)
	if err != nil {
		return shim.Error("Failed to decode json of Party: " + err.Error())
	}

	result := t.updateObject(stub, jsonBytes, newParty.ID)

	if result.Status == 200 {
		fmt.Println("- end updateParty (success)")
	}
	return shim.Success(nil)
}

// Methods on Location
// ========================================
func (t *FoodChaincode) createLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createLocation", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newLocation := Location{}
	err := json.Unmarshal(jsonBytes, &newLocation)
	if err != nil {
		return shim.Error("Failed to decode json of Location: " + err.Error())
	}

	result := t.createObject(stub, jsonBytes, newLocation.ID)

	if result.Status == 200 {
		fmt.Println("- end createLocation (success)")
	}
	return result
}

func (t *FoodChaincode) updateLocation(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateLocation", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newLocation := Location{}
	err := json.Unmarshal(jsonBytes, &newLocation)
	if err != nil {
		return shim.Error("Failed to decode json of Location: " + err.Error())
	}

	result := t.updateObject(stub, jsonBytes, newLocation.ID)

	if result.Status == 200 {
		fmt.Println("- end updateLocation (success)")
	}
	return shim.Success(nil)
}

// Methods on Product
// ========================================
func (t *FoodChaincode) createProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start createProduct", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newProduct := Product{}
	err := json.Unmarshal(jsonBytes, &newProduct)
	if err != nil {
		return shim.Error("Failed to decode json of Product: " + err.Error())
	}

	result := t.createObject(stub, jsonBytes, newProduct.ID)

	if result.Status == 200 {
		fmt.Println("- end createProduct (success)")
	}
	return result
}

func (t *FoodChaincode) updateProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start updateProduct", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	jsonBytes := []byte(args[0])
	newProduct := Product{}
	err := json.Unmarshal(jsonBytes, &newProduct)
	if err != nil {
		return shim.Error("Failed to decode json of Product: " + err.Error())
	}

	result := t.updateObject(stub, jsonBytes, newProduct.ID)

	if result.Status == 200 {
		fmt.Println("- end updateProduct (success)")
	}
	return shim.Success(nil)
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

	result := t.createObject(stub, jsonBytes, newLog.ID)

	if result.Status == 200 {
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

	result := t.updateObject(stub, jsonBytes, newLog.ID)

	if result.Status == 200 {
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

	result := t.createObject(stub, jsonBytes, newAuditor.ID)

	if result.Status == 200 {
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

	result := t.updateObject(stub, jsonBytes, newAuditor.ID)

	if result.Status == 200 {
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

	result := t.createObject(stub, jsonBytes, newAuditAction.ID)

	cKey, err := stub.CreateCompositeKey("auditedObject~audit", []string{newAuditAction.ObjectID, newAuditAction.ID})
	if err != nil {
		return shim.Error("Failed to create composite key: " + err.Error())
	}
	fmt.Println("save new composite key:", cKey)
	err = stub.PutState(cKey, []byte{0x00})
	if err != nil {
		return shim.Error("Failed to save composite key: " + err.Error())
	}

	if result.Status == 200 {
		fmt.Println("- end createAuditAction (success)")
	}
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

	result := t.updateObject(stub, jsonBytes, newAuditActions.ID)

	if result.Status == 200 {
		fmt.Println("- end updateAuditAction (success)")
	}
	return shim.Success(nil)
}

// Query methods
// ========================================
func (t *FoodChaincode) getObject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start getObject", args)
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]
	existedObjectAsBytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Failed to get existed Object with ID: " + ID + ", error: " + err.Error())
	} else if existedObjectAsBytes == nil {
		return shim.Error("Object with ID " + ID + " does not exist")
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
	resultsIterator, err := stub.GetStateByPartialCompositeKey("auditedObject~audit", []string{ID})
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
