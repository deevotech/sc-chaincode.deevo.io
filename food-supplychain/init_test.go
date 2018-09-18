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

	newOrg := Org{ObjectType: "org", ID: "org_1", Name: "org 1", Address: "address 1"}
	newParty1 := Party{ObjectType: "party", ID: "party_1", Name: "party 1", ORG: "org_1"}
	newParty2 := Party{ObjectType: "party", ID: "party_2", Name: "party 2", ORG: "org_1"}
	newLocation1 := Location{ObjectType: "location", ID: "location_1", Name: "location 1", Party: "party_1"}
	newLocation2 := Location{ObjectType: "location", ID: "location_2", Name: "location 2", Party: "party_2"}
	newProduct1 := Product{ObjectType: "product", ID: "product_1", Name: "product 1", Location: "product_1"}
	newProduct2 := Product{ObjectType: "product", ID: "product_2", Name: "product 2", Location: "product_2"}
	newAuditor := Auditor{ObjectType: "Auditor", ID: "Auditor_1", Name: "Auditor 1"}

	newData := InitData{
		ORG:       newOrg,
		Parties:   []Party{newParty1, newParty2},
		Locations: []Location{newLocation1, newLocation2},
		Products:  []Product{newProduct1, newProduct2},
		Auditors:  []Auditor{newAuditor},
	}

	newOrgAsBytes, err := json.Marshal(newData)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkInitORGdata(t, stub, newOrgAsBytes, newData)
}

func checkInitORGdata(t *testing.T, stub *shim.MockStub, orgAsJSON []byte, value InitData) {
	res := stub.MockInvoke("1", [][]byte{[]byte("initOrgData"), orgAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	// Check org
	originOrg := value.ORG
	bytes := stub.State[originOrg.ID]
	if bytes == nil {
		fmt.Println("State", originOrg.ID, "failed to get value")
		t.FailNow()
	}
	resOrg := Org{}
	err := json.Unmarshal(bytes, &resOrg)
	if err != nil {
		fmt.Println("Failed to decode json of ORG:", err.Error())
		t.FailNow()
	}
	if resOrg != originOrg {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}

	// Check parties
	for _, originParty := range value.Parties {
		bytes = stub.State[originParty.ID]
		if bytes == nil {
			fmt.Println("State", originParty.ID, "failed to get value")
			t.FailNow()
		}
		resParty := Party{}
		err = json.Unmarshal(bytes, &resParty)
		if err != nil {
			fmt.Println("Failed to decode json of Party:", err.Error())
			t.FailNow()
		}
		if resParty != originParty {
			fmt.Println("Query value was not as expected")
			t.FailNow()
		}
	}

	// Check locations
	for _, originLocation := range value.Locations {
		bytes = stub.State[originLocation.ID]
		if bytes == nil {
			fmt.Println("State", originLocation.ID, "failed to get value")
			t.FailNow()
		}
		resLocation := Location{}
		err = json.Unmarshal(bytes, &resLocation)
		if err != nil {
			fmt.Println("Failed to decode json of Location:", err.Error())
			t.FailNow()
		}
		if resLocation != originLocation {
			fmt.Println("Query value was not as expected")
			t.FailNow()
		}
	}

	// Check product
	for _, originProduct := range value.Products {
		bytes = stub.State[originProduct.ID]
		if bytes == nil {
			fmt.Println("State", originProduct.ID, "failed to get value")
			t.FailNow()
		}
		resProduct := Product{}
		err = json.Unmarshal(bytes, &resProduct)
		if err != nil {
			fmt.Println("Failed to decode json of Product:", err.Error())
			t.FailNow()
		}
		if resProduct != originProduct {
			fmt.Println("Query value was not as expected")
			t.FailNow()
		}
	}
}
