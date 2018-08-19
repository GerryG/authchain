
package chain

import (
	"os"
	"math/big"
	"encoding/xml"
	"crypto"
	"crypto/hmac"
	"crypto/sha256"
)

type devI interface {
	privateKey() crypto.PrivateKey
}

type SecureDevice struct {
	pKey crypto.PrivateKey
	hashMethod crypto.Hash

	devI
}

func NewDevice( fd *os.File, isDev bool ) (sd *SecureDevice) {
	buf := make([]byte, 8192)
 	sd.hashMethod = crypto.SHA256
	var info interface{}
	if isDev {
	} else {
		_, err := fd.Read(buf)
		if err != nil {
			// log error (file too long to read)
		}
		xml.Unmarshal(buf, info)
	}
	//sd.pKey = info.PrivateKey
	return
}

// This needs to be in a secure device or process
func (this *SecureDevice) sign( message []byte ) (id ChainID) {
	if this.hashMethod == crypto.SHA256 {
		mac := hmac.New(sha256.New, this.privateKey().([]byte))
		mac.Write(message)
		(*big.Int)(&id).SetBytes( mac.Sum(nil) )
	} else {
		// log error, hash not implemented
	}
	return
}
