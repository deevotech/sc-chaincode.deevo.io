package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreareORG(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newOrg := Org{ObjectType: "org", ID: "org_1", Name: "org 1", Address: "address 1"}
	newOrgAsBytes, err := json.Marshal(newOrg)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateORG(t, stub, newOrgAsBytes, newOrg)
}

func TestFood_UpdateORG(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newOrg := Org{ObjectType: "org", ID: "org_1", Name: "org 1", Address: "address 1"}
	newOrgAsBytes, err := json.Marshal(newOrg)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateORG(t, stub, newOrgAsBytes, newOrg)

	updatedOrg := Org{ObjectType: "org", ID: "org_1", Name: "org 1", Address: "address 2"}
	updatedOrgAsBytes, err := json.Marshal(updatedOrg)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateORG(t, stub, updatedOrgAsBytes, updatedOrg)
}

func checkCreateORG(t *testing.T, stub *shim.MockStub, orgAsJSON []byte, value Org) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createORG"), orgAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resOrg := Org{}
	err := json.Unmarshal(bytes, &resOrg)
	if err != nil {
		fmt.Println("Failed to decode json of ORG:", err.Error())
		t.FailNow()
	}
	if resOrg != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkUpdateORG(t *testing.T, stub *shim.MockStub, orgAsJSON []byte, value Org) {
	res := stub.MockInvoke("1", [][]byte{[]byte("updateORG"), orgAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resOrg := Org{}
	err := json.Unmarshal(bytes, &resOrg)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if resOrg != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
