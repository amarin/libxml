package libxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ParseXMLReader parses XML taken from io.Reader interface with XmlStreamParser.
func ParseXMLReader(reader io.Reader, streamParser XmlStreamParser) (err error) {
	var token xml.Token
	decoder := xml.NewDecoder(reader)

	for {
		if token, err = decoder.Token(); err != nil && errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return errors.New(fmt.Sprintf("token error: %v", err))
		} else if token == nil {
			return nil
		} else {
			t := xml.CopyToken(token)
			switch v := t.(type) {
			case xml.CharData:
				if err = streamParser.ProcessCharData(v); err != nil {
					return err
				}
			case xml.Comment:
				if err = streamParser.ProcessComment(v); err != nil {
					return err
				}
			case xml.Directive:
				if err = streamParser.ProcessDirective(v); err != nil {
					return err
				}
			case xml.ProcInst:
				if err = streamParser.ProcessProcInst(v); err != nil {
					return err
				}
			case xml.StartElement:
				if err = streamParser.ProcessStartElement(v); err != nil {
					return err
				}
			case xml.EndElement:
				if err = streamParser.ProcessEndElement(v); err != nil {
					return err
				}
			}
		}
	}
}

// ParseXMLFile parses XML file specified by fileName with XmlStreamParser.
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

	defer func() {
		if p := recover(); p != nil {
			_, _ = fmt.Fprintf(os.Stderr, "file %v\npanic: %v", fileName, p)
		}
		if closeErr := reader.Close(); closeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "file %v\nclose: %v", fileName, closeErr)
		}
	}()

	return ParseXMLReader(reader, streamParser)
}
