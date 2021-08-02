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

func Test_Scan_StructStructs(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `dconv:"password1"`
		Pass2 string `dconv:"password2"`
	}
	dtest.C(t, func(t *dtest.T) {
		var (
			user   = new(User)
			params = d.Map{
				"uid":   1,
				"name":  "john",
				"PASS1": "123",
				"PASS2": "456",
			}
		)
		err := dconv.Scan(params, user)
		t.Assert(err, nil)
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	dtest.C(t, func(t *dtest.T) {
		var (
			users  []User
			params = d.Slice{
				d.Map{
					"uid":   1,
					"name":  "john1",
					"PASS1": "111",
					"PASS2": "222",
				},
				d.Map{
					"uid":   2,
					"name":  "john2",
					"PASS1": "333",
					"PASS2": "444",
				},
			}
		)
		err := dconv.Scan(params, &users)
		t.AssertNil(err)
		t.Assert(users, d.Slice{
			&User{
				Uid:   1,
				Name:  "john1",
				Pass1: "111",
				Pass2: "222",
			},
			&User{
				Uid:   2,
				Name:  "john2",
				Pass1: "333",
				Pass2: "444",
			},
		})
	})
}

func Test_Scan_StructStr(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `dconv:"password1"`
		Pass2 string `dconv:"password2"`
	}
	dtest.C(t, func(t *dtest.T) {
		var (
			user   = new(User)
			params = `{"uid":1,"name":"john", "pass1":"123","pass2":"456"}`
		)
		err := dconv.Scan(params, user)
		t.Assert(err, nil)
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	dtest.C(t, func(t *dtest.T) {
		var (
			users  []User
			params = `[
{"uid":1,"name":"john1", "pass1":"111","pass2":"222"},
{"uid":2,"name":"john2", "pass1":"333","pass2":"444"}
]`
		)
		err := dconv.Scan(params, &users)
		t.Assert(err, nil)
		t.Assert(users, d.Slice{
			&User{
				Uid:   1,
				Name:  "john1",
				Pass1: "111",
				Pass2: "222",
			},
			&User{
				Uid:   2,
				Name:  "john2",
				Pass1: "333",
				Pass2: "444",
			},
		})
	})
}

func Test_Scan_Map(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var m map[string]string
		data := d.Map{
			"k1": "v1",
			"k2": "v2",
		}
		err := dconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(data, m)
	})
	dtest.C(t, func(t *dtest.T) {
		var m map[int]int
		data := d.Map{
			"1": "11",
			"2": "22",
		}
		err := dconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(data, m)
	})
	// json string parameter.
	dtest.C(t, func(t *dtest.T) {
		var m map[string]string
		data := `{"k1":"v1","k2":"v2"}`
		err := dconv.Scan(data, &m)
		t.AssertNil(err)
		t.Assert(m, d.Map{
			"k1": "v1",
			"k2": "v2",
		})
	})
}

func Test_Scan_Maps(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var maps []map[string]string
		data := d.Slice{
			d.Map{
				"k1": "v1",
				"k2": "v2",
			},
			d.Map{
				"k3": "v3",
				"k4": "v4",
			},
		}
		err := dconv.Scan(data, &maps)
		t.AssertNil(err)
		t.Assert(data, maps)
	})
	// json string parameter.
	dtest.C(t, func(t *dtest.T) {
		var maps []map[string]string
		data := `[{"k1":"v1","k2":"v2"},{"k3":"v3","k4":"v4"}]`
		err := dconv.Scan(data, &maps)
		t.AssertNil(err)
		t.Assert(maps, d.Slice{
			d.Map{
				"k1": "v1",
				"k2": "v2",
			},
			d.Map{
				"k3": "v3",
				"k4": "v4",
			},
		})
	})
}
