
package chain

import (
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
}

type Ident struct {
	Name string
	secDev *SecureServer

	idInterface
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
	var loadid Identity
	err = xml.Unmarshal(data, &loadid)
	if err != nil {
		panic(err)
	}

	id = &Ident{Name: loadid.Name}
	id.secDev = NewSecure( &loadid )
	return
}

func (this *Ident) PublicKey() *rsa.PublicKey {
	pk, ok := this.secDev.Cert.PublicKey.(*rsa.PublicKey)
	if ok { return pk }
	return nil
}

// Use secureDevice to compute signature using key associated with this
func (this *Ident) Sign( message []byte ) (mID ChainID) {
	return this.secDev.sign( message )
}

func (this *Ident) VerifyID( id *ChainID, message []byte ) error {
	hash := sha256.Sum256(message)
	pk := this.PublicKey()
        idBytes := (*big.Int)(id).Bytes()
	return rsa.VerifyPSS( pk, crypto.SHA256, hash[:], idBytes, pssOpts )
}
