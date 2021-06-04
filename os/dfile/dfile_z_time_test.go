// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dfile_test

import (
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/test/dtest"
	"os"
	"testing"
	"time"
)

func Test_MTime(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {

		var (
			file1   = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		t.Assert(err, nil)

		t.Assert(dfile.MTime(testpath()+file1), fileobj.ModTime())
		t.Assert(dfile.MTime(""), "")
	})
}

func Test_MTimeMillisecond(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			file1   = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		t.Assert(err, nil)

		time.Sleep(time.Millisecond * 100)
		t.AssertGE(
			dfile.MTimestampMilli(testpath()+file1),
			fileobj.ModTime().UnixNano()/1000000,
		)
		t.Assert(dfile.MTimestampMilli(""), -1)
	})
}
