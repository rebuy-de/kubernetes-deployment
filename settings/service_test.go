package settings

import "testing"

func TestClean(t *testing.T) {
	var cleantests = []struct {
		in  Service
		out Service
	}{
		{
			Service{
				Repository: "github.com/rebuy-de/something",
			},
			Service{
				Name:       "something",
				Repository: "git@github.com:rebuy-de/something.git",
				Path:       "/deployment/kubernetes",
				Branch:     "master",
			},
		},
		{
			Service{
				Repository: "http://veryspecial.com/repos/foo",
			},
			Service{
				Name:       "http:--veryspecial.com-repos-foo",
				Repository: "http://veryspecial.com/repos/foo",
				Path:       "/deployment/kubernetes",
				Branch:     "master",
			},
		},
		{
			Service{
				Repository: "github.com/rebuy-de/blubber",
				Branch:     "special",
			},
			Service{
				Name:       "blubber-special",
				Repository: "git@github.com:rebuy-de/blubber.git",
				Path:       "/deployment/kubernetes",
				Branch:     "special",
			},
		},
		{
			Service{
				Repository: "github.com/rebuy-de/project",
				Branch:     "some-branch",
				Path:       "i/dont/care/about/conventions",
			},
			Service{
				Name:       "project-i-dont-care-about-conventions-some-branch",
				Repository: "git@github.com:rebuy-de/project.git",
				Path:       "/i/dont/care/about/conventions",
				Branch:     "some-branch",
			},
		},
	}

	for i, tt := range cleantests {
		generated := tt.in
		generated.Clean()

		if tt.out != generated {
			t.Errorf("Test %d failed:", i)
			t.Errorf("  Input:    %#v", tt.in)
			t.Errorf("  Output:   %#v", generated)
			t.Errorf("  Expected: %#v", tt.out)
		}
	}
}
