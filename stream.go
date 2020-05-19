package libxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type XmlStreamParser interface {
	ProcessStartElement(xml.StartElement) error
	ProcessEndElement(xml.EndElement) error
	ProcessCharData(xml.CharData) error
	ProcessComment(xml.Comment) error
	ProcessProcInst(xml.ProcInst) error
	ProcessDirective(xml.Directive) error
}

// Parse XML from io.Reader interface with XmlStreamParser
func ParseXMLReader(reader io.Reader, streamParser XmlStreamParser) error {
	decoder := xml.NewDecoder(reader)
	for {
		if token, err := decoder.Token(); err != nil && errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return errors.New(fmt.Sprintf("token error: %v", err))
		} else if token == nil {
			return nil
		} else {
			t := xml.CopyToken(token)
			switch v := t.(type) {
			case xml.CharData:
				if err := streamParser.ProcessCharData(v); err != nil {
					return err
				}
			case xml.Comment:
				if err := streamParser.ProcessComment(v); err != nil {
					return err
				}
			case xml.Directive:
				if err := streamParser.ProcessDirective(v); err != nil {
					return err
				}
			case xml.ProcInst:
				if err := streamParser.ProcessProcInst(v); err != nil {
					return err
				}
			case xml.StartElement:
				if err := streamParser.ProcessStartElement(v); err != nil {
					return err
				}
			case xml.EndElement:
				if err := streamParser.ProcessEndElement(v); err != nil {
					return err
				}
			}
		}
	}
}

// Parse XML from fileName with XmlStreamParser
func ParseXMLFile(fileName string, streamParser interface{}) error {
	if absPath, err := filepath.Abs(fileName); err != nil {
		return err
	} else if reader, err := os.Open(absPath); err != nil {
		return err
	} else if parser, ok := streamParser.(XmlStreamParser); ok {
		return ParseXMLReader(reader, parser)
	} else if parser, ok := streamParser.(XMLInstanceGrabber); ok {
		return GrabTargetsFromXML(reader, parser)
	} else {
		return errors.New(fmt.Sprintf("unexpected parser %T", streamParser))
	}
}
