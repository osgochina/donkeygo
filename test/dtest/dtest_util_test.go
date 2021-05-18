package dtest_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestC(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(123.456, 123.456)
	})
}

func TestCase(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(123.456, 123.456)
	})
}
