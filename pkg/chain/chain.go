
// Chain: These are created with a creation message and an active identity.
// This will create a header record that is then signed by the active identity.
// This signature is the message ID used to index the message in a store.
// Given an existing Chain and an active identity, a new message can be added
// and signed generating a new message ID for the new state of the Chain.
package chain

import (
	"math/big"
	"encoding/xml"
)

type chainI interface {
	// interface methods
	createMessage() *Chain
	createID() ChainID
	author() *Ident  // Identity of the chain's authorized identity
	isEmpty() bool
	previous() *Chain
}

type Chain struct {
	mID, prevID, cID ChainID  // Storage indexes of this and related messages
	auth *Ident // Set to active identity when signatures are set
	Serialized []byte // Record as it will be stored

	chainI
}

var bigZero big.Int // so we don't need to new this for isNil each time

// For now a big integer, probably should also have some interfaces
type ChainID big.Int

// Load a message from the storage using the index.
func (this ChainID) loadMessage() (m *Chain, e error) {
	return
}

// nil ChainID is represented by a zero value
func (this ChainID) isNil() bool {
	return (*big.Int)(&this).Cmp(&bigZero) == 0
}

// New( activeIdentity, value ) Create a new message chain from a generalized object.
// Creates a new Chain with no entries (just the header/creation message that
// will specify what this chain is about, etc. Returns a chain object with no entries.
// Compute message signature for an identity
//func (this *Chain) signature( active *Ident ) ChainID {
func New( active *Ident, m interface{} ) ( chainRoot Chain, err error ) {
	chainRoot.Serialized, err = xml.Marshal(m)
	if err == nil {
		chainRoot.auth = active
		chainRoot.mID = active.sign( chainRoot.Serialized )
	}
	return
}

// load the header record from its ID
func (this *Chain) createMessage() *Chain {
	msg, err := this.cID.loadMessage()
	if err == nil {
		return msg
	} else {
		// log error loading chain header
		return nil
	}
}

// The message ID of the chain header (zero entry). Follow the chain back to the beginning
// if it isn't already set on the object.
func (this *Chain) createID() ChainID {
	if this.cID.isNil() {
		if this.isEmpty() {
			this.cID = this.mID
		} else {
			this.cID = this.previous().createID()
		}
	}
	return this.cID
}

// return true when the chain is empty (Chain is a chain header message)
func (this *Chain) isEmpty() bool {
	this.createID()
	return (* big.Int)(&this.cID).Cmp((* big.Int)(&this.mID)) == 0
}

// Get the previous entry in the chain
func (this *Chain) previous() *Chain {
	prev, err := this.prevID.loadMessage()
	if err == nil {
		return prev
	} else {
		// log error loading previous in chain
		return nil
	}
}

// Identity of the chain's authorized identity
func (this *Chain) author() *Ident {
	return this.author()
}

// AuthChain.get( ChainID ) -> chainEntry
// Lookup and load a specific chain state. If the ID is a chainID,
// the chain is empty, it has no messages added to it.
func Get( mID ChainID ) *Chain {
	c, err := mID.loadMessage()
	if err != nil {
		// log error
	}
	return c
}

// getID returns the ChainID (authenticating signature) of a chain entry
// Returns ID of the current (last) entry in the chain.
func (this *Chain) getID() (id ChainID) {
	return
}

// getChainID: Returns ID of the zero entry (header/creation message) of the chain.
func (this *Chain) getChainID() (id ChainID) {
	return
}

// getMessage -> message: Returns the deserialized message as an object. It will be the header message when the chain is empty
func (this *Chain) getMessage() (m *Chain) {
	return
}

// addEntry( Chain ) Add message to the chain and returns the resulting chain.
func (this *Chain) addEntry( m *Chain ) (c *Chain) {
	return
}

