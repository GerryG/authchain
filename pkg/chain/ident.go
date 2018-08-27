
package chain

import (
	//"fmt"
	"os"
	"math/big"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"io/ioutil"
	"encoding/xml"
)

type idInterface interface {
	Sign( []byte ) ChainID
	Store() error
	VerifyID( *ChainID, []byte ) error
	PublicKey() *rsa.PublicKey
	Id() string
	Domain() string
}

// for identity file
type Ident struct {
	XMLName xml.Name `xml:"identity"`
	Name	string   `xml:"name"`
	Certificate Certificate
	Secure	SecureServer `xml:"secure"`

	idInterface
}


func (this *Ident) Id() string {
	return this.Secure.Id
}
func (this *Ident) Domain() string {
	return this.Secure.Domain
}
/*
func NewIdent( sDev *SecureServer, user interface{} ) (id *Ident) {
	return nil
}

func (this *Ident) store( ) error {
	return nil
}
*/

func LoadIdent( istore interface{}, ext string ) (id *Ident) {
	//IsDevice := false
	fd, ok := istore.(* os.File)
	if !ok {
		fn, ok := istore.(string)
		if ok {
			var err error
			// if fn looks like a URL
			// fd = connect to device/service
			// IsDevice = true
			// or ...
			fd, err = os.Open( fn )
			if err != nil {
				// log error (panic?)
				return nil
			}
		}
	}
	if fd == nil {
		return nil
	}

	// read and decode XML file
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	var loadid Ident
	err = xml.Unmarshal(data, &loadid)
	//fmt.Printf("Ident: %#v\n", loadid)
	if err != nil {
		panic(err)
	}

	return &loadid
}

func (this *Ident) PublicKey() *rsa.PublicKey {
	pk, ok := this.Secure.Cert.PublicKey.(*rsa.PublicKey)
	if ok { return pk }
	return nil
}

// Use secureDevice to compute signature using key associated with this
func (this *Ident) Sign( message []byte ) (mID ChainID) {
	return this.Secure.sign( message )
}

func (this *Ident) VerifyID( id *ChainID, message []byte ) error {
	hash := sha256.Sum256(message)
	pk := this.PublicKey()
        idBytes := (*big.Int)(id).Bytes()
	return rsa.VerifyPSS( pk, crypto.SHA256, hash[:], idBytes, pssOpts )
}
