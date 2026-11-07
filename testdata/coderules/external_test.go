package coderules_test

import (
	"testing"

	"example.com/coderules"
)

func TestExternal(t *testing.T) {
	if coderules.Exported(2) != 2 {
		t.Fail()
	}
	_ = coderules.Exported(2) * 3.0 // want fused
}
