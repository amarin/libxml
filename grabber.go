package libxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type XmlGrabberMakeInstanceFun func() interface{}
type XmlGrabberProcessInstanceFun func(interface{}) error

type XMLInstanceGrabber interface {
	SetName(xml.Name)
	GetName() xml.Name
	MakeInstance() interface{}
	ProcessInstance(interface{}) error
}

type XmlGrabber struct {
	catchName  xml.Name
	makeFun    XmlGrabberMakeInstanceFun
	processFun XmlGrabberProcessInstanceFun
}

func (x *XmlGrabber) MakeInstance() interface{} {
	return x.makeFun()
}

func (x *XmlGrabber) ProcessInstance(i interface{}) error {
	return x.processFun(i)
}

func (x *XmlGrabber) SetName(name xml.Name) {
	x.catchName = name
}

func (x XmlGrabber) GetName() xml.Name {
	return x.catchName
}

func NewXmlGrabber(tag string, ns string, makeInstanceFun XmlGrabberMakeInstanceFun, processInstanceFun XmlGrabberProcessInstanceFun) *XmlGrabber {
	name := xml.Name{
		Space: ns,
		Local: tag,
	}
	return &XmlGrabber{
		catchName:  name,
		makeFun:    makeInstanceFun,
		processFun: processInstanceFun,
	}
}

func NewXmlTargetGrabber(tag string, ns string, target interface{}, processInstanceFun XmlGrabberProcessInstanceFun) *XmlGrabber {
	targetType := reflect.TypeOf(target)
	return NewXmlGrabber(tag, ns, func() interface{} { return reflect.New(targetType).Interface() }, processInstanceFun)
}

// Parse XML from io.Reader interface with XmlStreamParser
func GrabTargetsFromXML(reader io.Reader, streamParser XMLInstanceGrabber) error {
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
			case xml.StartElement:
				if v.Name.Local == streamParser.GetName().Local && v.Name.Space == streamParser.GetName().Space {
					newInstance := streamParser.MakeInstance()
					if err := decoder.DecodeElement(newInstance, &v); err != nil {
						return err
					} else if err := streamParser.ProcessInstance(newInstance); err != nil {
						return err
					}
				}
			}
		}
	}
}
