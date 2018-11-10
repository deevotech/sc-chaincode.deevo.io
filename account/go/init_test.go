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
		Address: "113yvjFhnmGYN2PaXfD5XT9TDHGbRUyTykiBJ7X3fFG9CMsMCXkr4JksWG2oRy7rpWLkGTM48HhHKLPyDNv8jXoh7jjSYy9zLS9sJw1X2vE2P4Pc66hJtoirwxN8j",
		Publickey: `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEUkaGIAmlbgE9lfFz2wdMlZSMyTyh
KnVw7s2wQEgkCA7yrKr8iEXxtGflsBLtqLH7LE071/G3lXn0+tjhlv1Uww==
-----END PUBLIC KEY-----
`,
		Balance: 1000,
	}
	newAccount2 := account{
		Address: "113yvjFhnmGYN2PaXfD5XT9TDHGbRUyTykj3GWqqgK4rJFFrtswbE7xghrX9GRkqVPaYpf4GsSh3jGDeW8MFvubXzAzEEmLbZqvDoueLf8oPv8p5iNEFnsgSA9MeM",
		Publickey: `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE0pMl4REOfV+19c8+QFLAco5EnM6I
+kamXYuxYj9fulZidArnsVBD3WoHkSxESuyTpdCGB3YCNxXeaR9wI1gWgg==
-----END PUBLIC KEY-----
`,
		Balance: 2000,
	}

	newAccount3 := account{
		Address: "113yvjFhnmGYN2PaXfD5XT9TDHGbRUyTykiQNQhvigEvs4EQyKKDGp9dz8gxwejg5ADNut35369oriF1wEEqHXZkYKzRvFCcr6RRi5r5Tqdze4uYi6PgohrTzyENi",
		Publickey: `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEc9eFHSThwRVyy6cUr1OOBAeRTHLY
Puqrx+WvavUNw8TMS09IBkh9pyNwjiWyyynNisF+9wqqgyzfi56aGyxDRA==
-----END PUBLIC KEY-----
`,
		Balance: 500,
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

	res = stub.MockInvoke("3", [][]byte{[]byte("initAcc"), []byte(newAccount2.Publickey), []byte(newAccount2.Address), []byte("2000")})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}
	res = stub.MockInvoke("3", [][]byte{[]byte("initAcc"), []byte(newAccount3.Publickey), []byte(newAccount3.Address), []byte("500")})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}

	checkAccount(t, stub, newAccount1.Address, "1000")
	checkAccount(t, stub, newAccount2.Address, "2000")
	checkAccount(t, stub, newAccount3.Address, "500")

	r1 := "29978764154139864880696030835938256737287232151003611159021314956225371873730"
	s1 := "75348994682775390770944787851125569805092606536265710298111367961967701172281"
	data1 := "113yvjFhnmGYN2PaXfD5XT9TDHGbRUyTykj3GWqqgK4rJFFrtswbE7xghrX9GRkqVPaYpf4GsSh3jGDeW8MFvubXzAzEEmLbZqvDoueLf8oPv8p5iNEFnsgSA9MeM1002018-11-10 12:17:29.014634"
	checkTransfer(t, stub, data1, r1, s1, newAccount1.Address, newAccount2.Address, "100")
	checkAccount(t, stub, newAccount1.Address, "900")
	checkAccount(t, stub, newAccount2.Address, "2100")
}

func checkAccount(t *testing.T, stub *shim.MockStub, address string, value string) {
	// Check org
	res := stub.MockInvoke("1", [][]byte{[]byte("getBalance"), []byte(address)})
	if res.Status != shim.OK {
		fmt.Println("failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("failed: no payload account")
		t.FailNow()
	}

	resData := account{}
	err := json.Unmarshal(res.Payload, &resData)
	if err != nil {
		fmt.Println("Failed to decode json:", err.Error())
		t.FailNow()
	}
	fmt.Println(resData.Balance)
	//value1 := fmt.Sprintf("%f", value)
	value2 := strconv.FormatFloat(resData.Balance, 'f', -1, 64)
	if value != value2 {
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
}
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}
