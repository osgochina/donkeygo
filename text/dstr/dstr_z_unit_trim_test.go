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

func Test_Trim(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Trim(" 123456\n "), "123456")
		t.Assert(dstr.Trim("#123456#;", "#;"), "123456")
	})
}

func Test_TrimStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimStr("gogo我爱gogo", "go"), "我爱")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimStr("gogo我爱gogo", "go", 1), "go我爱go")
		t.Assert(dstr.TrimStr("gogo我爱gogo", "go", 2), "我爱")
		t.Assert(dstr.TrimStr("gogo我爱gogo", "go", -1), "我爱")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimStr("啊我爱中国人啊", "啊"), "我爱中国人")
	})
}

func Test_TrimRight(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimRight(" 123456\n "), " 123456")
		t.Assert(dstr.TrimRight("#123456#;", "#;"), "#123456")
	})
}

func Test_TrimRightStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimRightStr("gogo我爱gogo", "go"), "gogo我爱")
		t.Assert(dstr.TrimRightStr("gogo我爱gogo", "go我爱gogo"), "go")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimRightStr("gogo我爱gogo", "go", 1), "gogo我爱go")
		t.Assert(dstr.TrimRightStr("gogo我爱gogo", "go", 2), "gogo我爱")
		t.Assert(dstr.TrimRightStr("gogo我爱gogo", "go", -1), "gogo我爱")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimRightStr("我爱中国人", "人"), "我爱中国")
		t.Assert(dstr.TrimRightStr("我爱中国人", "爱中国人"), "我")
	})
}

func Test_TrimLeft(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimLeft(" \r123456\n "), "123456\n ")
		t.Assert(dstr.TrimLeft("#;123456#;", "#;"), "123456#;")
	})
}

func Test_TrimLeftStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimLeftStr("gogo我爱gogo", "go"), "我爱gogo")
		t.Assert(dstr.TrimLeftStr("gogo我爱gogo", "gogo我爱go"), "go")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimLeftStr("gogo我爱gogo", "go", 1), "go我爱gogo")
		t.Assert(dstr.TrimLeftStr("gogo我爱gogo", "go", 2), "我爱gogo")
		t.Assert(dstr.TrimLeftStr("gogo我爱gogo", "go", -1), "我爱gogo")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimLeftStr("我爱中国人", "我爱"), "中国人")
		t.Assert(dstr.TrimLeftStr("我爱中国人", "我爱中国"), "人")
	})
}

func Test_TrimAll(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimAll("gogo我go\n爱gogo\n", "go"), "我爱")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimAll("gogo\n我go爱gogo", "go"), "我爱")
		t.Assert(dstr.TrimAll("gogo\n我go爱gogo\n", "go"), "我爱")
		t.Assert(dstr.TrimAll("gogo\n我go\n爱gogo", "go"), "我爱")
	})
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.TrimAll("啊我爱\n啊中国\n人啊", "啊"), "我爱中国人")
	})
}
