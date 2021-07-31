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

func Test_Bytes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewBytes([]byte("abc"))
		iClone := i.Clone()
		t.AssertEQ(iClone.Set([]byte("123")), []byte("abc"))
		t.AssertEQ(iClone.Val(), []byte("123"))

		//空参测试
		i1 := dtype.NewBytes()
		t.AssertEQ(i1.Val(), nil)
	})
}

func Test_Bytes_JSON(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		b := []byte("i love gf")
		i := dtype.NewBytes(b)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := dtype.NewBytes()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), b)
	})
}

func Test_Bytes_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *dtype.Bytes
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
