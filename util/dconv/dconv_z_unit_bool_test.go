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

type boolStruct struct {
}

func Test_Bool(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var any interface{} = nil
		t.AssertEQ(dconv.Bool(any), false)
		t.AssertEQ(dconv.Bool(false), false)
		t.AssertEQ(dconv.Bool(nil), false)
		t.AssertEQ(dconv.Bool(0), false)
		t.AssertEQ(dconv.Bool("0"), false)
		t.AssertEQ(dconv.Bool(""), false)
		t.AssertEQ(dconv.Bool("false"), false)
		t.AssertEQ(dconv.Bool("off"), false)
		t.AssertEQ(dconv.Bool([]byte{}), false)
		t.AssertEQ(dconv.Bool([]string{}), false)
		t.AssertEQ(dconv.Bool([]interface{}{}), false)
		t.AssertEQ(dconv.Bool([]map[int]int{}), false)

		t.AssertEQ(dconv.Bool("1"), true)
		t.AssertEQ(dconv.Bool("on"), true)
		t.AssertEQ(dconv.Bool(1), true)
		t.AssertEQ(dconv.Bool(123.456), true)
		t.AssertEQ(dconv.Bool(boolStruct{}), true)
		t.AssertEQ(dconv.Bool(&boolStruct{}), true)
	})
}
