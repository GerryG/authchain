
package chain

import (
	"os"
)

type idInterface interface {
	sign( []byte ) ChainID
	store() error
}

type Ident struct {
	Name string
	secDev *SecureDevice

	idInterface
}

func NewIdent( sDev *SecureDevice, user interface{} ) (id *Ident) {
	return nil
}

func (this *Ident) store( ) error {
	return nil
}

func LoadIdent( istore interface{} ) (id *Ident) {
	isDevice := false
	fd, ok := istore.(* os.File)
	if !ok {
		fn, ok := istore.(string)
		if ok {
			var err error
			// if fn looks like a URL
			// fd = connect to device/service
			// isDevice = true
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
	id.secDev = NewDevice( fd, isDevice )
	return
}

// Use secureDevice to compute signature using key associated with this
func (this *Ident) sign( message []byte ) (mID ChainID) {
	return this.secDev.sign( message )
}
