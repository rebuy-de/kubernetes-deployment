package templates

import (
	"reflect"
	"testing"
)

func TestValueDefaults(t *testing.T) {
	def := Variables{
		"A": "1",
		"B": "2",
	}

	values := Variables{
		"B": "3",
		"C": "4",
	}

	expect := Variables{
		"A": "1",
		"B": "3",
		"C": "4",
	}

	values.Defaults(def)

	if !reflect.DeepEqual(expect, values) {
		t.Errorf("Merged values doesn't match.")
		t.Errorf("    Expected:  %#v", expect)
		t.Errorf("    Generated: %#v", values)

	}
}
