package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreareProduct(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newProduct := Product{ObjectType: "Product", ID: "Product_1", Name: "Product 1", Location: "Location_1"}
	newProductAsBytes, err := json.Marshal(newProduct)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateProduct(t, stub, newProductAsBytes, newProduct)
}

func TestFood_UpdateProduct(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newProduct := Product{ObjectType: "Product", ID: "Product_1", Name: "Product 1", Location: "Location_1"}
	newProductAsBytes, err := json.Marshal(newProduct)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateProduct(t, stub, newProductAsBytes, newProduct)

	updatedProduct := Product{ObjectType: "Product", ID: "Product_1", Name: "Product 2", Location: "Location_1"}
	updatedProductAsBytes, err := json.Marshal(updatedProduct)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateProduct(t, stub, updatedProductAsBytes, updatedProduct)
}
func checkCreateProduct(t *testing.T, stub *shim.MockStub, productAsJSON []byte, value Product) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createProduct"), productAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resProduct := Product{}
	err := json.Unmarshal(bytes, &resProduct)
	if err != nil {
		fmt.Println("Failed to decode json of Product:", err.Error())
		t.FailNow()
	}
	if resProduct != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkUpdateProduct(t *testing.T, stub *shim.MockStub, productAsJSON []byte, value Product) {
	res := stub.MockInvoke("1", [][]byte{[]byte("updateProduct"), productAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resProduct := Product{}
	err := json.Unmarshal(bytes, &resProduct)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if resProduct != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
