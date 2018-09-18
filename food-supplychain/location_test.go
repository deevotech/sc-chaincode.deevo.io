package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreareLocation(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newLocation := Location{ObjectType: "Location", ID: "Location_1", Name: "Location 1", Party: "Party_1"}
	newLocationAsBytes, err := json.Marshal(newLocation)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateLocation(t, stub, newLocationAsBytes, newLocation)
}

func TestFood_UpdateLocation(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newLocation := Location{ObjectType: "Location", ID: "Location_1", Name: "Location 1", Party: "Party_1"}
	newLocationAsBytes, err := json.Marshal(newLocation)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateLocation(t, stub, newLocationAsBytes, newLocation)

	updatedLocation := Location{ObjectType: "Location", ID: "Location_1", Name: "Location 2", Party: "Party_1"}
	updatedLocationAsBytes, err := json.Marshal(updatedLocation)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateLocation(t, stub, updatedLocationAsBytes, updatedLocation)
}

func checkCreateLocation(t *testing.T, stub *shim.MockStub, locationAsJSON []byte, value Location) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createLocation"), locationAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resLocation := Location{}
	err := json.Unmarshal(bytes, &resLocation)
	if err != nil {
		fmt.Println("Failed to decode json of Location:", err.Error())
		t.FailNow()
	}
	if resLocation != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkUpdateLocation(t *testing.T, stub *shim.MockStub, locationAsJSON []byte, value Location) {
	res := stub.MockInvoke("1", [][]byte{[]byte("updateLocation"), locationAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resLocation := Location{}
	err := json.Unmarshal(bytes, &resLocation)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if resLocation != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
