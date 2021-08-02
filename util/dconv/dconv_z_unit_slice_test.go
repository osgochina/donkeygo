// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

func Test_Slice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := 123.456
		t.AssertEQ(dconv.Bytes("123"), []byte("123"))
		t.AssertEQ(dconv.Bytes([]interface{}{1}), []byte{1})
		t.AssertEQ(dconv.Bytes([]interface{}{300}), []byte("[300]"))
		t.AssertEQ(dconv.Strings(value), []string{"123.456"})
		t.AssertEQ(dconv.Ints(value), []int{123})
		t.AssertEQ(dconv.Floats(value), []float64{123.456})
		t.AssertEQ(dconv.Interfaces(value), []interface{}{123.456})
	})
	dtest.C(t, func(t *dtest.T) {
		s := []*dvar.Var{
			dvar.New(1),
			dvar.New(2),
		}
		t.AssertEQ(dconv.SliceInt64(s), []int64{1, 2})
	})
}

func Test_Slice_Empty(t *testing.T) {
	// Int.
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Ints(""), []int{})
		t.Assert(dconv.Ints(nil), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Int32s(""), []int32{})
		t.Assert(dconv.Int32s(nil), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Int64s(""), []int64{})
		t.Assert(dconv.Int64s(nil), nil)
	})
	// Uint.
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Uints(""), []uint{})
		t.Assert(dconv.Uints(nil), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Uint32s(""), []uint32{})
		t.Assert(dconv.Uint32s(nil), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Uint64s(""), []uint64{})
		t.Assert(dconv.Uint64s(nil), nil)
	})
	// Float.
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Floats(""), []float64{})
		t.Assert(dconv.Floats(nil), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Float32s(""), []float32{})
		t.Assert(dconv.Float32s(nil), nil)
	})
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Float64s(""), []float64{})
		t.Assert(dconv.Float64s(nil), nil)
	})
}

func Test_Strings(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		array := []*d.Var{
			d.NewVar(1),
			d.NewVar(2),
			d.NewVar(3),
		}
		t.AssertEQ(dconv.Strings(array), []string{"1", "2", "3"})
	})
}

func Test_Slice_Interfaces(t *testing.T) {
	// map
	dtest.C(t, func(t *dtest.T) {
		array := dconv.Interfaces(d.Map{
			"id":   1,
			"name": "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(d.Map)["id"], 1)
		t.Assert(array[0].(d.Map)["name"], "john")
	})
	// struct
	dtest.C(t, func(t *dtest.T) {
		type A struct {
			Id   int `json:"id"`
			Name string
		}
		array := dconv.Interfaces(&A{
			Id:   1,
			Name: "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(*A).Id, 1)
		t.Assert(array[0].(*A).Name, "john")
	})
}

func Test_Slice_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int    `json:"id"`
		name string `json:"name"`
	}
	dtest.C(t, func(t *dtest.T) {
		user := &User{1, "john"}
		array := dconv.Interfaces(user)
		t.Assert(len(array), 1)
		t.Assert(array[0].(*User).Id, 1)
		t.Assert(array[0].(*User).name, "john")
	})
}

func Test_Slice_Structs(t *testing.T) {
	type Base struct {
		Age int
	}
	type User struct {
		Id   int
		Name string
		Base
	}

	dtest.C(t, func(t *dtest.T) {
		users := make([]User, 0)
		params := []d.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := dconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, params[0]["id"])
		t.Assert(users[0].Name, params[0]["name"])
		t.Assert(users[0].Age, 18)

		t.Assert(users[1].Id, params[1]["id"])
		t.Assert(users[1].Name, params[1]["name"])
		t.Assert(users[1].Age, 20)
	})
}
