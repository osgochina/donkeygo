// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar_test

import (
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

func TestVar_Struct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type StTest struct {
			Test int
		}

		Kv := make(map[string]int, 1)
		Kv["Test"] = 100

		testObj := &StTest{}

		objOne := dvar.New(Kv, true)

		objOne.Struct(testObj)

		t.Assert(testObj.Test, Kv["Test"])
	})
	dtest.C(t, func(t *dtest.T) {
		type StTest struct {
			Test int8
		}
		o := &StTest{}
		v := dvar.New(d.Slice{"Test", "-25"})
		v.Struct(o)
		t.Assert(o.Test, -25)
	})
}

func TestVar_Var_Attribute_Struct(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid  int
			Name string
		}
		user := new(User)
		err := dconv.Struct(
			d.Map{
				"uid":  dvar.New(1),
				"name": dvar.New("john"),
			}, user)
		t.Assert(err, nil)
		t.Assert(user.Uid, 1)
		t.Assert(user.Name, "john")
	})
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid  int
			Name string
		}
		var user *User
		err := dconv.Struct(
			d.Map{
				"uid":  dvar.New(1),
				"name": dvar.New("john"),
			}, &user)
		t.Assert(err, nil)
		t.Assert(user.Uid, 1)
		t.Assert(user.Name, "john")
	})
}
