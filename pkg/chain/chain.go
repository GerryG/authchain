
// Chain: These are created with a creation message and an active identity.
// This will create a header record that is then signed by the active identity.
// This signature is the message ID used to index the message in a store.
// Given an existing Chain and an active identity, a new message can be added
// and signed generating a new message ID for the new state of the Chain.
package chain

import (
	"fmt"
	"time"
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
	at time.Time

	chainI
}

// xml for the chain header message
type ChainHeader struct {
	XMLName xml.Name `xml:"header"`
	Author string    `xml:"author"`
	At time.Time     `xml:"at"`
	App *Element     `xml:"app"`

	xml.Unmarshaler
}

// xml for the chain entry message
type ChainEntry struct {
	XMLName xml.Name `xml:"entry"`
	Entry *Element

	xml.Unmarshaler
}

var bigZero big.Int // so we don't need to new this for IsNil each time

// For now a big integer, probably should also have some interfaces
type ChainID big.Int

// New( activeIdentity, value ) Create a new message chain from a generalized object.
// Creates a new Chain with no entries (just the header/creation message that
// will specify what this chain is about, etc. Returns a chain object with no entries.
// Compute message signature for an identity
//func (this *Chain) signature( active *Ident ) ChainID {
func New( active *Ident, m *Element ) ( chainRoot *Chain, err error ) {
	if m == nil {
		return nil, errors.New("No Create Message to New\n")
	}
	chainRoot = new(Chain)
	chainRoot.author = active

	mHead := ChainHeader{
		Author: active.Domain()+":"+active.Id(),
		At: time.Now(),
		App: m,
	}
	chainRoot.Serialized, err = xml.Marshal(&mHead)
	if err == nil {
		chainRoot.cID = active.Sign( chainRoot.Serialized )
		chainRoot.eID = chainRoot.cID
	} else {
		fmt.Printf("Error marshalling: %s\n", err)
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
			this.cID = this.Previous().chainZeroID()
		}
	}
	return this.cID
}

// return true when the chain is empty (Chain is a chain header message)
func (this *Chain) IsEmpty() bool {
	return (* big.Int)(&this.cID).Cmp((* big.Int)(&this.eID)) == 0
}

// Get the previous entry in the chain
func (this *Chain) Previous() *Chain {
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
func (this *Chain) AddEntry( m *Element ) (chain *Chain, err error) {
	chain = new(Chain)

	chain.Serialized, err = Marshal("entry", m)
	if err == nil {
		chain.author = this.author
		chain.cID = this.cID
		chain.prevID = this.eID
		chain.author = this.author
		chain.eID = chain.author.Sign( chain.Serialized )
		//fmt.Printf("AddE:%s\nNewI:%s\n", (*big.Int)(&chain.cID).Text(62), (*big.Int)(&chain.eID).Text(62))
	} else {
		fmt.Printf("Error from marshal: %s\n", err)
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

