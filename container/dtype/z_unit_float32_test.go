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
	"math"
	"testing"
)

func Test_Float32(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewFloat32(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(0.1), float32(0))
		t.AssertEQ(iClone.Val(), float32(0.1))

		//空参测试
		i1 := dtype.NewFloat32()
		t.AssertEQ(i1.Val(), float32(0))
	})
}

func Test_Float32_JSON(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		v := float32(math.MaxFloat32)
		i := dtype.NewFloat32(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())

		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := dtype.NewFloat32()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), v)
	})
}

func Test_Float32_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *dtype.Float32
	}
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123.456",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123.456")
	})
}
