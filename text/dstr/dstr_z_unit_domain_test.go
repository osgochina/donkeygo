// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package dstr_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"testing"
)

func Test_IsSubDomain(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		main := "goframe.org"
		t.Assert(dstr.IsSubDomain("goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.s.goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.s.goframe.org:8080", main), true)
		t.Assert(dstr.IsSubDomain("johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.s.johng.cn", main), false)
	})
	dtest.C(t, func(t *dtest.T) {
		main := "*.goframe.org"
		t.Assert(dstr.IsSubDomain("goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.goframe.org:80", main), true)
		t.Assert(dstr.IsSubDomain("s.s.goframe.org", main), false)
		t.Assert(dstr.IsSubDomain("johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.s.johng.cn", main), false)
	})
	dtest.C(t, func(t *dtest.T) {
		main := "*.*.goframe.org"
		t.Assert(dstr.IsSubDomain("goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.s.goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.s.goframe.org:8000", main), true)
		t.Assert(dstr.IsSubDomain("s.s.s.goframe.org", main), false)
		t.Assert(dstr.IsSubDomain("johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.s.johng.cn", main), false)
	})
	dtest.C(t, func(t *dtest.T) {
		main := "*.*.goframe.org:8080"
		t.Assert(dstr.IsSubDomain("goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.s.goframe.org", main), true)
		t.Assert(dstr.IsSubDomain("s.s.goframe.org:8000", main), true)
		t.Assert(dstr.IsSubDomain("s.s.s.goframe.org", main), false)
		t.Assert(dstr.IsSubDomain("johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.johng.cn", main), false)
		t.Assert(dstr.IsSubDomain("s.s.johng.cn", main), false)
	})
}
