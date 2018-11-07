package lib

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"hash"
	"os"
	"strconv"
)

// wallet
type Wallet struct {
	HomeDir string
	TypeKey string
	Length  int
	Ops     string
}

// create
func (w *Wallet) Create() error {
	key, err := rsa.GenerateKey(rand.Reader, w.Length)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	publicKey := key.PublicKey
	w.savePEMKey("private.key", key)
	w.savePublicPEMKey("publickey.key", publicKey)
	return nil
}

// verify
func (w *Wallet) Verify(publickey *rsa.PublicKey, newhash crypto.Hash, hashed []byte, signature []byte, opts *rsa.PSSOptions) {
	//Verify Signature
	err := rsa.VerifyPSS(publickey, newhash, hashed, signature, opts)

	if err != nil {
		fmt.Println("Who are U? Verify Signature failed")
		os.Exit(1)
	} else {
		fmt.Println("Verify Signature successful")
	}
}
func (w *Wallet) GetInfor() {

}
func (w *Wallet) saveKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}
func (w *Wallet) savePEMKey(fileName string, key *rsa.PrivateKey) {
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
func (w *Wallet) savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
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
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
func (w *Wallet) Sign(privateKey *rsa.PrivateKey, data string) ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
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
func (w *Wallet) Transfer(from string, to string, value int, privateKey *rsa.PrivateKey) ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
	data := from + to + strconv.Itoa(value)
	signature, ciphertext, label, hash, newhash, hashed, opts := w.Sign(privateKey, data)

	return signature, ciphertext, label, hash, newhash, hashed, opts
}
