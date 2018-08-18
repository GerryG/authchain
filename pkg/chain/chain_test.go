
package chain

import (
	"github.com/GerryG/authchain/pkg/testing"
)

// Chain.New( activeIdentity, headerMessage ) -> chain:
// Creates a new Chain with no entries (just the header/creation message that
// will specify what this chain is about, etc. Returns a chain object with no entries.
//func New( active *Ident, m *Msg ) *Chain { }

// AuthChain.get( MsgID ) -> chainEntry
// Lookup and load a specific chain state. If the ID is a chainID,
// the chain is empty, it has no messages added to it.
//func (this *Chain) Get( mID MsgID ) *Chain { }

// getID returns the MsgID (authenticating signature) of a chain entry
// Returns ID of the current (last) entry in the chain.
//func (this *Chain) getID() *MsgID { }

// getChainID: Returns ID of the zero entry (header/creation message) of the chain.
//func (this *Chain) getChainID() *MsgID { }

// getMessage -> message: Returns the deserialized message as an object. It will be the header message when the chain is empty
//func (this *Chain) getMessage() *Msg { }

// addEntry( Msg ) Add message to the chain and returns the resulting chain.
//func (this *Chain) addEntry( m *Msg ) *Chain { }

