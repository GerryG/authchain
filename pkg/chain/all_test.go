
package chain_test

import (
	"os"
	"io"
	"fmt"
	//"io/ioutil"
	//"encoding/xml"
	"testing"
	"path/filepath"
	. "github.com/GerryG/authchain/pkg/chain"
)

// shared testing methods to set up fixtures and test objects
type testInst struct {
	Name string
}

func testdataInputFrom( t *testing.T, base, tn string, exts ...string ) (in io.Reader, ext string) {
	var err error
	tn = base + "_" + tn
	for _, ext = range exts {
		in, err = os.Open( tn+ext )
		if err == nil {
			//fmt.Printf("found: %s\n", tn+ext)
			return
		}
	}

	for _, ext = range exts {
		in, err = os.Open( base+ext )
		if err == nil {
			//fmt.Printf("found base: %s\n", base+ext)
			return
		}
	}

	return nil, ""
}

func testdataIdent( t *testing.T, tn string ) (id *Ident) {
	 //filepath.Join("testdata", "testcert.pem"), filepath.Join("testdata", "testkey.pem") )
	testBase := filepath.Join("testdata", "ident")
	in, ext := testdataInputFrom( t, testBase, tn, "", ".xml", ".json", ".crt" )
	id = LoadIdent( in, ext)
	id.Secure.Unlock(id.Certificate.PEM)
	return
}

func testdataCreate( t *testing.T, tn string ) (doc *Element) {
	testBase := filepath.Join("testdata", "create")
	in, ext := testdataInputFrom( t, testBase, tn, "", ".xml", ".json" )
	defer in.(*os.File).Close()
	if ext == ".xml" {
		// read and decode XML file
		doc := &Element{}
		err := Unmarshal("app", in, doc)
		if err != nil {
			panic(fmt.Sprintf("Unmarshal error: %s", err))
		}
		return doc
	}
	return nil
}

func testdataMessage( t *testing.T, tn string ) *Element {
	testBase := filepath.Join("testdata", "message")
	in, ext := testdataInputFrom( t, testBase, tn, "", ".xml", ".json", ".crt" )
	if ext == ".xml" {
		// read and decode XML file
		entry := &Element{}
		err := Unmarshal("entry", in, entry)
		if err != nil {
			panic(err)
		}
		return entry
	}
	return nil
}

