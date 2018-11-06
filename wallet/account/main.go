package main

import (
	"hash"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
	"flag"
	"encoding/gob"
	"encoding/pem"
	"crypto/x509"
	"io/ioutil"
	"strconv"
)

func main() {
	var action string
	var size int
	var pass string
	var value int
	var from string
	var to string
	var privPath string
	flag.StringVar(&action, "action", "create", "a action")
	flag.IntVar(&size, "size", 2048, "a number")
	flag.StringVar(&pass, "pass", "password", "a pass")
	flag.IntVar(&value, "value", 1, "a number")
	flag.StringVar(&from, "from", "xxx", "a address")
	flag.StringVar(&to, "to", "xxx", "a address")
	flag.StringVar(&privPath, "privPath", "./private.key", "file path")
	flag.Parse()
	if action == "create" {
		fmt.Print("create account with ", size, pass)
		createAccount(size)
	} else if action == "transfer" {
		priv, err := ioutil.ReadFile(privPath)
		if err != nil {
			checkError(err)
		}
		privPem, _ := pem.Decode(priv)
		privPemBytes := privPem.Bytes
		key, err := x509.ParsePKCS1PrivateKey(privPemBytes)
		if err != nil {
			checkError(err)
		}
		fmt.Println(key.PublicKey)
		signature, ciphertext, label, hash, newhash, hashed, opts := transfer(from, to, value, key)
		fmt.Println("transfer")
		fmt.Println("singature: ", signature)
		fmt.Println("ciphertext: ", ciphertext)
		fmt.Println("label: ", label)
		fmt.Println("hash: ", hash)
		fmt.Println("newhash: ", newhash)
		fmt.Println("hashed: ", hashed)
		fmt.Println("verify")
		plainText, err := rsa.DecryptOAEP(hash, rand.Reader, key, ciphertext, label)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	
		fmt.Printf("OAEP decrypted [%x] to \n[%s]\n", ciphertext, plainText)
		//Verify Signature
		err = rsa.VerifyPSS(&key.PublicKey, newhash, hashed, signature, &opts)

		if err != nil {
			fmt.Println("Who are U? Verify Signature failed")
			os.Exit(1)
		} else {
			fmt.Println("Verify Signature successful")
		}
	} else if action == "verify" {
		fmt.Print("verify")
	} else {
		fmt.Print("not found action\n")
	}
}
func createAccount(size int) {
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	publicKey := key.PublicKey
	savePEMKey("private.key", key)
	savePublicPEMKey("publickey.key", publicKey)
}
func sign(privateKey *rsa.PrivateKey, data string) ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
	message := []byte(data)
	label := []byte("")
	hash := sha256.New()
	publickey := &privateKey.PublicKey
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publickey, message, label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privateKey, newhash, hashed, &opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return signature, ciphertext, label, hash, newhash, hashed, opts
}
func transfer(from string, to string, value int, privateKey *rsa.PrivateKey) ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
	data := from + to + strconv.Itoa(value)
	signature, ciphertext, label, hash, newhash, hashed, opts := sign(privateKey, data)

	return signature, ciphertext, label, hash, newhash, hashed, opts
}
func verify(publicKey string, data string, singnature string) bool {
	return true
}
func getInfor() {
	fmt.Print("getInfor")
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func saveKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}
func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}
func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	//asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}