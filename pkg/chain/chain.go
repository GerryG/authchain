
// Chain: These are created with a creation message and an active identity.
// This will create a header record that is then signed by the active identity.
// This signature is the message ID used to index the message in a store.
// Given an existing Chain and an active identity, a new message can be added
// and signed generating a new message ID for the new state of the Chain.
package chain

import (
	//"fmt"
	"errors"
	"math/big"
	"encoding/xml"
)

type chainI interface {
	// interface methods
	IsEmpty() bool
	Previous() *Chain
	Verify() bool
	Author() *Ident
}

type Chain struct {
	eID, prevID, cID ChainID  // Storage indexes of this and related messages
	author *Ident // Set to active identity when signatures are set
	Serialized []byte // Record as it will be stored

	chainI
}

var bigZero big.Int // so we don't need to new this for IsNil each time

// For now a big integer, probably should also have some interfaces
type ChainID big.Int

// New( activeIdentity, value ) Create a new message chain from a generalized object.
// Creates a new Chain with no entries (just the header/creation message that
// will specify what this chain is about, etc. Returns a chain object with no entries.
// Compute message signature for an identity
//func (this *Chain) signature( active *Ident ) ChainID {
func New( active *Ident, m interface{} ) ( chainRoot *Chain, err error ) {
	if m == nil {
		return nil, errors.New("No Create Message to New\n")
	}
	chainRoot = new(Chain)

	// This (*Element) is what our generalized xml parser return, so ...
	el, ok := m.(*Element)
	if ok { // us matching serializer
		chainRoot.Serialized, err = Marshal(el)
		//fmt.Printf("SS1:%s\n", chainRoot.Serialized)
	} else { // this version should be Parsed according the xml structs and tags
		chainRoot.Serialized, err = xml.Marshal(m)
		//fmt.Printf("SS2:%s\n", chainRoot.Serialized)
	}
	if err == nil {
		chainRoot.author = active
		chainRoot.cID = active.Sign( chainRoot.Serialized )
		chainRoot.eID = chainRoot.cID
	//	fmt.Printf("Signed:%s\n", (*big.Int)(&chainRoot.eID).Text(56) )
	//} else {
	//	fmt.Printf("Error marshalling: %s\n", err)
	}
	return
}

// verifies that the ID (eID) is the signature of the data (Serialized)
// true when ID is valid, false otherwise
func (this *Chain) CheckID() bool {
	return this.author.VerifyID( &this.eID, this.Serialized ) == nil
}

// load the header record from its ID
func (this *Chain) getCreate() *Chain {
	if this.cID.IsNil() {
		this.chainZeroID()
	}
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
func (this *Chain) chainZeroID() ChainID {
	if this.cID.IsNil() {
		if this.IsEmpty() {
			this.cID = this.eID
		} else {
			this.cID = this.previous().chainZeroID()
		}
	}
	return this.cID
}

// return true when the chain is empty (Chain is a chain header message)
func (this *Chain) IsEmpty() bool {
	//this.createID()
	//fmt.Printf("IsE:c %s\nIsE:e %s\n", (*big.Int)(&this.eID).Text(56), (*big.Int)(&this.cID).Text(56) )
	return (* big.Int)(&this.cID).Cmp((* big.Int)(&this.eID)) == 0
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
func (this *Chain) Author() *Ident {
	return this.author
}

// AuthChain.get( ChainID ) -> chainEntry
// Lookup and load a specific chain state. If the ID is a chainID,
// the chain is empty, it has no messages added to it.
func Get( eID ChainID ) (c *Chain) {
	c, err := eID.loadMessage()
	if err != nil {
		// log error
	}
	return
}

// addEntry( Chain ) Add message to the chain and returns the resulting chain.
// returns nil when there was a problem (when correctly configured, shouldn't happen)
func (this *Chain) addEntry( m interface{} ) (c *Chain, err error) {
	c = new(Chain)
	c.cID = this.cID
	c.Serialized, err = xml.Marshal(m)
	if err == nil {
		c.prevID = this.eID
		c.author = this.author
		c.eID = c.author.Sign( c.Serialized )
	}
	return
}

/**** ChainID methods ******/

// Load a message from the storage using the index.
func (this ChainID) loadMessage() (m *Chain, e error) {
	return
}

// nil ChainID is represented by a zero value
func (this ChainID) IsNil() bool {
	return (*big.Int)(&this).Cmp(&bigZero) == 0
}

