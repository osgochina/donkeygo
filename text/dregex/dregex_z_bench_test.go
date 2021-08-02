// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package dregex_test

import (
	"github.com/osgochina/donkeygo/text/dregex"
	"regexp"
	"testing"
)

var pattern = `(\w+).+\-\-\s*(.+)`
var src = `DK is best! -- John`

func Benchmark_DK_IsMatchString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dregex.IsMatchString(pattern, src)
	}
}

func Benchmark_DK_MatchString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dregex.MatchString(pattern, src)
	}
}

func Benchmark_Compile(b *testing.B) {
	var wcdRegexp = regexp.MustCompile(pattern)
	for i := 0; i < b.N; i++ {
		wcdRegexp.MatchString(src)
	}
}

func Benchmark_Compile_Actual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wcdRegexp := regexp.MustCompile(pattern)
		wcdRegexp.MatchString(src)
	}
}
