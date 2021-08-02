// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dutil_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

func Test_StructToSlice(t *testing.T) {
	type A struct {
		K1 int
		K2 string
	}
	dtest.C(t, func(t *dtest.T) {
		a := &A{
			K1: 1,
			K2: "v2",
		}
		s := dutil.StructToSlice(a)
		t.Assert(len(s), 4)
		t.AssertIN(s[0], d.Slice{"K1", "K2", 1, "v2"})
		t.AssertIN(s[1], d.Slice{"K1", "K2", 1, "v2"})
		t.AssertIN(s[2], d.Slice{"K1", "K2", 1, "v2"})
		t.AssertIN(s[3], d.Slice{"K1", "K2", 1, "v2"})
	})
	dtest.C(t, func(t *dtest.T) {
		s := dutil.StructToSlice(1)
		t.Assert(s, nil)
	})
}
