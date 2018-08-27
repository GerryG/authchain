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
		fmt.Printf("Has marshaler for %T %T\n", e, mar)
		err = mar.MarshalXML( enc, xml.StartElement{Name: xml.Name{Local: tag}} )
		return []byte(buf.String()), err
	} else {
		err = errors.New(fmt.Sprintf("Not a marshaler: %T", mar))
		return
	}
}

func (this *Element) MarshalXML(enc *xml.Encoder, st xml.StartElement) (err error) {
	//debug.PrintStack()
	//fmt.Printf("MaXML(%#v) %T, %#v, %d\n", st, this, this.Tok, len(this.Children))
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
	//fmt.Printf("St:%#v\n", tg)
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
		//fmt.Printf("Err: %#v\n", err)
		return err
	}
	//fmt.Printf("CH Unmarshaler %d\n", len(el.Children))
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
				//fmt.Printf("App subs: %#v\n", subel)
				this.App = &Element{Tok:subel.Tok, Children:subel.Children}
			}
		} else { return errors.New(fmt.Sprintf("Non-element in header: %#v", subel)) }
	}
	//fmt.Printf("CHMarsh:%#v\n", this)
	return nil
	//return decoder.DecodeElement(this, &st)
}

func (this *Element) UnmarshalXML(decoder *xml.Decoder, st xml.StartElement) error {
	var stack []*Element
	var top, closing *Element
	//fmt.Printf("Star UmXML %#v\n", st)
	for first := true; ; first = false {
		token, err := decoder.Token()
		//fmt.Printf("Tok[%T], %#v, %s\n", token, token, err)
		v, ok := token.(xml.StartElement)
		if first && ok && st.Name != v.Name {
			fmt.Printf("Push pre consumed start %#v\n", st)
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
					//fmt.Printf("UmXML Done: %#v\n", this)
					return nil
				}
				top = stack[len(stack)-1] // new top after pop
				top.Children = append(top.Children, *closing)
			case xml.CharData:
				if nonWhite.Find(v) != nil {
					top.Children = append(top.Children, Element{Tok:xml.CopyToken(v)})
				//} else {
				//	fmt.Printf("CD white[%d] %T\n", len(top.Children), v)
				}
			//case xml.Comment:
				//top.Children = append(top.Children, v)
			}
		} else {
			if err == io.EOF {
				if len(stack) > 0 {
					return errors.New(fmt.Sprintf("Unclosed: %v", stack))
				} else {
					*this = *top
					//fmt.Printf("UmXML DoneEOF: %#v\n", this)
					return nil
				}
			} else {
				return err
			}
		}
	}
}

