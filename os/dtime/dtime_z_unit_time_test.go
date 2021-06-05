// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtime_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	// time.Time
	dtest.C(t, func(t *dtest.T) {
		timeNow := time.Now()
		timeTemp := dtime.New(timeNow)
		t.Assert(timeTemp.Time.UnixNano(), timeNow.UnixNano())

		timeTemp1 := dtime.New()
		t.Assert(timeTemp1.Time, time.Time{})
	})
	// string
	dtest.C(t, func(t *dtest.T) {
		timeNow := dtime.Now()
		timeTemp := dtime.New(timeNow.String())
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), timeNow.Time.Format("2006-01-02 15:04:05"))
	})
	dtest.C(t, func(t *dtest.T) {
		timeNow := dtime.Now()
		timeTemp := dtime.New(timeNow.TimestampMicroStr())
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), timeNow.Time.Format("2006-01-02 15:04:05"))
	})
	// int64
	dtest.C(t, func(t *dtest.T) {
		timeNow := dtime.Now()
		timeTemp := dtime.New(timeNow.TimestampMicro())
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), timeNow.Time.Format("2006-01-02 15:04:05"))
	})
	// short datetime.
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.New("2021-2-9 08:01:21")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2021-02-09 08:01:21")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2021-02-09 08:01:21")
	})
}

func Test_Nil(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var t1 *dtime.Time
		t.Assert(t1.String(), "")
	})
	dtest.C(t, func(t *dtest.T) {
		var t1 dtime.Time
		t.Assert(t1.String(), "")
	})
}

func Test_NewFromStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2006-01-02 15:04:05")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := dtime.NewFromStr("2006.0102")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t1 := dtime.NewFromStr("2006-01-02 15:04:05")
		t.Assert(t1.String(), "2006-01-02 15:04:05")
		t.Assert(fmt.Sprintf("%s", t1), "2006-01-02 15:04:05")

		t2 := *t1
		t.Assert(t2.String(), "2006-01-02 15:04:05")
		t.Assert(fmt.Sprintf("{%s}", t2.String()), "{2006-01-02 15:04:05}")
	})
}

func Test_NewFromStrFormat(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStrFormat("2006-01-02 15:04:05", "Y-m-d H:i:s")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := dtime.NewFromStrFormat("2006-01-02 15:04:05", "aabbcc")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})

	dtest.C(t, func(t *dtest.T) {
		t1 := dtime.NewFromStrFormat("2019/2/1", "Y/n/j")
		t.Assert(t1.Format("Y-m-d"), "2019-02-01")

		t2 := dtime.NewFromStrFormat("2019/10/12", "Y/n/j")
		t.Assert(t2.Format("Y-m-d"), "2019-10-12")
	})
}

func Test_NewFromStrLayout(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStrLayout("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := dtime.NewFromStrLayout("2006-01-02 15:04:05", "aabbcc")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_NewFromTimeStamp(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromTimeStamp(1554459846000)
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2019-04-05 18:24:06")
		timeTemp1 := dtime.NewFromTimeStamp(0)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "0001-01-01 00:00:00")
	})
}

func Test_Time_Second(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		t.Assert(timeTemp.Second(), timeTemp.Time.Second())
	})
}

func Test_Time_Nanosecond(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		t.Assert(timeTemp.Nanosecond(), timeTemp.Time.Nanosecond())
	})
}

func Test_Time_Microsecond(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		t.Assert(timeTemp.Microsecond(), timeTemp.Time.Nanosecond()/1e3)
	})
}

func Test_Time_Millisecond(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		t.Assert(timeTemp.Millisecond(), timeTemp.Time.Nanosecond()/1e6)
	})
}

func Test_Time_String(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		t.Assert(timeTemp.String(), timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
}

func Test_Time_ISO8601(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		now := dtime.Now()
		t.Assert(now.ISO8601(), now.Format("c"))
	})
}

func Test_Time_RFC822(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		now := dtime.Now()
		t.Assert(now.RFC822(), now.Format("r"))
	})
}

func Test_Clone(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		timeTemp1 := timeTemp.Clone()
		t.Assert(timeTemp.Time.Unix(), timeTemp1.Time.Unix())
	})
}

func Test_ToTime(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		timeTemp1 := timeTemp.Time
		t.Assert(timeTemp.Time.UnixNano(), timeTemp1.UnixNano())
	})
}

func Test_Add(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2006-01-02 15:04:05")
		timeTemp = timeTemp.Add(time.Second)
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:06")
	})
}

func Test_ToZone(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		timeTemp, _ = timeTemp.ToZone("America/Los_Angeles")
		t.Assert(timeTemp.Time.Location().String(), "America/Los_Angeles")

		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			t.Error("test fail")
		}
		timeTemp = timeTemp.ToLocation(loc)
		t.Assert(timeTemp.Time.Location().String(), "Asia/Shanghai")

		timeTemp1, _ := timeTemp.ToZone("errZone")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_AddDate(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2006-01-02 15:04:05")
		timeTemp = timeTemp.AddDate(1, 2, 3)
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2007-03-05 15:04:05")
	})
}

func Test_UTC(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.UTC()
		t.Assert(timeTemp.UnixNano(), timeTemp1.UTC().UnixNano())
	})
}

func Test_Local(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.Local()
		t.Assert(timeTemp.UnixNano(), timeTemp1.Local().UnixNano())
	})
}

func Test_Round(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp = timeTemp.Round(time.Hour)
		t.Assert(timeTemp.UnixNano(), timeTemp1.Round(time.Hour).UnixNano())
	})
}

func Test_Truncate(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp = timeTemp.Truncate(time.Hour)
		t.Assert(timeTemp.UnixNano(), timeTemp1.Truncate(time.Hour).UnixNano())
	})
}

func Test_StartOfMinute(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfMinute()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 18:24:00")
	})
}

func Test_EndOfMinute(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfMinute()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 18:24:59")
	})
}

func Test_StartOfHour(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfHour()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 18:00:00")
	})
}

func Test_EndOfHour(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfHour()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 18:59:59")
	})
}

func Test_StartOfDay(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfDay()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 00:00:00")
	})
}

func Test_EndOfDay(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfDay()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 23:59:59")
	})
}

func Test_StartOfWeek(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfWeek()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-06 00:00:00")
	})
}

func Test_EndOfWeek(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfWeek()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 23:59:59")
	})
}

func Test_StartOfMonth(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfMonth()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-01 00:00:00")
	})
}

func Test_EndOfMonth(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfMonth()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-31 23:59:59")
	})
}

func Test_StartOfQuarter(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.StartOfQuarter()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-10-01 00:00:00")
	})
}

func Test_EndOfQuarter(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfQuarter()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-31 23:59:59")
	})
}

func Test_StartOfHalf(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.StartOfHalf()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-07-01 00:00:00")
	})
}

func Test_EndOfHalf(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfHalf()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-31 23:59:59")
	})
}

func Test_StartOfYear(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.StartOfYear()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-01-01 00:00:00")
	})
}

func Test_EndOfYear(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		timeTemp := dtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfYear()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-31 23:59:59")
	})
}
