// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtype_test

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

func Test_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewString("abc")
		iClone := i.Clone()
		t.AssertEQ(iClone.Set("123"), "abc")
		t.AssertEQ(iClone.Val(), "123")

		//空参测试
		i1 := dtype.NewString()
		t.AssertEQ(i1.Val(), "")
	})
}

func Test_String_JSON(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := "i love gf"
		i1 := dtype.NewString(s)
		b1, err1 := json.Marshal(i1)
		b2, err2 := json.Marshal(i1.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := dtype.NewString()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), s)
	})
}

func Test_String_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *dtype.String
	}
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123")
	})
}
