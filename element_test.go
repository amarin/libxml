package libxml_test

import (
	"reflect"
	"testing"

	"github.com/amarin/libxml"
)

func TestElement_Tag(t *testing.T) {
	for _, tt := range []struct {
		name  string
		path  string
		names []libxml.TagName
	}{
		{"1st_level", "/first", []libxml.TagName{"first"}},
		{"2nd_level", "/first/second", []libxml.TagName{"first", "second"}},
		{"3rd_level", "/first/second/third", []libxml.TagName{"first", "second", "third"}},
	} {
		tt := tt // pin tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt // pin tt
			root := libxml.MakeElement(nil, "", nil, nil, nil)
			got := root.Tag(tt.names...)
			path := got.Path()

			if path != tt.path {
				t.Errorf("Tag() = %v, want %v", path, tt.path)
			}

			got1 := root.Tag(tt.names...)
			if !reflect.DeepEqual(got, got1) {
				t.Errorf("Tag() = \n%v, same \n%v", got, got1)
			}
		})
	}
}
