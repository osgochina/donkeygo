// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dutil_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

func Test_ComparatorString(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorString(1, 1), 0)
		t.Assert(dutil.ComparatorString(1, 2), -1)
		t.Assert(dutil.ComparatorString(2, 1), 1)
	})
}

func Test_ComparatorInt(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorInt(1, 1), 0)
		t.Assert(dutil.ComparatorInt(1, 2), -1)
		t.Assert(dutil.ComparatorInt(2, 1), 1)
	})
}

func Test_ComparatorInt8(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorInt8(1, 1), 0)
		t.Assert(dutil.ComparatorInt8(1, 2), -1)
		t.Assert(dutil.ComparatorInt8(2, 1), 1)
	})
}

func Test_ComparatorInt16(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorInt16(1, 1), 0)
		t.Assert(dutil.ComparatorInt16(1, 2), -1)
		t.Assert(dutil.ComparatorInt16(2, 1), 1)
	})
}

func Test_ComparatorInt32(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorInt32(1, 1), 0)
		t.Assert(dutil.ComparatorInt32(1, 2), -1)
		t.Assert(dutil.ComparatorInt32(2, 1), 1)
	})
}

func Test_ComparatorInt64(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorInt64(1, 1), 0)
		t.Assert(dutil.ComparatorInt64(1, 2), -1)
		t.Assert(dutil.ComparatorInt64(2, 1), 1)
	})
}

func Test_ComparatorUint(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorUint(1, 1), 0)
		t.Assert(dutil.ComparatorUint(1, 2), -1)
		t.Assert(dutil.ComparatorUint(2, 1), 1)
	})
}

func Test_ComparatorUint8(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorUint8(1, 1), 0)
		t.Assert(dutil.ComparatorUint8(2, 6), 252)
		t.Assert(dutil.ComparatorUint8(2, 1), 1)
	})
}

func Test_ComparatorUint16(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorUint16(1, 1), 0)
		t.Assert(dutil.ComparatorUint16(1, 2), 65535)
		t.Assert(dutil.ComparatorUint16(2, 1), 1)
	})
}

func Test_ComparatorUint32(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorUint32(1, 1), 0)
		t.Assert(dutil.ComparatorUint32(-1000, 2147483640), 2147482656)
		t.Assert(dutil.ComparatorUint32(2, 1), 1)
	})
}

func Test_ComparatorUint64(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorUint64(1, 1), 0)
		t.Assert(dutil.ComparatorUint64(1, 2), -1)
		t.Assert(dutil.ComparatorUint64(2, 1), 1)
	})
}

func Test_ComparatorFloat32(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorFloat32(1, 1), 0)
		t.Assert(dutil.ComparatorFloat32(1, 2), -1)
		t.Assert(dutil.ComparatorFloat32(2, 1), 1)
	})
}

func Test_ComparatorFloat64(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorFloat64(1, 1), 0)
		t.Assert(dutil.ComparatorFloat64(1, 2), -1)
		t.Assert(dutil.ComparatorFloat64(2, 1), 1)
	})
}

func Test_ComparatorByte(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorByte(1, 1), 0)
		t.Assert(dutil.ComparatorByte(1, 2), 255)
		t.Assert(dutil.ComparatorByte(2, 1), 1)
	})
}

func Test_ComparatorRune(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorRune(1, 1), 0)
		t.Assert(dutil.ComparatorRune(1, 2), -1)
		t.Assert(dutil.ComparatorRune(2, 1), 1)
	})
}

func Test_ComparatorTime(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		j := dutil.ComparatorTime("2019-06-14", "2019-06-14")
		t.Assert(j, 0)

		k := dutil.ComparatorTime("2019-06-15", "2019-06-14")
		t.Assert(k, 1)

		l := dutil.ComparatorTime("2019-06-13", "2019-06-14")
		t.Assert(l, -1)
	})
}

func Test_ComparatorFloat32OfFixed(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorFloat32(0.1, 0.1), 0)
		t.Assert(dutil.ComparatorFloat32(1.1, 2.1), -1)
		t.Assert(dutil.ComparatorFloat32(2.1, 1.1), 1)
	})
}

func Test_ComparatorFloat64OfFixed(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		t.Assert(dutil.ComparatorFloat64(0.1, 0.1), 0)
		t.Assert(dutil.ComparatorFloat64(1.1, 2.1), -1)
		t.Assert(dutil.ComparatorFloat64(2.1, 1.1), 1)
	})
}
