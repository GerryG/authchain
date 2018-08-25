package chain

import (
	"regexp"
	"fmt"
	"os"
	"io"
	"bytes"
	"errors"
	//"math/big"
	"encoding/xml"
)

// xml for the chain header message
type ChainHeader struct {
	XMLName xml.Name `xml:"header"`
	App Application
}

type Application struct {
	XMLName xml.Name `xml:"app"`
	Appname string   `xml:"name,attr"`
	Dna string       `xml:"dna,attr"`
	Auths []Author   `xml:"player"`
}

type Author struct {
	XMLName xml.Name `xml:"player"`
	Role string      `xml:"role,attr"`
	ID string        `xml:"id,attr"`
	Sec *Secure      `xml:"secure"`
	Name string      `xml:"name"`
	Nick string      `xml:"nick"`
	ChainID *string   `xml:"chain"`
}

type Secure struct {
	XMLName xml.Name `xml:"secure"`
	DevID string  `xml:"deviceid,attr"`
}

// xml for the chain entry message
type ChainEntry struct {
	XMLName xml.Name `xml:"entry"`
	Acts []interface{} `xml:"move"`
}

// for identity file
type Identity struct {
	XMLName xml.Name `xml:"identity"`
	Name	string
	Domain	string
	Certificate Certificate
	Secure	SecureID
}

type SecureID struct {
	XMLName xml.Name `xml:"secure"`
	Id	string   `xml:"id,attr"`
	Device	string   `xml:"device,attr"`
	PEM	string   `xml:"pem"`
}

type Certificate struct {
	XMLName xml.Name `xml:"certificate"`
	Id	string   `xml:"id,attr"`
	Name	string   `xml:"name,attr"`
	PEM	string   `xml:"pem"`
}

type Element struct {
	Open xml.StartElement
	Children []xml.Token
}

func Marshal(e *Element) (b []byte, err error) {
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	err = e.Encode( enc, buf )
	if err != nil { return }
	err = enc.Flush()
	return []byte(buf.String()), err
}

func (e *Element) Encode(enc *xml.Encoder, buf *bytes.Buffer) (err error) {
	err = enc.EncodeToken(e.Open)
	if err != nil { return }
	for _, child := range(e.Children) {
		ce, ok := child.(Element)
		if ok {
			err = ce.Encode( enc, buf )
			if err != nil { return }
		} else {
			err = enc.EncodeToken(child)
			if err != nil { return }
		}
	}
	return enc.EncodeToken(e.Open.End())
}

var nonWhite *regexp.Regexp = regexp.MustCompile(`[^\s\p{Zs}]`)

func Unmarshal(in os.File) (*Element, error) {
	decoder := xml.NewDecoder(io.Reader(&in))
	var stack []*Element
	var top, closing *Element
	for {
		token, err := decoder.Token()
		if err == nil {
			switch v := token.(type) {
			case xml.StartElement:
				top = &Element{Open: v}
				stack = append(stack, top) // push stack
			case xml.EndElement:
				if v.Name != top.Open.Name {
					return top, errors.New(fmt.Sprintf("Open and close differ %s :: %s", v.Name, top.Open.Name))
				}
				closing = top // top is what would be popped
				stack = stack[:len(stack)-1] // pop stack
				if len(stack) == 0 {
					return closing, nil
				}
				top = stack[len(stack)-1] // new top after pop
				top.Children = append(top.Children, *closing)
			case xml.CharData:
				if nonWhite.Find(v) != nil {
					top.Children = append(top.Children, xml.CopyToken(v))
				//} else {
				//	fmt.Printf("CD white[%d] %T\n", len(top.Children), v)
				}
			//case xml.Comment:
				//top.Children = append(top.Children, v)
			}
		} else {
			if err == io.EOF {
				if len(stack) > 0 {
					return top, errors.New(fmt.Sprintf("Unclosed: %v", stack))
				} else { return top, nil }
			} else { return top, err }
		}
	}
}

