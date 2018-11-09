package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestAccount_InitAccountdata(t *testing.T) {
	scc := new(AccountChaincode)
	stub := shim.NewMockStub("account", scc)

	checkInit(t, stub, [][]byte{})

	newAccount1 := account{
		Address:   "113yvjFhnmGYN2PaXfD5XT9TDHGbRUyTykiBJ7X3fFG9CMsMCXkr4JksWG2oRy7rpWLkGTM48HhHKLPyDNv8jXoh7jjSYy9zLS9sJw1X2vE2P4Pc66hJtoirwxN8j",
		Publickey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEUkaGIAmlbgE9lfFz2wdMlZSMyTyh\nKnVw7s2wQEgkCA7yrKr8iEXxtGflsBLtqLH7LE071/G3lXn0+tjhlv1Uww==\n-----END PUBLIC KEY-----",
		Balance:   1000,
	}
	newAccount2 := account{
		Address:   "113yvjFhnmGYN2PaXfD5XT9TDHGbRUyTykj3GWqqgK4rJFFrtswbE7xghrX9GRkqVPaYpf4GsSh3jGDeW8MFvubXzAzEEmLbZqvDoueLf8oPv8p5iNEFnsgSA9MeM",
		Publickey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE0pMl4REOfV+19c8+QFLAco5EnM6I\n+kamXYuxYj9fulZidArnsVBD3WoHkSxESuyTpdCGB3YCNxXeaR9wI1gWgg==\n-----END PUBLIC KEY-----",
		Balance:   0,
	}

	/*newAccount1AsBytes, err := json.Marshal(newAccount1)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}*/

	res := stub.MockInvoke("3", [][]byte{[]byte("initAcc"), []byte(newAccount1.Publickey), []byte(newAccount1.Address), []byte("1000")})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	/*newAccount2AsBytes, err := json.Marshal(newAccount2)
	if err != nil {
		fmt.Println("Failed to encode json")
		t.FailNow()
	}*/

	res = stub.MockInvoke("3", [][]byte{[]byte("initAcc"), []byte(newAccount2.Publickey), []byte(newAccount2.Address), []byte("1000")})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	checkAccount(t, stub, newAccount1)
	checkAccount(t, stub, newAccount2)
}

func checkAccount(t *testing.T, stub *shim.MockStub, value account) {
	// Check org
	res := stub.MockInvoke("1", [][]byte{[]byte("getBalance"), []byte(value.Address)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("failed: no payload")
		t.FailNow()
	}

	resData := account{}
	err := json.Unmarshal(res.Payload, &resData)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	fmt.Println(resData.Balance)
	if resData != value {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}

func checkTransfer(t *testing.T, stub *shim.MockStub, hash string, r string, s string, fromAddress string, toAddress string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("transfer"), []byte(value), []byte(hash), []byte(r), []byte(s), []byte(fromAddress), []byte(toAddress)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("failed: no payload")
		t.FailNow()
	}

	resData := account{}
	err := json.Unmarshal(res.Payload, &resData)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	strBalance, err := strconv.ParseFloat(value, 6)
	if resData.Balance != strBalance {
		fmt.Println("Query value was not as expected")
		t.FailNow()
	}
}
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}
