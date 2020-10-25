package libxml

import (
	"encoding/xml"
)

// TagName defines string type to define XML tag name.
type TagName string

// TagElementMapping defines mapping of tag names to created elements.
// Used in Element structure to map children tag names to children elements.
type TagElementMapping map[TagName]*Element

// EnterFunc declares function interface to execute when entering some Element.
// It takes current element f.e. to extract data from tag attributes.
type EnterFunc func(element xml.StartElement) error

// ExitFunc declares function interface to execute when exiting some Element.
type ExitFunc func() error

// DataFunc declares function interface to execute when collecting data in specified Element.
type DataFunc func(data xml.CharData) error
