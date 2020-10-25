package libxml

// Element type declares XML tree element data to use during parse.
type Element struct {
	parent    *Element
	tagName   TagName
	enter     EnterFunc
	data      DataFunc
	exit      ExitFunc
	structure TagElementMapping
}

// MakeElement creates new instance of Element.
func MakeElement(parent *Element, tag TagName, enter EnterFunc, data DataFunc, out ExitFunc) *Element {
	return &Element{
		parent:    parent,
		tagName:   tag,
		enter:     enter,
		data:      data,
		exit:      out,
		structure: make(TagElementMapping),
	}
}

// Tag defines element using direct tag name of child or chain of children lead to target element.
func (elem *Element) Tag(names ...TagName) *Element {
	firstTag := names[0]
	if existed, ok := elem.structure[firstTag]; ok {
		if len(names) > 1 {
			return existed.Tag(names[1:]...)
		}

		return existed
	}

	newElement := MakeElement(elem, firstTag, nil, nil, nil)
	elem.structure[firstTag] = newElement

	if len(names) > 1 {
		return newElement.Tag(names[1:]...)
	}

	return newElement
}

// Root returns root element of current element tree.
func (elem *Element) Root() *Element {
	if elem.parent == nil {
		return elem
	}

	return elem.parent.Root()
}

// Path returns tag names chain to arrive into current element from root.
func (elem *Element) Path() string {
	prefix := ""
	if elem.parent != nil {
		prefix = elem.parent.Path() + "/"
	}

	return prefix + string(elem.tagName)
}

// Parse adds current element child specified by tag name
// as well as functions to execute when child enter, child data processing and child exit.
// Returns created child element.
func (elem *Element) Parse(tag TagName, enter EnterFunc, data DataFunc, out ExitFunc) *Element {
	child := elem.Tag(tag)
	child.OnData(data)
	child.OnEnter(enter)
	child.OnExit(out)

	return child
}

// HasChild returns true if element has child with specified tag name.
func (elem Element) HasChild(tag TagName) bool {
	_, ok := elem.structure[tag]
	return ok
}

// OnEnter set function to execute on entering current element during parse.
func (elem *Element) OnEnter(fun EnterFunc) { elem.enter = fun }

// OnData set function to process current element data during parse.
func (elem *Element) OnData(fun DataFunc) { elem.data = fun }

// OnEnter set function to execute on exiting current element during parse.
func (elem *Element) OnExit(fun ExitFunc) { elem.exit = fun }
