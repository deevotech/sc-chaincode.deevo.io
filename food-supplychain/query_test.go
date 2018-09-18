package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_QueryLog(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newLog := Log{ObjectType: "Log", ID: "Log_1", Content: "Log 1", Time: time.Now().Unix(), Location: "Location_1", ObjectID: "Product_1"}
	newLogAsBytes, err := json.Marshal(newLog)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateLog(t, stub, newLogAsBytes, newLog)

	res := stub.MockInvoke("1", [][]byte{[]byte("getObject"), []byte(newLog.ID)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload == nil {
		fmt.Println("failed: no payload")
		t.FailNow()
	}

	resLog := Log{}
	err = json.Unmarshal(res.Payload, &resLog)
	if err != nil {
		fmt.Println("Failed to decode json of Log:", err.Error())
		t.FailNow()
	}
	if resLog != newLog {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func TestFood_QueryAudit(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newLog := Log{ObjectType: "Log", ID: "Log_1", Content: "Log 1", Time: time.Now().Unix(), Location: "Location_1", ObjectID: "Product_1"}
	newLogAsBytes, err := json.Marshal(newLog)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateLog(t, stub, newLogAsBytes, newLog)

	newAuditAction := AuditAction{ObjectType: "AuditAction", ID: "AuditAction_1", Time: time.Now().Unix(), Location: "Location_1", ObjectID: "Log_1"}
	newAuditActionAsBytes, err := json.Marshal(newAuditAction)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateAuditAction(t, stub, newAuditActionAsBytes, newAuditAction)

	res := stub.MockInvoke("1", [][]byte{[]byte("getAuditOfObject"), []byte(newLog.ID)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	if res.Payload == nil {
		fmt.Println("failed: no payload")
		t.FailNow()
	}

	resAudit := AuditAction{}
	err = json.Unmarshal(res.Payload, &resAudit)
	if err != nil {
		fmt.Println("Failed to decode json of Log:", err.Error())
		t.FailNow()
	}
	if resAudit != newAuditAction {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
