package libxml

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestGrabTargetsFromXML(t *testing.T) {
	type testTarget struct {
		Value string `xml:",chardata"`
	}
	grabbedItems := make([]testTarget, 0)
	grabber := NewXmlGrabber(
		"e", "",
		func() interface{} {
			return &testTarget{}
		},
		func(item interface{}) error {
			if expectedItem, ok := item.(*testTarget); ok {
				grabbedItems = append(grabbedItems, *expectedItem)
				return nil
			} else {
				return errors.New(fmt.Sprintf("unexpected item %T", item))
			}
		})
	exampleXML := `<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<e>1</e>
<e>2</e>
`
	itemsCount := 2
	reader := bytes.NewBufferString(exampleXML)
	if err := GrabTargetsFromXML(reader, grabber); err != nil {
		t.Errorf("GrabTargetsFromXML() error = %v\n\n%v", err, exampleXML)
	} else if len(grabbedItems) != itemsCount {
		t.Errorf("Expected %v items, got %v\n\n%v", itemsCount, len(grabbedItems), exampleXML)
	}
}

func TestNewXmlTargetGrabber(t *testing.T) {
	type testTarget struct {
		Value string `xml:",chardata"`
	}
	grabbedItems := make([]testTarget, 0)
	grabber := NewXmlTargetGrabber("e", "", testTarget{},
		func(item interface{}) error {
			if expectedItem, ok := item.(*testTarget); ok {
				grabbedItems = append(grabbedItems, *expectedItem)
				return nil
			} else {
				return errors.New(fmt.Sprintf("unexpected item %T", item))
			}
		})
	exampleXML := `<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<e>1</e>
<e>2</e>
`
	itemsCount := 2
	reader := bytes.NewBufferString(exampleXML)
	if err := GrabTargetsFromXML(reader, grabber); err != nil {
		t.Errorf("GrabTargetsFromXML() error = %v\n\n%v", err, exampleXML)
	} else if len(grabbedItems) != itemsCount {
		t.Errorf("Expected %v items, got %v\n\n%v", itemsCount, len(grabbedItems), exampleXML)
	}
}
