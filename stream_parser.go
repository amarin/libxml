package libxml

/* StreamParser implements declarative XML structure parsing.
It takes structure definition to define document tags or can setup structure and tags processing after creation.
*/
import (
	"encoding/xml"
	"errors"
	"fmt"
)

var (
	errUnexpectedTag        = errors.New("unexpected tag")
	errUnexpectedTransition = errors.New("unexpected transition")
)

// StreamParser implements libxml.XmlStreamParser.
// It requires parsing structure created before parse with Tag and Parse methods.
type StreamParser struct {
	strict  bool
	current *Element
}

// Tag adds and returns new element with specified tag name.
// Parsing functions for tag should set using Element OnEnter, OnData and OnExit methods.
func (parser *StreamParser) Tag(names ...TagName) *Element {
	return parser.current.Root().Tag(names...)
}

// Parse creates new root element one time with parsing functions set for enter, processing data and tag exit.
func (parser *StreamParser) Parse(name TagName, in EnterFunc, dataFunc DataFunc, out ExitFunc) *Element {
	return parser.current.Root().Parse(name, in, dataFunc, out)
}

// ProcessComment implements comments parse. Implements XmlStreamParser.
func (parser StreamParser) ProcessComment(_ xml.Comment) error { return nil }

// ProcessProcInst implements processing instruction parse. Implements XmlStreamParser.
func (parser StreamParser) ProcessProcInst(_ xml.ProcInst) error { return nil }

// ProcessDirective implements directive parse. Implements XmlStreamParser.
func (parser StreamParser) ProcessDirective(_ xml.Directive) error { return nil }

// Root returns root Element.
func (parser *StreamParser) Root() *Element {
	return parser.current.Root()
}

// ProcessStartElement processes token start. Implements XmlStreamParser.
// It checks if current Element has specified Tag child, set detected child element as current
// and executes OnEnter function if defined.
// If Strict mode set and no child defined returns with error.
func (parser *StreamParser) ProcessStartElement(token xml.StartElement) (err error) {
	var (
		next *Element
		ok   bool
	)

	tagName := TagName(token.Name.Local)
	next, ok = parser.current.structure[tagName]

	switch {
	case parser.strict && !ok:
		return fmt.Errorf("%w: %v.%v", errUnexpectedTransition, parser.current.tagName, token.Name.Local)
	case !ok && !parser.strict:
		next = parser.current.Tag(tagName)
	}
	// set current element
	parser.current = next
	// take element enter function and call it if set
	if next.enter != nil {
		return next.enter(token)
	}

	return nil
}

// ProcessCharData processes tag (string) data. Processing Element should have OnData function set.
// Returns error from current tag processing function if encountered.
func (parser *StreamParser) ProcessCharData(data xml.CharData) error {
	if parser.current.data != nil {
		return parser.current.data(data)
	}

	return nil
}

// ProcessEndElement processes tag closing.
// Calls ExitFunc for current element if set with OnExit method of Element.
// If no exit processing function set to current Element it simply returns previous element.
// Returns ExitFunc error if encountered.
func (parser *StreamParser) ProcessEndElement(element xml.EndElement) (err error) {
	tagName := TagName(element.Name.Local)

	if parser.current.tagName != tagName {
		return fmt.Errorf("%w: closing %v in %v", errUnexpectedTag, tagName, parser.current.tagName)
	}

	defer func() {
		if parser.current.parent != nil { // move to outer element
			parser.current = parser.current.parent
		}
	}()

	if parser.current.exit != nil {
		return parser.current.exit()
	}

	return nil
}

// NewParser creates new parser.
// If strict is false, any unexpected tag will be silently ignored.
// If strict is true structure should define all required tag elements
// with parser root or some element tag method, f.e. parser.Tag("ignore", "this", "tree").
// In struct mode any unexpected tag will produce parsing error.
func NewParser(strict bool) *StreamParser {
	p := &StreamParser{
		strict:  strict,
		current: MakeElement(nil, RootTag, nil, nil, nil),
	}

	return p
}
