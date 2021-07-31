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

func Test_Bool(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewBool(true)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(false), true)
		t.AssertEQ(iClone.Val(), false)

		i1 := dtype.NewBool(false)
		iClone1 := i1.Clone()
		t.AssertEQ(iClone1.Set(true), false)
		t.AssertEQ(iClone1.Val(), true)

		//空参测试
		i2 := dtype.NewBool()
		t.AssertEQ(i2.Val(), false)
	})
}

func Test_Bool_JSON(t *testing.T) {
	// Marshal
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewBool(true)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)
	})
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewBool(false)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)
	})
	// Unmarshal
	dtest.C(t, func(t *dtest.T) {
		var err error
		i := dtype.NewBool()
		err = json.UnmarshalUseNumber([]byte("true"), &i)
		t.Assert(err, nil)
		t.Assert(i.Val(), true)
		err = json.UnmarshalUseNumber([]byte("false"), &i)
		t.Assert(err, nil)
		t.Assert(i.Val(), false)
		err = json.UnmarshalUseNumber([]byte("1"), &i)
		t.Assert(err, nil)
		t.Assert(i.Val(), true)
		err = json.UnmarshalUseNumber([]byte("0"), &i)
		t.Assert(err, nil)
		t.Assert(i.Val(), false)
	})

	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewBool(true)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := dtype.NewBool()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), i.Val())
	})
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewBool(false)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := dtype.NewBool()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), i.Val())
	})
}

func Test_Bool_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *dtype.Bool
	}
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "true",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), true)
	})
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "false",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), false)
	})
}
