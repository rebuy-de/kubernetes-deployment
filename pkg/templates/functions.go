package templates

import (
	"regexp"
	"strings"
)

// Matches all invalid chars from rfc1035/rfc1123.
// See also: https://github.com/kubernetes/community/blob/master/contributors/design-proposals/identifiers.md
var IdentifierInvalidRe = regexp.MustCompile("[^a-z0-9]+")

func IdentifierFunc(s string) string {
	s = strings.ToLower(s)
	s = IdentifierInvalidRe.ReplaceAllLiteralString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
