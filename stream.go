package libxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Parse XML from io.Reader interface with XmlStreamParser.
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

// Parse XML from fileName with XmlStreamParser.
func ParseXMLFile(fileName string, streamParser XmlStreamParser) (err error) {
	var (
		absPath string
		reader  io.ReadCloser
	)
	absPath, err = filepath.Abs(fileName)
	if err != nil {
		return err
	}

	reader, err = os.Open(absPath)
	if err != nil {
		return err
	}

	defer reader.Close()

	return ParseXMLReader(reader, streamParser)
}
