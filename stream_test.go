package libxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
)

type testParser struct {
	startElement xml.StartElement
	endElement   xml.EndElement
	data         xml.CharData
	comment      xml.Comment
	inst         xml.ProcInst
	directive    xml.Directive
}

func (t *testParser) ProcessStartElement(element xml.StartElement) error {
	t.startElement = element.Copy()
	return nil
}

func (t *testParser) ProcessEndElement(element xml.EndElement) error {
	t.endElement = element
	return nil
}

func (t *testParser) ProcessCharData(data xml.CharData) error {
	t.data = data.Copy()
	return nil
}

func (t *testParser) ProcessComment(comment xml.Comment) error {
	t.comment = comment.Copy()
	return nil
}

func (t *testParser) ProcessProcInst(inst xml.ProcInst) error {
	t.inst = inst.Copy()
	return nil
}

func (t *testParser) ProcessDirective(directive xml.Directive) error {
	t.directive = directive.Copy()
	return nil
}

func (t *testParser) String() string {
	attrStrings := make([]string, 0)
	for _, attr := range t.startElement.Attr {
		attrStrings = append(attrStrings, fmt.Sprintf(`%s="%s"`, attr.Name, attr.Value))
	}
	return fmt.Sprintf(
		"inst: `%v %v`\ndirective: `%v`\nstart: <%s:%s %s>\nend: <%s:%v>\ncomment: `%v`\ndata: `%v`\n",
		t.inst.Target,
		string(t.inst.Inst),
		string(t.directive),
		t.startElement.Name.Space, t.startElement.Name.Local,
		strings.Join(attrStrings, ","),
		t.endElement.Name.Space, t.endElement.Name.Local,
		string(t.comment),
		string(t.data),
	)
}

func TestParseXMLReader(t *testing.T) {
	target := "xml"
	inst := `version="1.0" encoding="utf-8" standalone="yes"`
	instruction := fmt.Sprintf("%v %v", target, inst)
	comment := "comment"

	tag := "example"
	chars := "content"

	directive := "COMMENT "
	exampleXML := fmt.Sprintf(
		`<?%v?>
<!--%v-->
<!%v>
<%v>%v</%v>`,
		instruction,
		comment,
		directive,
		tag, chars, tag,
	)
	reader := bytes.NewBufferString(exampleXML)
	parser := new(testParser)

	if err := ParseXMLReader(reader, parser); err != nil {
		t.Errorf("ParseXMLReader() error = %v\n\n%v\n\n%v", err, exampleXML, parser)
	} else if string(parser.directive) != directive {
		t.Errorf("expected directive `%v`, got `%v`\n\n%v", directive, string(parser.directive), exampleXML)
	} else if string(parser.comment) != comment {
		t.Errorf("expected directive `%v`, got `%v`\n\n%v", comment, string(parser.comment), exampleXML)
	} else if parser.inst.Target != target {
		t.Errorf("expected instruction target `%v`, got `%v`\n\n%v", target, parser.inst.Target, exampleXML)
	} else if inst != string(parser.inst.Inst) {
		t.Errorf("expected instruction `%v`, got `%v`\n\n%v", inst, string(parser.inst.Inst), exampleXML)
	} else if tag != parser.startElement.Name.Local {
		t.Errorf("expected tag `%v`, got `%v`\n\n%v", tag, parser.startElement.Name.Local, exampleXML)
	} else if tag != parser.endElement.Name.Local {
		t.Errorf("expected tag `/%v`, got `/%v`\n\n%v", tag, parser.endElement.Name.Local, exampleXML)
	} else if chars != string(parser.data) {
		t.Errorf("expected data `%v`, got `%v`\n\n%v\n\n%v", chars, string(parser.data), exampleXML, parser)
	}
}
