package main

import (
	"fmt"
	"testing"
)

func TestWallet_ShareData(t *testing.T) {
	myWallet := &WalletCmd{
		homeDirectory: "/home/datlv/wallet",
		length:        256,
		typeKey:       "ecdsa",
	}
	ecdsaWallet, err := myWallet.getECDSAWallet().LoadPrivateKey()

	if err != nil {
		fmt.Println("Error read file")
		t.FailNow()
	}
	data := "Helloworld!"
	dataEncrypt, err := ecdsaWallet.Encrypt([]byte(data))
	if err != nil {
		fmt.Println("Error encrypt")
		t.FailNow()
	}
	dataDecrypt, err := ecdsaWallet.Decrypt([]byte(dataEncrypt))
	if err != nil {
		fmt.Println("Error Decrypt")
		t.FailNow()
	}
	if data != string(dataDecrypt) {
		fmt.Println("Decrypt not the same")
		t.FailNow()
	}
}
