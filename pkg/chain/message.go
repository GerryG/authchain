
package chain

import (
	"math/big"
	"encoding/xml"
)

type msgI interface {
	// interface methods
	createMessage() *Msg
	createID() MsgID
	author() *Ident  // Identity of the chain's authorized identity
	signature( active *Ident ) MsgID
	isEmpty() bool
	previous() *Msg
}

type Msg struct {
	mID, prevID, cID MsgID  // Storage indexes of this and related messages
	auth *Ident // Set to active identity when signatures are set
	Serialized []byte // Record as it will be stored

	msgI
}

var bigZero big.Int // so we don't need to new this for isNil each time

// For now a big integer, probably should also have some interfaces
type MsgID big.Int

// Load a message from the storage using the index.
func (this MsgID) loadMessage() (m *Msg, e error) {
	return
}

// nil MsgID is represented by a zero value
func (this MsgID) isNil() bool {
	return (*big.Int)(&this).Cmp(&bigZero) == 0
}

// NewMessage( value ) Create a new message from a generalized object.
func NewMessage( v interface{} ) (msg Msg, err error ) {
	serialValue, err := xml.Marshal(v)
	if err == nil {
		msg.Serialized = serialValue
	}
	return
}

func (this *Msg) createMessage() *Msg {
	msg, err := this.cID.loadMessage()
	if err == nil {
		return msg
	} else {
		// log error loading chain header
		return nil
	}
}

// The message ID of the chain header (zero entry)
func (this *Msg) createID() MsgID {
	if this.cID.isNil() {
		if this.isEmpty() {
			this.cID = this.mID
		} else {
			this.cID = this.previous().createID()
		}
	}
	return this.cID
}

// return true when the chain is empty (Msg is a chain header message)
func (this *Msg) isEmpty() bool {
	return false
}

// Get the previous entry in the chain
func (this *Msg) previous() *Msg {
	prev, err := this.prevID.loadMessage()
	if err == nil {
		return prev
	} else {
		// log error loading previous in chain
		return nil
	}
}

// Identity of the chain's authorized identity
func (this *Msg) author() *Ident {
	return this.author()
}

// Compute message signature for an identity
func (this *Msg) signature( active *Ident ) MsgID {
	if this.mID.isNil() {
		this.auth = active
		this.mID = active.sign( this )
	}

	return this.mID
}
