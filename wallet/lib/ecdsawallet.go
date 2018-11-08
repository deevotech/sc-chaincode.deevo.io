package lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type ECDSAWallet struct {
	HomeDir    string
	TypeKey    string
	KeyLen     int
	Ops        string
	Priv       *ecdsa.PrivateKey
	Pub        *ecdsa.PublicKey
	To_address string
	Value      int
}

// create
func (w *ECDSAWallet) Create() error {
	w.Generate()
	w.savePEMKey("private.key")
	w.savePublicPEMKey("publickey.key")
	return nil
}
func (w *ECDSAWallet) EcdsaGenerateKey(c elliptic.Curve) {
	priv, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		checkError(err)
		return
	}
	if !c.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		checkError(err)
		return
	}
	w.Priv = priv
	w.Pub = &priv.PublicKey
}

func (w *ECDSAWallet) Generate() {
	if w.KeyLen == 224 {
		w.EcdsaGenerateKey(elliptic.P224())
	} else if w.KeyLen == 256 {
		w.EcdsaGenerateKey(elliptic.P256())
	} else if w.KeyLen == 348 {
		w.EcdsaGenerateKey(elliptic.P384())
	} else if w.KeyLen == 521 {
		w.EcdsaGenerateKey(elliptic.P521())
	}
}
func (w *ECDSAWallet) CreateDirectory() {
	if w.HomeDir != "" {
		err := os.MkdirAll(w.HomeDir, 0755)
		if err != nil {
			checkError(err)
			return
		}
	}
}

func (w *ECDSAWallet) savePEMKey(fileName string) {
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
	keyDer, err := x509.MarshalECPrivateKey(w.Priv)
	if err != nil {
		log.Fatalf("Failed to serialize ECDSA key: %s\n", err)
	}
	var privateKey = &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyDer,
	}
	err = pem.Encode(outFile, privateKey)
	checkError(err)
}
func (w *ECDSAWallet) savePublicPEMKey(fileName string) {
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
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
func (w *ECDSAWallet) LoadPrivateKey() (*ECDSAWallet, error) {
	priv, err := ioutil.ReadFile(w.HomeDir + "/private.key")
	if err != nil {
		checkError(err)
	}
	privPem, _ := pem.Decode(priv)
	privPemBytes := privPem.Bytes
	key, err := x509.ParseECPrivateKey(privPemBytes)
	if err != nil {
		checkError(err)
	}
	w.Priv = key
	return w, nil
}
