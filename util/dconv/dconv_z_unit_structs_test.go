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

func Test_Structs_WithTag(t *testing.T) {
	type User struct {
		Uid      int    `json:"id"`
		NickName string `json:"name"`
	}
	dtest.C(t, func(t *dtest.T) {
		var users []User
		params := d.Slice{
			d.Map{
				"id":   1,
				"name": "name1",
			},
			d.Map{
				"id":   2,
				"name": "name2",
			},
		}
		err := dconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
	dtest.C(t, func(t *dtest.T) {
		var users []*User
		params := d.Slice{
			d.Map{
				"id":   1,
				"name": "name1",
			},
			d.Map{
				"id":   2,
				"name": "name2",
			},
		}
		err := dconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
}

func Test_Structs_WithoutTag(t *testing.T) {
	type User struct {
		Uid      int
		NickName string
	}
	dtest.C(t, func(t *dtest.T) {
		var users []User
		params := d.Slice{
			d.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			d.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := dconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
	dtest.C(t, func(t *dtest.T) {
		var users []*User
		params := d.Slice{
			d.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			d.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := dconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
}

func Test_Structs_SliceParameter(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid      int
			NickName string
		}
		var users []User
		params := d.Slice{
			d.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			d.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := dconv.Structs(params, users)
		t.AssertNE(err, nil)
	})
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid      int
			NickName string
		}
		type A struct {
			Users []User
		}
		var a A
		params := d.Slice{
			d.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			d.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := dconv.Structs(params, a.Users)
		t.AssertNE(err, nil)
	})
}

func Test_Structs_DirectReflectSet(t *testing.T) {
	type A struct {
		Id   int
		Name string
	}
	dtest.C(t, func(t *dtest.T) {
		var (
			a = []*A{
				{Id: 1, Name: "john"},
				{Id: 2, Name: "smith"},
			}
			b []*A
		)
		err := dconv.Structs(a, &b)
		t.Assert(err, nil)
		t.AssertEQ(a, b)
	})
	dtest.C(t, func(t *dtest.T) {
		var (
			a = []A{
				{Id: 1, Name: "john"},
				{Id: 2, Name: "smith"},
			}
			b []A
		)
		err := dconv.Structs(a, &b)
		t.Assert(err, nil)
		t.AssertEQ(a, b)
	})
}

func Test_Structs_IntSliceAttribute(t *testing.T) {
	type A struct {
		Id []int
	}
	type B struct {
		*A
		Name string
	}
	dtest.C(t, func(t *dtest.T) {
		var (
			array []*B
		)
		err := dconv.Structs(d.Slice{
			d.Map{"id": nil, "name": "john"},
			d.Map{"id": nil, "name": "smith"},
		}, &array)
		t.Assert(err, nil)
		t.Assert(len(array), 2)
		t.Assert(array[0].Name, "john")
		t.Assert(array[1].Name, "smith")
	})
}
