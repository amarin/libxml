package libxml

import (
	"encoding/xml"
)

// XmlStreamParser requires implementations defines processing functions to deal with XML tree elements stream.
// All processing functions should return error if parsing should stopped
// as implemented algorithm is not suitable for incoming data.
//
// ProcessStartElement is invoked for each opening tag.
// It takes xml.StartElement to process.
//
// ProcessEndElement is invoked when tag closing.
// It takes xml.EndElement data to process and should return error
//
// ProcessCharData is invoked when tag character data received.
// It takes xml.CharData data to process and should return error
//
// ProcessComment is invoked when comment line encountered in XML stream.
// It takes xml.CharData data as byte string of comment content do not including the <!-- and --> comment markers.
//
// ProcessComment is invoked when processing instructions like <?target inst?> taken from XML stream.
// It takes xml.ProcInst data.
//
// ProcessDirective is invoked when processing directive of the form <!text> goes next in XML Sstream.
// It takes xml.Directive data.
type XmlStreamParser interface {
	ProcessStartElement(xml.StartElement) error
	ProcessEndElement(xml.EndElement) error
	ProcessCharData(xml.CharData) error
	ProcessComment(xml.Comment) error
	ProcessProcInst(xml.ProcInst) error
	ProcessDirective(xml.Directive) error
}
