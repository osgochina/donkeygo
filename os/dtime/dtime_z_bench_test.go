// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtime_test

import (
	"github.com/osgochina/donkeygo/os/dtime"
	"testing"
	"time"
)

func Benchmark_Timestamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.Timestamp()
	}
}

func Benchmark_TimestampMilli(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.TimestampMilli()
	}
}

func Benchmark_TimestampMicro(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.TimestampMicro()
	}
}

func Benchmark_TimestampNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.TimestampNano()
	}
}

func Benchmark_StrToTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.StrToTime("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_StrToTime_Format(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.StrToTime("2018-02-09 20:46:17.897", "Y-m-d H:i:su")
	}
}

func Benchmark_StrToTime_Layout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.StrToTimeLayout("2018-02-09T20:46:17.897Z", time.RFC3339)
	}
}

func Benchmark_ParseTimeFromContent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.ParseTimeFromContent("2018-02-09T20:46:17.897Z")
	}
}

func Benchmark_NewFromTimeStamp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.NewFromTimeStamp(1542674930)
	}
}

func Benchmark_Date(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.Date()
	}
}

func Benchmark_Datetime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.Datetime()
	}
}

func Benchmark_SetTimeZone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dtime.SetTimeZone("Asia/Shanghai")
	}
}
