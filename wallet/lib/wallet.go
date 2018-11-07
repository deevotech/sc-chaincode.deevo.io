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
type Wallet struct {
	HomeDir    string
	TypeKey    string
	Length     int
	Ops        string
	privateKey *rsa.PrivateKey
	To_address string
	Value      int
}

// create
func (w *Wallet) Create() error {
	key, err := rsa.GenerateKey(rand.Reader, w.Length)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	publicKey := key.PublicKey
	var path string
	if w.HomeDir != "" {
		path = w.HomeDir + "/"
		err = os.MkdirAll(w.HomeDir, 0755)
		if err != nil {
			return err
		}
	} else {
		path = "./"
	}
	w.savePEMKey(path+"private.key", key)
	w.savePublicPEMKey(path+"publickey.key", publicKey)
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
func (w *Wallet) Sign(data string) ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
	message := []byte(data)
	label := []byte("")
	hash := sha256.New()
	publickey := &w.privateKey.PublicKey
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

	signature, err := rsa.SignPSS(rand.Reader, w.privateKey, newhash, hashed, &opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return signature, ciphertext, label, hash, newhash, hashed, opts
}
func (w *Wallet) Transfer() ([]byte, []byte, []byte, hash.Hash, crypto.Hash, []byte, rsa.PSSOptions) {
	/*asn1Bytes, err := x509.MarshalPKIXPublicKey(&w.privateKey.PublicKey)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	data := string(asn1Bytes[:]) + w.To_address + strconv.Itoa(w.Value)*/
	fmt.Println(w.privateKey.PublicKey)
	data := "xxx" + w.To_address + strconv.Itoa(w.Value)
	signature, ciphertext, label, hash, newhash, hashed, opts := w.Sign(data)

	return signature, ciphertext, label, hash, newhash, hashed, opts
}
func (w *Wallet) LoadPrivateKey() error {
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
	w.privateKey = key
	return nil
}
