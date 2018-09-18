package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestFood_CreareParty(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newParty := Party{ObjectType: "Party", ID: "Party_1", Name: "Party 1", ORG: "org_1"}
	newPartyAsBytes, err := json.Marshal(newParty)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateParty(t, stub, newPartyAsBytes, newParty)
}

func TestFood_UpdateParty(t *testing.T) {
	scc := new(FoodChaincode)
	stub := shim.NewMockStub("food", scc)

	checkInit(t, stub, [][]byte{})

	newParty := Party{ObjectType: "Party", ID: "Party_1", Name: "Party 1", ORG: "org_1"}
	newPartyAsBytes, err := json.Marshal(newParty)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkCreateParty(t, stub, newPartyAsBytes, newParty)

	updatedParty := Party{ObjectType: "Party", ID: "Party_1", Name: "Party 2", ORG: "org_1"}
	updatedPartyAsBytes, err := json.Marshal(updatedParty)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}
	checkUpdateParty(t, stub, updatedPartyAsBytes, updatedParty)
}
func checkCreateParty(t *testing.T, stub *shim.MockStub, partyAsJSON []byte, value Party) {
	res := stub.MockInvoke("1", [][]byte{[]byte("createParty"), partyAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resParty := Party{}
	err := json.Unmarshal(bytes, &resParty)
	if err != nil {
		fmt.Println("Failed to decode json of Party:", err.Error())
		t.FailNow()
	}
	if resParty != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkUpdateParty(t *testing.T, stub *shim.MockStub, partyAsJSON []byte, value Party) {
	res := stub.MockInvoke("1", [][]byte{[]byte("updateParty"), partyAsJSON})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	bytes := stub.State[value.ID]
	if bytes == nil {
		fmt.Println("State", value.ID, "failed to get value")
		t.FailNow()
	}
	resParty := Party{}
	err := json.Unmarshal(bytes, &resParty)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	if resParty != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
