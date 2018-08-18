
package chain

import (
	//"math/big"
	//"encoding/xml"
)

type ident interface {
	sign( *Msg ) MsgID
}

type Ident struct {
	Name string
	secDev *SecureDevice

	ident
}
