
package chain

import (
	"fmt"
	"math/big"
	"crypto"
	"crypto/tls"
	"crypto/rsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
)

type devI interface {
	privateKey() rsa.PrivateKey
}

// represents a resource that does the signing and holds the private key */
type SecureServer struct {
	XMLName xml.Name `xml:"secure"`
	pKey *rsa.PrivateKey // need to move this to SecureService
	Domain	string   `xml:"domain,attr"`
	Id	string   `xml:"id,attr"`
	PEM	[]byte   `xml:"pem"`
	Cert *x509.Certificate

	devI
}

type Certificate struct {
	XMLName xml.Name `xml:"certificate"`
	Id	string   `xml:"id,attr"`
	Name	string   `xml:"name,attr"`
	PEM	[]byte   `xml:"pem"`
}

func (this *SecureServer) Unlock(cert []byte) (err error) {
	certb, _ := pem.Decode([]byte(cert))
	this.Cert, err = x509.ParseCertificate(certb.Bytes)
	if err != nil { return }
	block, _ := pem.Decode(this.PEM)
	this.pKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	return
}

func NewSecure( loadid *Ident ) (sd *SecureServer) {
	//fmt.Printf("Loaded PEMs from xml: %s\nSec:%s\n", loadid.Certificate.PEM, loadid.Secure.PEM)
	t, err := tls.X509KeyPair( []byte(loadid.Certificate.PEM), []byte(loadid.Secure.PEM) )
	//fmt.Printf("Loaded id from xml: %-v\n", *loadid)
	if err != nil {
		panic(err)
	}
	sd = new(SecureServer)
	sd.Cert, err = x509.ParseCertificate( t.Certificate[0])
	if err != nil {
		panic(err)
	}
	sd.pKey = t.PrivateKey.(*rsa.PrivateKey)
	return
}

var pssOpts = &rsa.PSSOptions{}

// This needs to be in a secure device or process
func (this *SecureServer) sign( message []byte ) (id ChainID) {
        hash := sha256.Sum256(message)
	//fmt.Printf("sign( %#v, %#v\n", this, hash)
        sig, err := rsa.SignPSS(rand.Reader, this.pKey, crypto.SHA256, hash[:], pssOpts)
	if err != nil {
		fmt.Printf("Sig error %s\n", err)
	} else {
		(*big.Int)(&id).SetBytes( sig )
	}
	return
}
