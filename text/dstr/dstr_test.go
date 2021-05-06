package dstr_test

import (
	"donkeygo/test/dtest"
	"donkeygo/text/dstr"
	"testing"
)

func TestSnakeString(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.SnakeString("XxYy"), "xx_yy")
	})
}
