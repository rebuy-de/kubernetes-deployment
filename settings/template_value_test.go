package settings

import (
	"reflect"
	"testing"
)

func TestTemplateValueMerge(t *testing.T) {
	tv1 := TemplateValues{
		"bish": "a",
		"bash": "foo",
	}

	tv2 := TemplateValues{
		"bosh": "blubber",
		"bash": "bar",
	}

	tv12 := TemplateValues{
		"bosh": "blubber",
		"bash": "bar",
		"bish": "a",
	}

	tv21 := TemplateValues{
		"bosh": "blubber",
		"bash": "foo",
		"bish": "a",
	}

	tv12g := tv1.Merge(tv2)
	tv21g := tv2.Merge(tv1)

	if !reflect.DeepEqual(tv12, tv12g) {
		t.Errorf("Generated values are wrong for merging tv2 into tv1:")
		t.Errorf("  tv1:       %#v", tv1)
		t.Errorf("  tv2:       %#v", tv2)
		t.Errorf("  expected:  %#v", tv12)
		t.Errorf("  generated: %#v", tv12g)
		t.Fail()
	}

	if !reflect.DeepEqual(tv21, tv21g) {
		t.Errorf("Generated values are wrong for merging tv1 into tv2:")
		t.Errorf("  tv1:       %#v", tv1)
		t.Errorf("  tv2:       %#v", tv2)
		t.Errorf("  expected:  %#v", tv21)
		t.Errorf("  generated: %#v", tv21g)
		t.Fail()
	}

	var tvNil TemplateValues
	tvNilG := tvNil.Merge(tv1)

	if !reflect.DeepEqual(tv21, tv21g) {
		t.Errorf("Generated values are wrong for merging tv1 into nil:")
		t.Errorf("  tvNil:     %#v", tvNil)
		t.Errorf("  tv1:       %#v", tv1)
		t.Errorf("  expected:  %#v", tv1)
		t.Errorf("  generated: %#v", tvNilG)
		t.Fail()
	}

}