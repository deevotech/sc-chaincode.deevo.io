package main

import (
	"io/ioutil"
	//"bytes"
	//"encoding/json"
	"fmt"
	//"strconv"
	//"time"

	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	//"github.com/hyperledger/fabric/core/chaincode/shim"
	//pb "github.com/hyperledger/fabric/protos/peer"
	"crypto/sha256"
	//"crypto/rsa"
	//"crypto/rand"
	"crypto/x509"
	//"crypto/x509/pkix"
	//"crypto/x509/pkcs1"
	"encoding/pem"
	"math/big"
	//"time"
)

// interface permission
// ===========================
type PermissionChaincode interface {
	SetOwner(pub string, ca string, signature string) bool
	HashCa(ca string) bool
	DecryptSignature(signatrue string, ca string, publickey string)
	ChangeOwner() bool
	GetObject()
	ShareObject() bool
	UpdateObject() bool
	DeleteObject() bool
	ViewObject()
	IsOwner() bool
	HaveInShare() bool
}
// structure of permission
// ===========================
type PermissionStruct struct {
	Publickey string
	Ca string
}
type PermissionStructObject struct {
	Owner PermissionStruct
	Viewer []PermissionStruct
	Name string
	Id string
}
func (p *PermissionStruct) SetOwner(pub string, ca string, signature string) bool {
	// check publickey, ca, signature
	if len(pub) <= 0 || len(ca) <= 0 || len(signature) <=0 {
		return false
	}
	// tinh hash ca
	//var hash := h
	p.Publickey = pub
	p.Ca = ca
	return true
}
func (p *PermissionStruct) ChangeOwner(pub string, ca string, signature string) bool {
	// check publickey, ca, signature
	if len(pub) <= 0 || len(ca) <= 0 || len(signature) <=0 {
		return false
	}
	// tinh hash ca
	//var hash := h
	p.Publickey = pub
	p.Ca = ca
	return true
}
func (p *PermissionStructObject) IsOwner(pub string, ca string, signature string) bool {
	// check publickey, ca, signature
	if len(pub) <= 0 || len(ca) <= 0 || len(signature) <=0 {
		return false
	}

	return true
}
func CheckShare(check PermissionStruct, owner PermissionStruct, signature string) bool {
	if len(signature) <=0 {
		return false
	}
	return true
}
func (p PermissionStructObject) HaveInShare(owner PermissionStruct, signature string) bool {
	// check publickey, ca, signature
	if len(signature) <=0 {
		return false
	}
	for i:=0; i<len(p.Viewer); i++ {
		if CheckShare(p.Viewer[0], owner, signature) {
			return true
		}
	}
	return false
}
func (p PermissionStructObject) GetObject(owner PermissionStruct, signature string) (string, string) {
	// check publickey, ca, signature
	if p.IsOwner(owner.Publickey, owner.Ca, signature) == false {
		return "", ""
	}
	return p.Name, p.Id
}
func (p *PermissionStructObject) ShareObject(pub string, ca string, signature string) bool {
	// check publickey, ca, signature
	if len(pub) <= 0 || len(ca) <= 0 || len(signature) <=0 {
		return false
	}
	var newP PermissionStruct
	newP.Ca = ca
	newP.Publickey = pub
	p.Viewer = append(p.Viewer, newP)
	return true
}
func (p *PermissionStructObject) UpdateObject(pub string, ca string, signature string, name string, id string) bool {
	// check publickey, ca, signature
	if len(pub) <= 0 || len(ca) <= 0 || len(signature) <=0 {
		return false
	}
	p.Id = id
	p.Name = name
	return true
}
func (p *PermissionStructObject) DeleteObject(pub string, ca string, signature string) bool {
	// check publickey, ca, signature
	if len(pub) <= 0 || len(ca) <= 0 || len(signature) <=0 {
		return false
	}
	return true
}
func main() {
	var obj1 PermissionStructObject
	var obj2 PermissionStructObject
	var account1 PermissionStruct
	var account2 PermissionStruct
	b, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/signcerts/cert.pem")
	if err != nil {
		fmt.Print(err)
	}
	c, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/tlscacert/tls-rca-org2-deevo-com-7054.pem")
	e, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/keystore/key.pem")
	f, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/keystore/pub.pem")
	account1.Publickey = string(b)
	account1.Ca = string(c)
	b1, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/signcerts/cert.pem")
	c1, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/tlscacert/tls-rca-org2-deevo-com-7054.pem")
	e1, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/keystore/key.pem")
	f1, err := ioutil.ReadFile("config/tls-peer0.org2.deevo.com/keystore/pub.pem")
	account2.Publickey = string(b1)
	account2.Ca = string(c1)
	var signature1, signature2, priv1, priv2 string
	h := sha256.New()
	h.Write([]byte(b))
	fmt.Printf("%d", h.Sum(nil))
	signature1 = "signature1"
	signature2 = "signature2"
	priv1 = string(e)
	priv2 = string(e1)
	obj1.Owner.SetOwner(account1.Publickey, account1.Ca, signature1)
	obj2.Owner.SetOwner(account2.Publickey, account2.Ca, signature2)
	fmt.Printf("%s\n, %s\n, %s\n", account1.Publickey, f, priv1)
	fmt.Printf("%s\n, %s\n, %s\n", account2.Publickey, f1, priv2)
	block, _ := pem.Decode([]byte(e))
	rsaPriv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	
	if err != nil {
		panic("Failed to parse private key: " + err.Error())
	}
	/*template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "test",
			Organization: []string{"Σ Acme Co"},
		},
		NotBefore: time.Unix(1000, 0),
		NotAfter:  time.Unix(100000, 0),
		KeyUsage:  x509.KeyUsageCertSign,
	}*/
	fmt.Printf("%+v\n", rsaPriv)
	block1, _ := pem.Decode([]byte(b))
	// cert la file in cert.pem
	cert, err := x509.ParseCertificate(block1.Bytes)
	fmt.Printf("%+v\n", cert)
	fmt.Printf("%+v\n", cert.Signature)
	fmt.Printf("%+v\n", cert.PublicKey)
	/*if _, err = x509.CreateCertificateRequest(rand.Reader, &template, &template, rsaPriv) ; err != nil {
		panic("failed to create certificate with basic imports: " + err.Error())
	}*/
}
type pkcs1PublicKey struct {
	N *big.Int
	E int
}
