package dstr_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"testing"
)

func TestSnakeString(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.SnakeString("XxYy"), "xx_yy")
	})
}
