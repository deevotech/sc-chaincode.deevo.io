package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreareLog(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newLog := Log{
		ObjectType: TYPE_LOG,
		ID:         "Log_1",
		Time:       time.Now().Unix(),
		Ref:        []string{"Product_1"},
		CTE:        "test_action",
		Content:    "Log 1",
		Asset:      "Asset_1",
		Product:    "Product_1",
		Location:   "Location_1",
	}

	newLogAsBytes, err := json.Marshal(newLog)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateLog(t, stub, newLogAsBytes, newLog)
}

func TestFood_UpdateLog(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newLog := Log{
		ObjectType: TYPE_LOG,
		ID:         "Log_1",
		Time:       time.Now().Unix(),
		Ref:        []string{"Product_1"},
		CTE:        "test_action",
		Content:    "Log 1",
		Asset:      "Asset_1",
		Product:    "Product_1",
		Location:   "Location_1",
	}
	newLogAsBytes, err := json.Marshal(newLog)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateLog(t, stub, newLogAsBytes, newLog)

	updatedLog := Log{
		ObjectType: TYPE_LOG,
		ID:         "Log_1",
		Time:       time.Now().Unix(),
		Ref:        []string{"Product_1", "product_2"},
		CTE:        "test_action",
		Content:    "Log 2",
		Asset:      "Asset_1",
		Product:    "Product_1",
		Location:   "Location_1",
	}
	updatedLogAsBytes, err := json.Marshal(updatedLog)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateLog(t, stub, updatedLogAsBytes, updatedLog)
}

func checkCreateLog(t *testing.T, stub *shim.MockStub, logAsJSON []byte, value Log) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createLog"), logAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resLog := Log{}
	err := json.Unmarshal(bytes, &resLog)
	if err != nil {
		fmt.Println("Failed to decode json of Log:", err.Error())
		t.FailNow()
	}
	if !resLog.Equals(value) {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkUpdateLog(t *testing.T, stub *shim.MockStub, logAsJSON []byte, value Log) {
	res := stub.MockInvoke("1", [][]byte{[]byte("updateLog"), logAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resLog := Log{}
	err := json.Unmarshal(bytes, &resLog)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if !resLog.Equals(value) {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
