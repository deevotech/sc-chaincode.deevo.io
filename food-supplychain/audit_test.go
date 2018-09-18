package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreareAuditor(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newAuditor := Auditor{ObjectType: "Auditor", ID: "Auditor_1", Name: "Auditor 1"}
	newAuditorAsBytes, err := json.Marshal(newAuditor)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateAuditor(t, stub, newAuditorAsBytes, newAuditor)
}
func TestFood_UpdateAuditor(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newAuditor := Auditor{ObjectType: "Auditor", ID: "Auditor_1", Name: "Auditor 1"}
	newAuditorAsBytes, err := json.Marshal(newAuditor)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateAuditor(t, stub, newAuditorAsBytes, newAuditor)

	updatedAuditor := Auditor{ObjectType: "Auditor", ID: "Auditor_1", Name: "Auditor 2"}
	updatedAuditorAsBytes, err := json.Marshal(updatedAuditor)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateAuditor(t, stub, updatedAuditorAsBytes, updatedAuditor)
}

func TestFood_CreareAuditAction(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newAuditAction := AuditAction{ObjectType: "AuditAction", ID: "AuditAction_1", Auditor: "Auditor_1", Time: time.Now().Unix(), Location: "Location_1", ObjectID: "Product_1"}
	newAuditActionAsBytes, err := json.Marshal(newAuditAction)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateAuditAction(t, stub, newAuditActionAsBytes, newAuditAction)
}

func TestFood_UpdateAuditAction(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newAuditAction := AuditAction{ObjectType: "AuditAction", ID: "AuditAction_1", Auditor: "Auditor_1", Time: time.Now().Unix(), Location: "Location_1", ObjectID: "Product_1"}
	newAuditActionAsBytes, err := json.Marshal(newAuditAction)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateAuditAction(t, stub, newAuditActionAsBytes, newAuditAction)

	updatedAuditAction := AuditAction{ObjectType: "AuditAction", ID: "AuditAction_1", Auditor: "Auditor_1", Time: time.Now().Unix(), Location: "Location_2", ObjectID: "Product_2"}
	updatedAuditActionAsBytes, err := json.Marshal(updatedAuditAction)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateAuditAction(t, stub, updatedAuditActionAsBytes, updatedAuditAction)
}
func checkCreateAuditor(t *testing.T, stub *shim.MockStub, auditorAsJSON []byte, value Auditor) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createAuditor"), auditorAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resAuditor := Auditor{}
	err := json.Unmarshal(bytes, &resAuditor)
	if err != nil {
		fmt.Println("Failed to decode json of Auditor:", err.Error())
		t.FailNow()
	}
	if resAuditor != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkUpdateAuditor(t *testing.T, stub *shim.MockStub, auditorAsJSON []byte, value Auditor) {
	res := stub.MockInvoke("1", [][]byte{[]byte("updateAuditor"), auditorAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resAuditor := Auditor{}
	err := json.Unmarshal(bytes, &resAuditor)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if resAuditor != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkCreateAuditAction(t *testing.T, stub *shim.MockStub, auditActionAsJSON []byte, value AuditAction) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createAuditAction"), auditActionAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resAuditAction := AuditAction{}
	err := json.Unmarshal(bytes, &resAuditAction)
	if err != nil {
		fmt.Println("Failed to decode json of AuditAction:", err.Error())
		t.FailNow()
	}
	if resAuditAction != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkUpdateAuditAction(t *testing.T, stub *shim.MockStub, auditActionAsJSON []byte, value AuditAction) {
	res := stub.MockInvoke("1", [][]byte{[]byte("updateAuditAction"), auditActionAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resAuditAction := AuditAction{}
	err := json.Unmarshal(bytes, &resAuditAction)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if resAuditAction != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
