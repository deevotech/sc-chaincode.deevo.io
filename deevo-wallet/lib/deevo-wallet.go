package lib

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/redoctober/ecdh"
	"github.com/deevotech/sc-chaincode.deevo.io/deevo-wallet/util"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/bccsp/utils"
)

type DeevoWallet struct {
	DeevoHash []byte
	Priv      *ecdsa.PrivateKey
	Key       bccsp.Key
	PublicKey *ecdsa.PublicKey
	Address   string
	HomeDir   string
	TypeKey   string
	KeyLen    int
	Ops       string
	ToAddress string
	Value     float64
}

func (w *DeevoWallet) Save(path string, priv *ecdsa.PrivateKey) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		checkError(err)
	}
	outFile, err := os.Create(path + "private.key")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	defer outFile.Close()
	keyDer, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		log.Fatalf("Failed to serialize ECDSA key: %s\n", err)
	}
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyDer,
	}
	err = pem.Encode(outFile, privateKey)
	outFile, err = os.Create(path + "public.key")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	defer outFile.Close()
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	err = pem.Encode(outFile, pemkey)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	asn1Bytes, err = x509.MarshalPKIXPublicKey(&priv.PublicKey)
	checkError(err)
	address, err := w.P2PKHToAddress(asn1Bytes, false)
	outFile, err = os.Create(path + "address.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	defer outFile.Close()
	_, err = outFile.WriteString(address)
	checkError(err)
	w.Priv = priv
	w.PublicKey = &priv.PublicKey
	w.Address = address
	return nil
}

// New returns a new instance of the BCCSP implementation
func (w *DeevoWallet) KeyGen(myopts bccsp.KeyGenOpts) (k *DeevoWallet, err error) {

	bkr := csr.NewBasicKeyRequest()
	key, err := bkr.Generate()
	if err != nil {
		fmt.Println(err.Error())
	}
	mykey := key.(*ecdsa.PrivateKey)
	err = w.Save(w.HomeDir, mykey)
	myWallet, err := w.KeyImport(w.HomeDir + "private.key")
	checkError(err)
	w = myWallet
	return w, nil
}

// KeyDeriv derives a key from k using opts.
// The opts argument should be appropriate for the primitive used.
func (w *DeevoWallet) KeyDeriv(k bccsp.Key, opts bccsp.KeyDerivOpts) (dk bccsp.Key, err error) {
	return nil, nil
}

// KeyImport imports a key from its raw representation using opts.
// The opts argument should be appropriate for the primitive used.
func (w *DeevoWallet) KeyImport(fileKey string) (k *DeevoWallet, err error) {
	var myCSP bccsp.BCCSP
	var mspDir = "msp"
	opts := factory.GetDefaultOpts()
	opts.SwOpts.FileKeystore = &factory.FileKeystoreOpts{KeyStorePath: os.TempDir()}
	opts.SwOpts.Ephemeral = false
	myCSP, err = util.InitBCCSP(&opts, "", mspDir)
	myBCCSP, err := util.ImportBCCSPKeyFromPEM(fileKey, myCSP, false)
	checkError(err)
	keyBuff, err := ioutil.ReadFile(fileKey)
	if err != nil {
		return nil, err
	}
	key, err := utils.PEMtoPrivateKey(keyBuff, nil)
	priv := key.(*ecdsa.PrivateKey)
	w.Priv = priv
	w.PublicKey = &priv.PublicKey
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	checkError(err)
	address, err := w.P2PKHToAddress(asn1Bytes, false)
	w.Address = address
	w.Key = myBCCSP
	return w, nil
}

// GetKey returns the key this CSP associates to
// the Subject Key Identifier ski.
func (w *DeevoWallet) GetKey(ski []byte) (k bccsp.Key, err error) {
	return w.Key, nil
}

// Hash hashes messages msg using options opts.
// If opts is nil, the default hash function will be used.
func (w *DeevoWallet) Hash(msg []byte, opts bccsp.HashOpts) (hash []byte, err error) {
	// hash 256
	h := sha256.New()
	h.Write(msg)
	w.DeevoHash = h.Sum(nil)
	return h.Sum(nil), nil
}

// GetHash returns and instance of hash.Hash using options opts.
// If opts is nil, the default hash function will be returned.
func (w *DeevoWallet) GetHash(opts bccsp.HashOpts) (h hash.Hash, err error) {
	hnew := sha256.New()
	return hnew, nil
}

// Sign signs digest using key k.
// The opts argument should be appropriate for the algorithm used.
//
// Note that when a signature of a hash of a larger message is needed,
// the caller is responsible for hashing the larger message and passing
// the hash (as digest).
func (w *DeevoWallet) Sign(k *ecdsa.PrivateKey, digest []byte, opts bccsp.SignerOpts) (signature []byte, err error) {
	r, s, err := ecdsa.Sign(rand.Reader, k, digest)
	if err != nil {
		return nil, err
	}

	s, _, err = utils.ToLowS(&k.PublicKey, s)
	if err != nil {
		return nil, err
	}

	return utils.MarshalECDSASignature(r, s)
}

// Verify verifies signature against key k and digest
// The opts argument should be appropriate for the algorithm used.
func (w *DeevoWallet) Verify(k *ecdsa.PublicKey, signature []byte, digest []byte, opts bccsp.SignerOpts) (valid bool, err error) {
	r, s, err := utils.UnmarshalECDSASignature(signature)
	if err != nil {
		return false, fmt.Errorf("Failed unmashalling signature [%s]", err)
	}

	lowS, err := utils.IsLowS(k, s)
	if err != nil {
		return false, err
	}

	if !lowS {
		return false, fmt.Errorf("Invalid S. Must be smaller than half the order [%s][%s].", s, utils.GetCurveHalfOrdersAt(k.Curve))
	}

	return ecdsa.Verify(k, digest, r, s), nil
}

// Encrypt encrypts plaintext using key k.
// The opts argument should be appropriate for the algorithm used.
func (w *DeevoWallet) Encrypt(k bccsp.Key, plaintext []byte, opts bccsp.EncrypterOpts) (ciphertext []byte, err error) {
	return ecdh.Encrypt(&w.Priv.PublicKey, plaintext)
}

// Decrypt decrypts ciphertext using key k.
// The opts argument should be appropriate for the algorithm used.
func (w *DeevoWallet) Decrypt(k bccsp.Key, ciphertext []byte, opts bccsp.DecrypterOpts) (plaintext []byte, err error) {
	// sha256
	return ecdh.Decrypt(w.Priv, ciphertext)
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
func (w *DeevoWallet) P2PKHToAddress(pkscript []byte, isTestnet bool) (string, error) {
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
func (w *DeevoWallet) Transfer() ([]byte, string, error) {
	current := time.Now()
	value := fmt.Sprintf("%f", w.Value)
	data := w.Address + w.ToAddress + value + current.Format("2006010215:04:05.000000")
	signature, erra := w.Sign(w.Priv, []byte(data), nil)
	return signature, data, erra
}
func (w *DeevoWallet) Receive(publickey *ecdsa.PublicKey, signature []byte, data string) (string, bool) {
	check, err := w.Verify(publickey, signature, []byte(data), nil)
	checkError(err)
	return data, check
}
