package interceptors

import "testing"

func TestMultiType(t *testing.T) {
	var multi Interface
	multi = &Multi{}
	_ = multi
}
