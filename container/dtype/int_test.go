package dtype_test

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestInt_Reduce(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		n := dtype.NewInt(10)
		t.Assert(n.Reduce(2), 8)
	})
}
