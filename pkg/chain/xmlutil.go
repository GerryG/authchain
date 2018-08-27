package chain

import (
	//"runtime/debug"
	"regexp"
	"fmt"
	"io"
	//"os"
	"time"
	"bytes"
	"errors"
	//"math/big"
	"encoding/xml"
)

type Element struct {
	Tok xml.Token      `xml:"-"`
	Children []Element `xml:"app"`
	xml.Marshaler
	xml.Unmarshaler
}

func Marshal(tag string, e interface{}) (b []byte, err error) {
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	mar, ok := e.(xml.Marshaler)
	if ok {
		err = mar.MarshalXML( enc, xml.StartElement{Name: xml.Name{Local: tag}} )
		return []byte(buf.String()), err
	} else {
		err = errors.New(fmt.Sprintf("Not a marshaler: %T", mar))
		return
	}
}

func (this *Element) MarshalXML(enc *xml.Encoder, st xml.StartElement) (err error) {
	err = this.Encode( enc )
	if err != nil { return }
	err = enc.Flush()
	return
}

func (this *Element) Encode(enc *xml.Encoder) (err error) {
	err = enc.EncodeToken(this.Tok)
	if err != nil { return }
	for _, child := range(this.Children) {
		switch v := child.Tok.(type) {
		case xml.StartElement:
			err = child.Encode( enc )
			if err != nil { return }
		case xml.CharData:
			err = enc.EncodeToken(v)
			if err != nil { return }
		}
	}
	return enc.EncodeToken(this.Tok.(xml.StartElement).End())
}

var nonWhite *regexp.Regexp = regexp.MustCompile(`[^\s\p{Zs}]`)

func Unmarshal(tag string, in io.Reader, el xml.Unmarshaler) error {
	decoder := xml.NewDecoder(in)
	tg := xml.StartElement{Name:xml.Name{Local: tag}}
	ret := el.UnmarshalXML(decoder, tg)
	return ret
}

func FindAttr(aname string, ats []xml.Attr) string {
	for _, at := range ats {
		if at.Name.Local == aname { return string(at.Value) }
	}
	return ""
}
func CharData(children []Element) (cd string) {
	for _, el := range children {
		s, ok := el.Tok.(xml.CharData)
		if ok { cd = cd+string([]byte(s)) }
	}
	return
}

func (this *ChainHeader) UnmarshalXML(decoder *xml.Decoder, st xml.StartElement) error {
	el := &Element{}
	err := el.UnmarshalXML(decoder, st)
	if err != nil {
		return err
	}
	for _, subel := range el.Children {
		st, ok := subel.Tok.(xml.StartElement)
		if ok {
			switch st.Name.Local {
			case "author":
				this.Author = CharData(subel.Children)
			case "at":
				t, err := time.Parse(time.RFC3339, CharData(subel.Children))
				if err != nil { return err }
				this.At = t
			case "app":
				this.App = &Element{Tok:subel.Tok, Children:subel.Children}
			}
		} else { return errors.New(fmt.Sprintf("Non-element in header: %#v", subel)) }
	}
	return nil
}

func (this *Element) UnmarshalXML(decoder *xml.Decoder, st xml.StartElement) error {
	var stack []*Element
	var top, closing *Element
	for first := true; ; first = false {
		token, err := decoder.Token()
		v, ok := token.(xml.StartElement)
		if first && ok && st.Name != v.Name {
			top = &Element{Tok: st}
			stack = append(stack, top)
		}
		if err == nil {
			switch v := token.(type) {
			case xml.StartElement:
				top = &Element{Tok: v}
				stack = append(stack, top) // push stack
			case xml.EndElement:
				if v.Name != top.Tok.(xml.StartElement).Name {
					return errors.New(fmt.Sprintf("Open and close differ %s :: %s", v.Name, top.Tok.(xml.StartElement).Name))
				}
				closing = top // top is what would be popped
				stack = stack[:len(stack)-1] // pop stack
				if len(stack) == 0 {
					*this = *top
					return nil
				}
				top = stack[len(stack)-1] // new top after pop
				top.Children = append(top.Children, *closing)
			case xml.CharData:
				if nonWhite.Find(v) != nil {
					top.Children = append(top.Children, Element{Tok:xml.CopyToken(v)})
			}
		} else {
			if err == io.EOF {
				if len(stack) > 0 {
					return errors.New(fmt.Sprintf("Unclosed: %v", stack))
				} else {
					*this = *top
					return nil
				}
			} else {
				return err
			}
		}
	}
}

