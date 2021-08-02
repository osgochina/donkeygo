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

func Test_OctStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.OctStr(`\346\200\241`), "怡")
	})
}
