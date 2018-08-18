
package chain

import (
	//"math/big"
	//"encoding/xml"
)

type ident interface {
	sign() MsgID
}

type Ident struct {
	Name string
	secDev *SecureDevice

	ident
}
