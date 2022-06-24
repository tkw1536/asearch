package mws

import (
	"bytes"
	"encoding/xml"

	"github.com/pkg/errors"
)

// RawMWSQuery represents a (raw) MathWebSearch Query that is sent directly to MathWebSearch
type RawMWSQuery struct {
	From int64 `xml:"limitmin,attr"` // offset within the set of results
	Size int64 `xml:"answsize,attr"` // maximum number of results returned

	ReturnTotal  BooleanYesNo `xml:"totalreq,attr"` // if true also compute the total number of elements
	OutputFormat string       `xml:"output,attr"`   // output format, "xml" or "json"

	Expressions []MWSExpression // the expressions that we are searching for
}

// MarshalXML marshales a raw query as XML
func (raw RawMWSQuery) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type marshalRawQuery RawMWSQuery // to prevent infinite recursion
	r := struct {
		marshalRawQuery

		NamespaceMWS string `xml:"xmlns:mws,attr"`
		NamespaceM   string `xml:"xmlns:m,attr"`
	}{
		marshalRawQuery(raw),
		"http://www.mathweb.org/mws/ns",
		"http://www.w3.org/1998/Math/MathML",
	}
	start.Name = xml.Name{Local: "mws:query", Space: ""} // TODO: Fixme, why is this not working
	return errors.Wrap(e.EncodeElement(r, start), "e.EncodeElement failed")
}

// MWSExpression represents a single expression that is being searched for
type MWSExpression struct {
	XMLName xml.Name `xml:"mws:expr"`
	Term    string   `xml:",innerxml"` // the actual term being searched for
}

// BooleanYesNo represents a boolean that is xml encoded as "yes" or "no"
type BooleanYesNo bool

// MarshalText turns a BooleanYesNo into a string
func (byesno BooleanYesNo) MarshalText() (text []byte, err error) {
	if byesno {
		text = []byte("yes")
	} else {
		text = []byte("no")
	}
	return
}

// UnmarshalText unmarshals text into a string
func (byesno *BooleanYesNo) UnmarshalText(text []byte) (err error) {
	// load yes and no
	if bytes.EqualFold(yesBytes, text) {
		*byesno = true
	} else if bytes.EqualFold(noBytes, text) {
		*byesno = false

		// do not load the else
	} else {
		err = errors.Errorf("Boolean should be \"yes\" or \"no\", not %q", string(text))
	}

	return
}

var noBytes = []byte("no")
var yesBytes = []byte("yes")
