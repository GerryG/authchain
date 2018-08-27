
package chain_test

import ( 
	//"fmt"
	"bytes"
	. "github.com/GerryG/authchain/pkg/chain"
	"testing"
	"math/big"
	"encoding/xml"
)

// Load a message from the storage using the index.
//func (this ChainID) loadMessage() (m *Chain, e error) {
func TestLoadMessage(t *testing.T) {
}

// nil ChainID is represented by a zero value
func TestIsNil(t *testing.T) {
	var testChain ChainID
	if !testChain.IsNil() {
		t.Errorf("IsNil reports false for a nil")
	}
	(*big.Int)(&testChain).Add((*big.Int)(&testChain), new(big.Int).SetInt64(1))
	if testChain.IsNil() {
		t.Errorf("IsNil reports true for a 1")
	}
	(*big.Int)(&testChain).Sub((*big.Int)(&testChain), new(big.Int).SetInt64(1000))
	if testChain.IsNil() {
		t.Errorf("IsNil reports true for a negative")
	}
}

//func New( active *Ident, m interface{} ) ( chainRoot *Chain, err error ) {
func TestNewChainVerify( t *testing.T ) {
	tname := "chainnew"
	aID := testdataIdent(t, tname)
	if aID == nil {
		return
	}
	// this load is using *Element method UnmarshalXML to decode the file input
	createMessage := testdataCreate(t, tname)

	if createMessage == nil {
		panic("No message from testdataCreate\n")
	}

	var cRoot *Chain
	var err error
	t.Run("New", func(t *testing.T) {
		cRoot, err = New( aID, createMessage )
		if err != nil {
			t.Errorf( "Unexpected error from New: %s\n", err )
		}
		if cRoot == nil {
			t.Errorf( "No chain object returned\n" )
		}
	})
	t.Run("StartEmpty", func(t *testing.T) {
		if !cRoot.IsEmpty() {
			t.Errorf( "New chain should be empty\n" )
		}
	})
	var head ChainHeader
	t.Run("CheckMarshalling", func(t *testing.T) {
		err := Unmarshal("header", bytes.NewBuffer(cRoot.Serialized), xml.Unmarshaler(&head))
		if err != nil {
			t.Errorf("Error from Unmarshal: %s\n", err)
		}
	})
	t.Run("CheckRemarshal", func(t *testing.T) {
		serAgain, err := xml.Marshal(head)
		if err != nil {
			t.Errorf("Error from (re)Marshal: %s\n", err)
		}
		if string(serAgain) != string(cRoot.Serialized) {
			t.Errorf("Differ:\nAgain:%s\nB orig:%s\n", serAgain, cRoot.Serialized )
		}
	})
	t.Run("VerifyID", func(t *testing.T) {
		if !cRoot.CheckID() {
			t.Error("Signature does not match\n")
		}
	})
}

// load the header record from its ID
//func (this *Chain) getCreate() *Chain {
func TestGetCreate( t *testing.T ) {
}

// The message ID of the chain header (zero entry). Follow the chain back to the beginning
// if it isn't already set on the object. (Or maybe store header ID in each entry?)
//func (this *Chain) createID() ChainID {
func TestCreateID( t *testing.T ) {
}

// return true when the chain is empty (Chain is a chain header message)
//func (this *Chain) isEmpty() bool {
func TestIsEmpty( t *testing.T ) {
}

// Get the previous entry in the chain
//func (this *Chain) Previous() *Chain {
func TestPrevious( t *testing.T ) {
}

// Identity of the chain's authorized identity
//func (this *Chain) author() *Ident {
func TestAuthor( t *testing.T ) {
}

// AuthChain.get( ChainID ) -> chainEntry
// Lookup and load a specific chain state. If the ID is a chainID,
// the chain is empty, it has no messages added to it.
//func Get( mID ChainID ) *Chain {
func TestGet( t *testing.T ) {
}

// getChainID: Returns ID of the zero entry (header/creation message) of the chain.
//func (this *Chain) getChainID() (id ChainID) {
func TestGetChainID( t *testing.T ) {
}

// addEntry( Chain ) Add message to the chain and returns the resulting chain.
//func (this *Chain) addEntry( m *Chain ) (c *Chain) {
func TestAddEntry( t *testing.T ) {
}


