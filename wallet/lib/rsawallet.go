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
	"io/ioutil"
	"os"
	"strconv"
)

// wallet
type RSAWallet struct {
	HomeDir    string
	TypeKey    string
	KeyLen     int
	Ops        string
	Priv       *rsa.PrivateKey
	Pub        *rsa.PublicKey
	To_address string
	Value      int
}

// create
func (w *RSAWallet) Create() error {
	w.Generate()
	w.savePEMKey("private.key")
	w.savePublicPEMKey("publickey.key")
	return nil
}

// verify
func (w *RSAWallet) Verify(publickey *rsa.PublicKey, newhash crypto.Hash, hashed []byte, signature []byte, opts *rsa.PSSOptions) {
	//Verify Signature
	err := rsa.VerifyPSS(publickey, newhash, hashed, signature, opts)

	if err != nil {
		fmt.Println("Who are U? Verify Signature failed")
		os.Exit(1)
	} else {
		fmt.Println("Verify Signature successful")
	}
}
func (w *RSAWallet) GetInfor() {

}
func (w *RSAWallet) saveKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}
func (w *RSAWallet) savePEMKey(fileName string) {
	w.CreateDirectory()
	var path string
	if w.HomeDir != "" {
		path = w.HomeDir + "/"
	} else {
		path = "./"
	}
	outFile, err := os.Create(path + fileName)
	checkError(err)
	defer outFile.Close()
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(w.Priv),
	}
	err = pem.Encode(outFile, privateKey)
	checkError(err)
}
func (w *RSAWallet) savePublicPEMKey(fileName string) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&w.Priv.PublicKey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	w.CreateDirectory()
	var path string
	if w.HomeDir != "" {
		path = w.HomeDir + "/"
	} else {
		path = "./"
	}
	outFile, err := os.Create(path + fileName)
	checkError(err)
	defer outFile.Close()

	err = pem.Encode(outFile, pemkey)
	checkError(err)
}
func (w *RSAWallet) Sign(data string) ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
	message := []byte(data)
	label := []byte("")
	hash := sha256.New()
	publickey := &w.Priv.PublicKey
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

	signature, err := rsa.SignPSS(rand.Reader, w.Priv, newhash, hashed, &opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return signature, ciphertext, label, hash, newhash, hashed, opts
}
func (w *RSAWallet) Transfer() ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&w.Priv.PublicKey)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	data := string(asn1Bytes[:]) + w.To_address + strconv.Itoa(w.Value)
	signature, ciphertext, label, hash, newhash, hashed, opts := w.Sign(data)

	return signature, ciphertext, label, hash, newhash, hashed, opts
}
func (w *RSAWallet) LoadPrivateKey() (*RSAWallet, error) {
	priv, err := ioutil.ReadFile(w.HomeDir + "/private.key")
	if err != nil {
		checkError(err)
	}
	privPem, _ := pem.Decode(priv)
	privPemBytes := privPem.Bytes
	key, err := x509.ParsePKCS1PrivateKey(privPemBytes)
	if err != nil {
		checkError(err)
	}
	w.Priv = key
	return w, nil
}

func (w *RSAWallet) CreateDirectory() {
	if w.HomeDir != "" {
		err := os.MkdirAll(w.HomeDir, 0755)
		if err != nil {
			checkError(err)
			return
		}
	}
}

func (w *RSAWallet) Generate() {
	key, err := rsa.GenerateKey(rand.Reader, w.KeyLen)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	w.Priv = key
	w.Pub = &key.PublicKey
}
