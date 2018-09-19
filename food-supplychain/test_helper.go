package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkTraceableData(t *testing.T, stub *shim.MockStub, value Traceable) {
	// Check org
	res := stub.MockInvoke("1", [][]byte{[]byte("getObject"), []byte(value.ID), []byte(value.ObjectType)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("failed: no payload")
		t.FailNow()
	}

	resData := Traceable{}
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
