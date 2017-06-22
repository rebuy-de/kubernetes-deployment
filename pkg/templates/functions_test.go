package templates

import "testing"

func TestIdentifierFunc(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"-ab-:_.fSDf123axxb-", "ab-fsdf123axxb"},
		{"------", ""},
		{"angular-4.2.0", "angular-4-2-0"},
		{"BLUE-1337__fancy-feature", "blue-1337-fancy-feature"},
	}

	for i, tc := range tests {
		out := IdentifierFunc(tc.in)
		if out != tc.out {
			t.Errorf("Test Case %d failed: %#v != %#v", i, out, tc.out)
		}
	}
}
