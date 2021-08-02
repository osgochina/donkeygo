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
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

func Test_Map_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		m1 := map[string]string{
			"k": "v",
		}
		m2 := map[int]string{
			3: "v",
		}
		m3 := map[float64]float32{
			1.22: 3.1,
		}
		t.Assert(dconv.Map(m1), d.Map{
			"k": "v",
		})
		t.Assert(dconv.Map(m2), d.Map{
			"3": "v",
		})
		t.Assert(dconv.Map(m3), d.Map{
			"1.22": "3.1",
		})
	})
}

func Test_Map_Slice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		slice1 := d.Slice{"1", "2", "3", "4"}
		slice2 := d.Slice{"1", "2", "3"}
		slice3 := d.Slice{}
		t.Assert(dconv.Map(slice1), d.Map{
			"1": "2",
			"3": "4",
		})
		t.Assert(dconv.Map(slice2), d.Map{
			"1": "2",
			"3": nil,
		})
		t.Assert(dconv.Map(slice3), d.Map{})
	})
}

func Test_Maps_Basic(t *testing.T) {
	params := d.Slice{
		d.Map{"id": 100, "name": "john"},
		d.Map{"id": 200, "name": "smith"},
	}
	dtest.C(t, func(t *dtest.T) {
		list := dconv.Maps(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
}

func Test_Maps_JsonStr(t *testing.T) {
	jsonStr := `[{"id":100, "name":"john"},{"id":200, "name":"smith"}]`
	dtest.C(t, func(t *dtest.T) {
		list := dconv.Maps(jsonStr)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
}

func Test_Map_StructWithdconvTag(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `dconv:"-"`
			NickName string `dconv:"nickname, omitempty"`
			Pass1    string `dconv:"password1"`
			Pass2    string `dconv:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := dconv.Map(user1)
		map2 := dconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithJsonTag(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `json:"-"`
			NickName string `json:"nickname,omitempty"`
			Pass1    string `json:"password1"`
			Pass2    string `json:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := dconv.Map(user1)
		map2 := dconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithCTag(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `c:"-"`
			NickName string `c:"nickname, omitempty"`
			Pass1    string `c:"password1"`
			Pass2    string `c:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := dconv.Map(user1)
		map2 := dconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	dtest.C(t, func(t *dtest.T) {
		user := &User{1, "john"}
		t.Assert(dconv.Map(user), d.Map{"Id": 1})
	})
}

func Test_Map_Embedded(t *testing.T) {
	type Base struct {
		Id int
	}
	type User struct {
		Base
		Name string
	}
	type UserDetail struct {
		User
		Brief string
	}
	dtest.C(t, func(t *dtest.T) {
		user := &User{}
		user.Id = 1
		user.Name = "john"

		m := dconv.Map(user)
		t.Assert(len(m), 2)
		t.Assert(m["Id"], user.Id)
		t.Assert(m["Name"], user.Name)
	})
	dtest.C(t, func(t *dtest.T) {
		user := &UserDetail{}
		user.Id = 1
		user.Name = "john"
		user.Brief = "john guo"

		m := dconv.Map(user)
		t.Assert(len(m), 3)
		t.Assert(m["Id"], user.Id)
		t.Assert(m["Name"], user.Name)
		t.Assert(m["Brief"], user.Brief)
	})
}

func Test_Map_Embedded2(t *testing.T) {
	type Ids struct {
		Id  int `c:"id"`
		Uid int `c:"uid"`
	}
	type Base struct {
		Ids
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := dconv.Map(user)
		t.Assert(m["id"], "100")
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], "2019")
	})
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := dconv.MapDeep(user)
		t.Assert(m["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], user.CreateTime)
	})
}

func Test_MapDeep2(t *testing.T) {
	type A struct {
		F string
		G string
	}

	type B struct {
		A
		H string
	}

	type C struct {
		A A
		F string
	}

	type D struct {
		I A
		F string
	}

	dtest.C(t, func(t *dtest.T) {
		b := new(B)
		c := new(C)
		d := new(D)
		mb := dconv.MapDeep(b)
		mc := dconv.MapDeep(c)
		md := dconv.MapDeep(d)
		t.Assert(dutil.MapContains(mb, "F"), true)
		t.Assert(dutil.MapContains(mb, "G"), true)
		t.Assert(dutil.MapContains(mb, "H"), true)
		t.Assert(dutil.MapContains(mc, "A"), true)
		t.Assert(dutil.MapContains(mc, "F"), true)
		t.Assert(dutil.MapContains(mc, "G"), false)
		t.Assert(dutil.MapContains(md, "F"), true)
		t.Assert(dutil.MapContains(md, "I"), true)
		t.Assert(dutil.MapContains(md, "H"), false)
		t.Assert(dutil.MapContains(md, "G"), false)
	})
}

func Test_MapDeep3(t *testing.T) {
	type Base struct {
		Id   int    `c:"id"`
		Date string `c:"date"`
	}
	type User struct {
		UserBase Base   `c:"base"`
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}

	dtest.C(t, func(t *dtest.T) {
		user := &User{
			UserBase: Base{
				Id:   1,
				Date: "2019-10-01",
			},
			Passport: "john",
			Password: "123456",
			Nickname: "JohnGuo",
		}
		m := dconv.MapDeep(user)
		t.Assert(m, d.Map{
			"base": d.Map{
				"id":   user.UserBase.Id,
				"date": user.UserBase.Date,
			},
			"passport": user.Passport,
			"password": user.Password,
			"nickname": user.Nickname,
		})
	})

	dtest.C(t, func(t *dtest.T) {
		user := &User{
			UserBase: Base{
				Id:   1,
				Date: "2019-10-01",
			},
			Passport: "john",
			Password: "123456",
			Nickname: "JohnGuo",
		}
		m := dconv.Map(user)
		t.Assert(m, d.Map{
			"base":     user.UserBase,
			"passport": user.Passport,
			"password": user.Password,
			"nickname": user.Nickname,
		})
	})
}

func Test_MapDeepWithAttributeTag(t *testing.T) {
	type Ids struct {
		Id  int `c:"id"`
		Uid int `c:"uid"`
	}
	type Base struct {
		Ids        `json:"ids"`
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base     `json:"base"`
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := dconv.Map(user)
		t.Assert(m["id"], "")
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], "")
	})
	dtest.C(t, func(t *dtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := dconv.MapDeep(user)
		t.Assert(m["base"].(map[string]interface{})["ids"].(map[string]interface{})["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["base"].(map[string]interface{})["create_time"], user.CreateTime)
	})
}
