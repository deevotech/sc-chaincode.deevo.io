package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreateLogs(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newSupplychain := Traceable{ObjectType: TYPE_SUPPLYCHAIN, ID: "sc_1", Name: "supplychain 1"}
	newSupplychainAsBytes, err := json.Marshal(newSupplychain)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateTraceable(t, stub, newSupplychainAsBytes, newSupplychain)

	newProduct := Traceable{ObjectType: TYPE_PRODUCT, ID: "Product_1", Name: "Product 1"}
	newProductAsBytes, err := json.Marshal(newProduct)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateTraceable(t, stub, newProductAsBytes, newProduct)

	newLog := Log{
		ObjectType:  TYPE_LOG,
		ID:          "Log_1",
		Time:        time.Now().Unix(),
		Ref:         []string{"Product_1"},
		CTE:         "test_action",
		Supplychain: "sc_1",
		Content:     "Log 1",
		Asset:       "Asset_1",
		Product:     "Product_1",
		Location:    "Location_1",
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

	newSupplychain := Traceable{ObjectType: TYPE_SUPPLYCHAIN, ID: "sc_1", Name: "supplychain 1"}
	newSupplychainAsBytes, err := json.Marshal(newSupplychain)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateTraceable(t, stub, newSupplychainAsBytes, newSupplychain)

	newProduct := Traceable{ObjectType: TYPE_PRODUCT, ID: "Product_1", Name: "Product 1"}
	newProductAsBytes, err := json.Marshal(newProduct)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateTraceable(t, stub, newProductAsBytes, newProduct)

	newLog := Log{
		ObjectType:  TYPE_LOG,
		ID:          "Log_1",
		Time:        time.Now().Unix(),
		Ref:         []string{"Product_1"},
		CTE:         "test_action",
		Supplychain: "sc_1",
		Content:     "Log 1",
		Asset:       "Asset_1",
		Product:     "Product_1",
		Location:    "Location_1",
	}
	newLogAsBytes, err := json.Marshal(newLog)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateLog(t, stub, newLogAsBytes, newLog)

	newProduct2 := Traceable{ObjectType: TYPE_PRODUCT, ID: "Product_2", Name: "Product 2"}
	newProduct2AsBytes, err := json.Marshal(newProduct2)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateTraceable(t, stub, newProduct2AsBytes, newProduct2)

	updatedLog := Log{
		ObjectType: TYPE_LOG,
		ID:         "Log_1",
		Time:       time.Now().Unix(),
		Ref:        []string{"Product_1", "product_2"},
		CTE:        "test_action",
		Content:    "Log 2",
		Asset:      "Asset_1",
		Product:    "Product_2",
		Location:   "Location_1",
	}
	updatedLogAsBytes, err := json.Marshal(updatedLog)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateLog(t, stub, updatedLogAsBytes, updatedLog, newLog)
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

	// check log of supplychain
	res = stub.MockInvoke("1", [][]byte{[]byte("getLogsOfSupplychain"), []byte(value.Supplychain)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	resLogs := []Log{}
	err = json.Unmarshal(res.Payload, &resLogs)
	if err != nil {
		fmt.Println("Failed to decode json of Logs:", err.Error())
		t.FailNow()
	}

	if len(resLogs) != 1 {
		fmt.Println("Size of response does not match")
		t.FailNow()
	}

	if !resLogs[0].Equals(value) {
		fmt.Println("Response does not match")
		t.FailNow()
	}

	// check log of product
	res = stub.MockInvoke("1", [][]byte{[]byte("getLogsOfProduct"), []byte(value.Product)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	resLogs = []Log{}
	err = json.Unmarshal(res.Payload, &resLogs)
	if err != nil {
		fmt.Println("Failed to decode json of Logs:", err.Error())
		t.FailNow()
	}

	if len(resLogs) != 1 {
		fmt.Println("Size of response does not match")
		t.FailNow()
	}

	if !resLogs[0].Equals(value) {
		fmt.Println("Response does not match")
		t.FailNow()
	}
}

func checkUpdateLog(t *testing.T, stub *shim.MockStub, logAsJSON []byte, value Log, oldValue Log) {
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

	// check log of supplychain
	res = stub.MockInvoke("1", [][]byte{[]byte("getLogsOfSupplychain"), []byte(oldValue.Supplychain)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	resLogs := []Log{}
	err = json.Unmarshal(res.Payload, &resLogs)
	if err != nil {
		fmt.Println("Failed to decode json of Logs:", err.Error())
		t.FailNow()
	}

	if len(resLogs) != 0 {
		fmt.Println("Size of response does not match")
		t.FailNow()
	}

	// check log of old product
	res = stub.MockInvoke("1", [][]byte{[]byte("getLogsOfProduct"), []byte(oldValue.Product)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	resLogs = []Log{}
	err = json.Unmarshal(res.Payload, &resLogs)
	if err != nil {
		fmt.Println("Failed to decode json of Logs:", err.Error())
		t.FailNow()
	}

	if len(resLogs) != 0 {
		fmt.Println("Size of response does not match")
		t.FailNow()
	}

	// check log of product
	res = stub.MockInvoke("1", [][]byte{[]byte("getLogsOfProduct"), []byte(value.Product)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	resLogs = []Log{}
	err = json.Unmarshal(res.Payload, &resLogs)
	if err != nil {
		fmt.Println("Failed to decode json of Logs:", err.Error())
		t.FailNow()
	}

	if len(resLogs) != 1 {
		fmt.Println("Size of response does not match")
		t.FailNow()
	}

	if !resLogs[0].Equals(value) {
		fmt.Println("Response does not match")
		t.FailNow()
	}
}
