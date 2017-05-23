package templates

import "testing"

func TestRender(t *testing.T) {
	var (
		template = `
			branch: {{.branch}}
			commit: {{.sha}}
		`
		values = Values{
			"branch": "foobar",
			"sha":    "123abc",
		}
		expected = `
			branch: foobar
			commit: 123abc
		`
	)

	result, err := Render(template, values)
	if err != nil {
		t.Fatal(err)
	}

	if expected != result {
		t.Errorf("Rendered templated is wrong")
		t.Errorf("  Expected: %#v", expected)
		t.Errorf("  Obtained: %#v", result)
	}
}
