// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

func Test_StructTag(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `orm:"password1"`
		Pass2 string `orm:"password2"`
	}
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		params1 := d.Map{
			"uid":       1,
			"Name":      "john",
			"password1": "123",
			"password2": "456",
		}
		if err := dconv.Struct(params1, user); err != nil {
			t.Error(err)
		}
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "",
			Pass2: "",
		})
	})
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		params1 := d.Map{
			"uid":       1,
			"Name":      "john",
			"password1": "123",
			"password2": "456",
		}
		if err := dconv.StructTag(params1, user, "orm"); err != nil {
			t.Error(err)
		}
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		params2 := d.Map{
			"uid":       2,
			"name":      "smith",
			"password1": "111",
			"password2": "222",
		}
		if err := dconv.StructTag(params2, user, "orm"); err != nil {
			t.Error(err)
		}
		t.Assert(user, &User{
			Uid:   2,
			Name:  "smith",
			Pass1: "111",
			Pass2: "222",
		})
	})
}
