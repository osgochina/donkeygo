// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

func Test_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		f32 := float32(123.456)
		i64 := int64(1552578474888)
		t.AssertEQ(dconv.Int(f32), int(123))
		t.AssertEQ(dconv.Int8(f32), int8(123))
		t.AssertEQ(dconv.Int16(f32), int16(123))
		t.AssertEQ(dconv.Int32(f32), int32(123))
		t.AssertEQ(dconv.Int64(f32), int64(123))
		t.AssertEQ(dconv.Int64(f32), int64(123))
		t.AssertEQ(dconv.Uint(f32), uint(123))
		t.AssertEQ(dconv.Uint8(f32), uint8(123))
		t.AssertEQ(dconv.Uint16(f32), uint16(123))
		t.AssertEQ(dconv.Uint32(f32), uint32(123))
		t.AssertEQ(dconv.Uint64(f32), uint64(123))
		t.AssertEQ(dconv.Float32(f32), float32(123.456))
		t.AssertEQ(dconv.Float64(i64), float64(i64))
		t.AssertEQ(dconv.Bool(f32), true)
		t.AssertEQ(dconv.String(f32), "123.456")
		t.AssertEQ(dconv.String(i64), "1552578474888")
	})

	dtest.C(t, func(t *dtest.T) {
		s := "-0xFF"
		t.Assert(dconv.Int(s), int64(-0xFF))
	})
}

func Test_Duration(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		d := dconv.Duration("1s")
		t.Assert(d.String(), "1s")
		t.Assert(d.Nanoseconds(), 1000000000)
	})
}
