package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreateTraceable(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newOrg := Traceable{ObjectType: "org", ID: "org_1", Name: "org 1", Content: "address 1"}
	newOrgAsBytes, err := json.Marshal(newOrg)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateTraceable(t, stub, newOrgAsBytes, newOrg)
}

func checkCreateTraceable(t *testing.T, stub *shim.MockStub, orgAsJSON []byte, value Traceable) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createTraceable"), orgAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	s := string(bytes[:])
	fmt.Println("json: ", s)
	resOrg := Traceable{}
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
