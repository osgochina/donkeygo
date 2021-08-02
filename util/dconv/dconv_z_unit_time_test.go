// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
	"time"
)

func Test_Time(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.AssertEQ(dconv.Duration(""), time.Duration(int64(0)))
		t.AssertEQ(dconv.GTime(""), dtime.New())
	})

	dtest.C(t, func(t *dtest.T) {
		s := "2011-10-10 01:02:03.456"
		t.AssertEQ(dconv.GTime(s), dtime.NewFromStr(s))
		t.AssertEQ(dconv.Time(s), dtime.NewFromStr(s).Time)
		t.AssertEQ(dconv.Duration(100), 100*time.Nanosecond)
	})
	dtest.C(t, func(t *dtest.T) {
		s := "01:02:03.456"
		t.AssertEQ(dconv.GTime(s).Hour(), 1)
		t.AssertEQ(dconv.GTime(s).Minute(), 2)
		t.AssertEQ(dconv.GTime(s).Second(), 3)
		t.AssertEQ(dconv.GTime(s), dtime.NewFromStr(s))
		t.AssertEQ(dconv.Time(s), dtime.NewFromStr(s).Time)
	})
	dtest.C(t, func(t *dtest.T) {
		s := "0000-01-01 01:02:03"
		t.AssertEQ(dconv.GTime(s).Year(), 0)
		t.AssertEQ(dconv.GTime(s).Month(), 1)
		t.AssertEQ(dconv.GTime(s).Day(), 1)
		t.AssertEQ(dconv.GTime(s).Hour(), 1)
		t.AssertEQ(dconv.GTime(s).Minute(), 2)
		t.AssertEQ(dconv.GTime(s).Second(), 3)
		t.AssertEQ(dconv.GTime(s), dtime.NewFromStr(s))
		t.AssertEQ(dconv.Time(s), dtime.NewFromStr(s).Time)
	})
	dtest.C(t, func(t *dtest.T) {
		t1 := dtime.NewFromStr("2021-05-21 05:04:51.206547+00")
		t2 := dconv.GTime(dvar.New(t1))
		t3 := dvar.New(t1).GTime()
		t.AssertEQ(t1, t2)
		t.AssertEQ(t1.Local(), t2.Local())
		t.AssertEQ(t1, t3)
		t.AssertEQ(t1.Local(), t3.Local())
	})
}

func Test_Time_Slice_Attribute(t *testing.T) {
	type SelectReq struct {
		Arr []*dtime.Time
		One *dtime.Time
	}
	dtest.C(t, func(t *dtest.T) {
		var s *SelectReq
		err := dconv.Struct(g.Map{
			"arr": g.Slice{"2021-01-12 12:34:56", "2021-01-12 12:34:57"},
			"one": "2021-01-12 12:34:58",
		}, &s)
		t.AssertNil(err)
		t.Assert(s.One, "2021-01-12 12:34:58")
		t.Assert(s.Arr[0], "2021-01-12 12:34:56")
		t.Assert(s.Arr[1], "2021-01-12 12:34:57")
	})
}
