package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_GetLogHistory(t *testing.T) {
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

	checkLogHistory(t, stub, updatedLog.ID)
}

func checkLogHistory(t *testing.T, stub *shim.MockStub, id string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getHistoryOfObject"), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}
}
