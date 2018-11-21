package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash"
	"log"
	"os"

	"github.com/cloudflare/cfssl/csr"
	"github.com/deevotech/sc-chaincode.deevo.io/deevo-wallet/util"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/bccsp/utils"
)

type DeevoWallet struct {
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
	return nil
}

// New returns a new instance of the BCCSP implementation
func (w *DeevoWallet) KeyGen(myopts bccsp.KeyGenOpts) (k bccsp.Key, err error) {

	bkr := csr.NewBasicKeyRequest()
	key, err := bkr.Generate()
	if err != nil {
		fmt.Println(err.Error())
	}
	mykey := key.(*ecdsa.PrivateKey)
	fmt.Println(mykey.PublicKey)
	err = w.Save("/home/datlv/bccsp-deevo/", mykey)
	myBCCSP, err := w.KeyImport("/home/datlv/bccsp-deevo/private.key")
	checkError(err)
	return myBCCSP, nil
}

// KeyDeriv derives a key from k using opts.
// The opts argument should be appropriate for the primitive used.
func (w *DeevoWallet) KeyDeriv(k bccsp.Key, opts bccsp.KeyDerivOpts) (dk bccsp.Key, err error) {
	return nil, nil
}

// KeyImport imports a key from its raw representation using opts.
// The opts argument should be appropriate for the primitive used.
func (w *DeevoWallet) KeyImport(fileKey string) (k bccsp.Key, err error) {
	var myCSP bccsp.BCCSP
	var mspDir = "msp"
	opts := factory.GetDefaultOpts()
	opts.SwOpts.FileKeystore = &factory.FileKeystoreOpts{KeyStorePath: os.TempDir()}
	opts.SwOpts.Ephemeral = false
	myCSP, err = util.InitBCCSP(&opts, "", mspDir)
	myBCCSP, err := util.ImportBCCSPKeyFromPEM(fileKey, myCSP, false)
	checkError(err)
	return myBCCSP, nil
}

// GetKey returns the key this CSP associates to
// the Subject Key Identifier ski.
func (w *DeevoWallet) GetKey(ski []byte) (k bccsp.Key, err error) {
	return nil, nil
}

// Hash hashes messages msg using options opts.
// If opts is nil, the default hash function will be used.
func (w *DeevoWallet) Hash(msg []byte, opts bccsp.HashOpts) (hash []byte, err error) {
	h :=
		h.Write(msg)
	return h.Sum(nil), nil
}

// GetHash returns and instance of hash.Hash using options opts.
// If opts is nil, the default hash function will be returned.
func (w *DeevoWallet) GetHash(opts bccsp.HashOpts) (h hash.Hash, err error) {
	return nil, nil
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
func (w *DeevoWallet) Verify(k *ecdsa.PublicKey, signature, digest []byte, opts bccsp.SignerOpts) (valid bool, err error) {
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
	return nil, nil
}

// Decrypt decrypts ciphertext using key k.
// The opts argument should be appropriate for the algorithm used.
func (w *DeevoWallet) Decrypt(k bccsp.Key, ciphertext []byte, opts bccsp.DecrypterOpts) (plaintext []byte, err error) {
	return nil, nil
}
func main() {
	var dwallet = &DeevoWallet{}
	dwallet.KeyGen(&bccsp.ECDSAP256KeyGenOpts{})
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
