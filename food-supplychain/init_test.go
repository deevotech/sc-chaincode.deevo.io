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

	newData := InitData{ORG: newOrg, Party1: newParty1, Party2: newParty2, Location1: newLocation1, Location2: newLocation2, Product1: newProduct1, Product2: newProduct2}

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
	originParty1 := value.Party1
	bytes = stub.State[originParty1.ID]
	if bytes == nil {
		fmt.Println("State", originParty1.ID, "failed to get value")
		t.FailNow()
	}
	resParty := Party{}
	err = json.Unmarshal(bytes, &resParty)
	if err != nil {
		fmt.Println("Failed to decode json of Party:", err.Error())
		t.FailNow()
	}
	if resParty != originParty1 {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}

	originParty2 := value.Party2
	bytes = stub.State[originParty2.ID]
	if bytes == nil {
		fmt.Println("State", originParty2.ID, "failed to get value")
		t.FailNow()
	}
	resParty = Party{}
	err = json.Unmarshal(bytes, &resParty)
	if err != nil {
		fmt.Println("Failed to decode json of Party:", err.Error())
		t.FailNow()
	}
	if resParty != originParty2 {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}

	// Check locations
	originLocation1 := value.Location1
	bytes = stub.State[originLocation1.ID]
	if bytes == nil {
		fmt.Println("State", originLocation1.ID, "failed to get value")
		t.FailNow()
	}
	resLocation := Location{}
	err = json.Unmarshal(bytes, &resLocation)
	if err != nil {
		fmt.Println("Failed to decode json of Location:", err.Error())
		t.FailNow()
	}
	if resLocation != originLocation1 {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}

	originLocation2 := value.Location2
	bytes = stub.State[originLocation2.ID]
	if bytes == nil {
		fmt.Println("State", originLocation2.ID, "failed to get value")
		t.FailNow()
	}
	resLocation = Location{}
	err = json.Unmarshal(bytes, &resLocation)
	if err != nil {
		fmt.Println("Failed to decode json of Party:", err.Error())
		t.FailNow()
	}
	if resLocation != originLocation2 {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}

	// Check product
	originProduct1 := value.Product1
	bytes = stub.State[originProduct1.ID]
	if bytes == nil {
		fmt.Println("State", originProduct1.ID, "failed to get value")
		t.FailNow()
	}
	resProduct := Product{}
	err = json.Unmarshal(bytes, &resProduct)
	if err != nil {
		fmt.Println("Failed to decode json of Product:", err.Error())
		t.FailNow()
	}
	if resProduct != originProduct1 {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}

	originProduct2 := value.Product2
	bytes = stub.State[originProduct2.ID]
	if bytes == nil {
		fmt.Println("State", originProduct2.ID, "failed to get value")
		t.FailNow()
	}
	resProduct = Product{}
	err = json.Unmarshal(bytes, &resProduct)
	if err != nil {
		fmt.Println("Failed to decode json of Product:", err.Error())
		t.FailNow()
	}
	if resProduct != originProduct2 {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
