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

func Test_Unsafe(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := "I love 小泽玛利亚"
		t.AssertEQ(dconv.UnsafeStrToBytes(s), []byte(s))
	})

	dtest.C(t, func(t *dtest.T) {
		b := []byte("I love 小泽玛利亚")
		t.AssertEQ(dconv.UnsafeBytesToStr(b), string(b))
	})
}
