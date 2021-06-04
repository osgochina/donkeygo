// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dfile_test

import (
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_Copy(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(dfile.Copy(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(dfile.IsFile(testpath()+topath), true)
		t.AssertNE(dfile.Copy("", ""), nil)
	})
}

func Test_CopyFile(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(dfile.CopyFile(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(dfile.IsFile(testpath()+topath), true)
		t.AssertNE(dfile.CopyFile("", ""), nil)
	})
	// Content replacement.
	dtest.C(t, func(t *dtest.T) {
		src := dfile.TempDir(dtime.TimestampNanoStr())
		dst := dfile.TempDir(dtime.TimestampNanoStr())
		srcContent := "1"
		dstContent := "1"
		t.Assert(dfile.PutContents(src, srcContent), nil)
		t.Assert(dfile.PutContents(dst, dstContent), nil)
		t.Assert(dfile.GetContents(src), srcContent)
		t.Assert(dfile.GetContents(dst), dstContent)

		t.Assert(dfile.CopyFile(src, dst), nil)
		t.Assert(dfile.GetContents(src), srcContent)
		t.Assert(dfile.GetContents(dst), srcContent)
	})
}

func Test_CopyDir(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			dirPath1 = "/test-copy-dir1"
			dirPath2 = "/test-copy-dir2"
		)
		haveList := []string{
			"t1.txt",
			"t2.txt",
		}
		createDir(dirPath1)
		for _, v := range haveList {
			t.Assert(createTestFile(dirPath1+"/"+v, ""), nil)
		}
		defer delTestFiles(dirPath1)

		var (
			yfolder  = testpath() + dirPath1
			tofolder = testpath() + dirPath2
		)

		if dfile.IsDir(tofolder) {
			t.Assert(dfile.Remove(tofolder), nil)
			t.Assert(dfile.Remove(""), nil)
		}

		t.Assert(dfile.CopyDir(yfolder, tofolder), nil)
		defer delTestFiles(tofolder)

		t.Assert(dfile.IsDir(yfolder), true)

		for _, v := range haveList {
			t.Assert(dfile.IsFile(yfolder+"/"+v), true)
		}

		t.Assert(dfile.IsDir(tofolder), true)

		for _, v := range haveList {
			t.Assert(dfile.IsFile(tofolder+"/"+v), true)
		}

		t.Assert(dfile.Remove(tofolder), nil)
		t.Assert(dfile.Remove(""), nil)
	})
	// Content replacement.
	dtest.C(t, func(t *dtest.T) {
		src := dfile.TempDir(dtime.TimestampNanoStr(), dtime.TimestampNanoStr())
		dst := dfile.TempDir(dtime.TimestampNanoStr(), dtime.TimestampNanoStr())
		defer func() {
			dfile.Remove(src)
			dfile.Remove(dst)
		}()
		srcContent := "1"
		dstContent := "1"
		t.Assert(dfile.PutContents(src, srcContent), nil)
		t.Assert(dfile.PutContents(dst, dstContent), nil)
		t.Assert(dfile.GetContents(src), srcContent)
		t.Assert(dfile.GetContents(dst), dstContent)

		err := dfile.CopyDir(dfile.Dir(src), dfile.Dir(dst))
		t.Assert(err, nil)
		t.Assert(dfile.GetContents(src), srcContent)
		t.Assert(dfile.GetContents(dst), srcContent)
	})
}
