// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package derror_test

import (
	"errors"
	"github.com/osgochina/donkeygo/errors/derror"
	"testing"
)

var (
	// base error for benchmark testing of Wrap* functions.
	baseError = errors.New("test")
)

func Benchmark_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.New("test")
	}
}

func Benchmark_Newf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.Newf("%s", "test")
	}
}

func Benchmark_Wrap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.Wrap(baseError, "test")
	}
}

func Benchmark_Wrapf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.Wrapf(baseError, "%s", "test")
	}
}

func Benchmark_NewSkip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.NewSkip(1, "test")
	}
}

func Benchmark_NewSkipf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.NewSkipf(1, "%s", "test")
	}
}

func Benchmark_NewCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.NewCode(500, "test")
	}
}

func Benchmark_NewCodef(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.NewCodef(500, "%s", "test")
	}
}

func Benchmark_NewCodeSkip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.NewCodeSkip(1, 500, "test")
	}
}

func Benchmark_NewCodeSkipf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.NewCodeSkipf(1, 500, "%s", "test")
	}
}

func Benchmark_WrapCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.WrapCode(500, baseError, "test")
	}
}

func Benchmark_WrapCodef(b *testing.B) {
	for i := 0; i < b.N; i++ {
		derror.WrapCodef(500, baseError, "test")
	}
}
