package lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

type ECDSAWallet struct {
	HomeDir    string
	TypeKey    string
	KeyLen     int
	Ops        string
	Priv       *ecdsa.PrivateKey
	Pub        *ecdsa.PublicKey
	Address    string
	To_address string
	Value      int
}

// create
func (w *ECDSAWallet) Create() error {
	w.Generate()
	w.savePEMKey("private.key")
	w.savePublicPEMKey("publickey.key")
	address, err := w.saveAddress("address.txt")
	if err != nil {
		checkError(err)
	}
	fmt.Println(address)
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
	w.Pub = &w.Priv.PublicKey
	return w, nil
}
func (w *ECDSAWallet) P2PKHToAddress(pkscript []byte, isTestnet bool) (string, error) {
	p := make([]byte, 1)
	p[0] = 0x00 // prefix with 00 if it's mainnet
	if isTestnet {
		p[0] = 0x6F // prefix with 0F if it's testnet
	}
	pub := pkscript[3 : len(pkscript)-2] // get pkhash
	pf := append(p[:], pub[:]...)        // add prefix
	h1 := sha256.Sum256(pf)              // hash it
	h2 := sha256.Sum256(h1[:])           // hash it again
	b := append(pf[:], h2[0:4]...)       // prepend the prefix to the first 5 bytes
	address := base58.Encode(b)          // encode to base58
	if !isTestnet {
		address = "1" + address // prefix with 1 if it's mainnet
	}

	return address, nil
}
func (w *ECDSAWallet) saveAddress(fileName string) (string, error) {

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
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&w.Priv.PublicKey)
	if err != nil {
		log.Fatalf("Failed to serialize ECDSA key: %s\n", err)
	}
	address, err := w.P2PKHToAddress(asn1Bytes, false)
	n, err := outFile.WriteString(address)
	fmt.Printf("wrote %d bytes\n", n)
	if err != nil {
		checkError(err)
	}
	return address, nil
}
func (w *ECDSAWallet) Sign(hash []byte) (r, s *big.Int, err error) {
	zero := big.NewInt(0)

	r, s, err = ecdsa.Sign(rand.Reader, w.Priv, hash)
	if err != nil {
		return zero, zero, err
	}
	return r, s, nil
}
func (w *ECDSAWallet) Verify(hash []byte, r *big.Int, s *big.Int) (result bool) {
	return ecdsa.Verify(&w.Priv.PublicKey, hash, r, s)
}
func (w *ECDSAWallet) Transfer() (r, s *big.Int, d string, err error) {
	current := time.Now()
	data := w.To_address + strconv.Itoa(w.Value) + current.Format("2006-01-02 15:04:05.000000")
	ra, sa, erra := w.Sign([]byte(data))
	return ra, sa, data, erra
}
func (w *ECDSAWallet) Receive(hash []byte, r *big.Int, s *big.Int, d string) (string, bool) {
	check := w.Verify(hash, r, s)
	return d, check
}
