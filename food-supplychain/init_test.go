package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_InitORGdata(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newOrg := Traceable{
		ObjectType: "org",
		ID:         "org_1",
		Name:       "org 1",
		Content:    "address 1",
	}
	newParty1 := Traceable{
		ObjectType: "party",
		ID:         "party_1",
		Name:       "party 1",
		Parent:     "org_1",
	}
	newParty2 := Traceable{
		ObjectType: "party",
		ID:         "party_2",
		Name:       "party 2",
		Parent:     "org_1",
	}
	newLocation1 := Traceable{
		ObjectType: "location",
		ID:         "location_1",
		Name:       "location 1",
		Parent:     "party_1",
	}
	newLocation2 := Traceable{
		ObjectType: "location",
		ID:         "location_2",
		Name:       "location 2",
		Parent:     "party_2",
	}
	newProduct1 := Traceable{
		ObjectType: "product",
		ID:         "product_1",
		Name:       "product 1",
		Parent:     "product_1",
	}
	newProduct2 := Traceable{
		ObjectType: "product",
		ID:         "product_2",
		Name:       "product 2",
		Parent:     "product_2",
	}
	newAuditor := Auditor{
		ObjectType: "auditor",
		ID:         "Auditor_1",
		Name:       "Auditor 1",
	}

	newData := InitData{
		Traceable: []Traceable{newOrg, newParty1, newParty2, newLocation1, newLocation2, newProduct1, newProduct2},
		Auditors:  []Auditor{newAuditor},
	}

	newOrgAsBytes, err := json.Marshal(newData)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}

	res := stub.MockInvoke("1", [][]byte{[]byte("initOrgData"), newOrgAsBytes})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	for _, data := range newData.Traceable {
		checkTraceableData(t, stub, data)
	}
	for _, a := range newData.Auditors {
		checkInitAuditor(t, stub, a)
	}
}

func checkInitAuditor(t *testing.T, stub *shim.MockStub, value Auditor) {
	// Check org
	res := stub.MockInvoke("1", [][]byte{[]byte("getObject"), []byte(value.ID), []byte(TYPE_AUDITOR)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("failed: no payload")
		t.FailNow()
	}

	resData := Auditor{}
	err := json.Unmarshal(res.Payload, &resData)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if resData != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
