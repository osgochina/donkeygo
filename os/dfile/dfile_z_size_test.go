// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dfile_test

import (
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
)

func Test_Size(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			paths1 string = "/testfile_t1.txt"
			sizes  int64
		)

		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)

		sizes = dfile.Size(testpath() + paths1)
		t.Assert(sizes, 14)

		sizes = dfile.Size("")
		t.Assert(sizes, 0)

	})
}

func Test_StrToSize(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dfile.StrToSize("0.00B"), 0)
		t.Assert(dfile.StrToSize("16.00B"), 16)
		t.Assert(dfile.StrToSize("1.00K"), 1024)
		t.Assert(dfile.StrToSize("1.00KB"), 1024)
		t.Assert(dfile.StrToSize("1.00KiloByte"), 1024)
		t.Assert(dfile.StrToSize("15.26M"), dconv.Int64(15.26*1024*1024))
		t.Assert(dfile.StrToSize("15.26MB"), dconv.Int64(15.26*1024*1024))
		t.Assert(dfile.StrToSize("1.49G"), dconv.Int64(1.49*1024*1024*1024))
		t.Assert(dfile.StrToSize("1.49GB"), dconv.Int64(1.49*1024*1024*1024))
		t.Assert(dfile.StrToSize("8.73T"), dconv.Int64(8.73*1024*1024*1024*1024))
		t.Assert(dfile.StrToSize("8.73TB"), dconv.Int64(8.73*1024*1024*1024*1024))
		t.Assert(dfile.StrToSize("8.53P"), dconv.Int64(8.53*1024*1024*1024*1024*1024))
		t.Assert(dfile.StrToSize("8.53PB"), dconv.Int64(8.53*1024*1024*1024*1024*1024))
		t.Assert(dfile.StrToSize("8.01EB"), dconv.Int64(8.01*1024*1024*1024*1024*1024*1024))
	})
}

func Test_FormatSize(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dfile.FormatSize(0), "0.00B")
		t.Assert(dfile.FormatSize(16), "16.00B")

		t.Assert(dfile.FormatSize(1024), "1.00K")

		t.Assert(dfile.FormatSize(16000000), "15.26M")

		t.Assert(dfile.FormatSize(1600000000), "1.49G")

		t.Assert(dfile.FormatSize(9600000000000), "8.73T")
		t.Assert(dfile.FormatSize(9600000000000000), "8.53P")
	})
}

func Test_ReadableSize(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {

		var (
			paths1 string = "/testfile_t1.txt"
		)
		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)
		t.Assert(dfile.ReadableSize(testpath()+paths1), "14.00B")
		t.Assert(dfile.ReadableSize(""), "0.00B")

	})
}
