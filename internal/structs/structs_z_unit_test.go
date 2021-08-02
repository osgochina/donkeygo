// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/internal/structs"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user User
		m, _ := structs.TagMapName(user, []string{"params"})
		t.Assert(m, d.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"params"})
		t.Assert(m, d.Map{"name": "Name", "pass": "Pass"})

		m, _ = structs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, d.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, d.Map{"name": "Name", "pass1": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, d.Map{"name": "Name", "pass2": "Pass"})
	})

	dtest.C(t, func(t *dtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithBase struct {
			Id   int
			Name string
			Base `params:"base"`
		}
		user := new(UserWithBase)
		m, _ := structs.TagMapName(user, []string{"params"})
		t.Assert(m, d.Map{
			"base":      "Base",
			"password1": "Pass1",
			"password2": "Pass2",
		})
	})

	dtest.C(t, func(t *dtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithEmbeddedAttribute struct {
			Id   int
			Name string
			Base
		}
		type UserWithoutEmbeddedAttribute struct {
			Id   int
			Name string
			Pass Base
		}
		user1 := new(UserWithEmbeddedAttribute)
		user2 := new(UserWithoutEmbeddedAttribute)
		m, _ := structs.TagMapName(user1, []string{"params"})
		t.Assert(m, d.Map{"password1": "Pass1", "password2": "Pass2"})
		m, _ = structs.TagMapName(user2, []string{"params"})
		t.Assert(m, d.Map{})
	})
}

func Test_StructOfNilPointer(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := structs.TagMapName(user, []string{"params"})
		t.Assert(m, d.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"params"})
		t.Assert(m, d.Map{"name": "Name", "pass": "Pass"})

		m, _ = structs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, d.Map{"name": "Name", "pass": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, d.Map{"name": "Name", "pass1": "Pass"})
		m, _ = structs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, d.Map{"name": "Name", "pass2": "Pass"})
	})
}

func Test_FieldMap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := structs.FieldMap(user, []string{"params"})
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, false)
		_, ok = m["name"]
		t.Assert(ok, true)
		_, ok = m["Pass"]
		t.Assert(ok, false)
		_, ok = m["pass"]
		t.Assert(ok, true)
	})
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := structs.FieldMap(user, nil)
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, true)
		_, ok = m["name"]
		t.Assert(ok, false)
		_, ok = m["Pass"]
		t.Assert(ok, true)
		_, ok = m["pass"]
		t.Assert(ok, false)
	})
}

//
//func Test_StructType(t *testing.T) {
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Name string
//		}
//		type A struct {
//			B
//		}
//		r, err := structs.StructType(new(A))
//		t.AssertNil(err)
//		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.A`)
//	})
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Name string
//		}
//		type A struct {
//			B
//		}
//		r, err := structs.StructType(new(A).B)
//		t.AssertNil(err)
//		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
//	})
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Name string
//		}
//		type A struct {
//			*B
//		}
//		r, err := structs.StructType(new(A).B)
//		t.AssertNil(err)
//		t.Assert(r.String(), `structs_test.B`)
//	})
//	// Error.
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Name string
//		}
//		type A struct {
//			*B
//			Id int
//		}
//		_, err := structs.StructType(new(A).Id)
//		t.AssertNE(err, nil)
//	})
//}
//
//func Test_StructTypeBySlice(t *testing.T) {
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Name string
//		}
//		type A struct {
//			Array []*B
//		}
//		r, err := structs.StructType(new(A).Array)
//		t.AssertNil(err)
//		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
//	})
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Name string
//		}
//		type A struct {
//			Array []B
//		}
//		r, err := structs.StructType(new(A).Array)
//		t.AssertNil(err)
//		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
//	})
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Name string
//		}
//		type A struct {
//			Array *[]B
//		}
//		r, err := structs.StructType(new(A).Array)
//		t.AssertNil(err)
//		t.Assert(r.Signature(), `github.com/gogf/gf/internal/structs_test/structs_test.B`)
//	})
//}
//
//func TestType_FieldKeys(t *testing.T) {
//	dtest.C(t, func(t *dtest.T) {
//		type B struct {
//			Id   int
//			Name string
//		}
//		type A struct {
//			Array []*B
//		}
//		r, err := structs.StructType(new(A).Array)
//		t.AssertNil(err)
//		t.Assert(r.FieldKeys(), d.Slice{"Id", "Name"})
//	})
//}
