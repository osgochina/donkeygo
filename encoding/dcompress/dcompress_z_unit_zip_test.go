// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dcompress_test

import (
	"bytes"
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/encoding/dcompress"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_ZipPath(t *testing.T) {
	// file
	dtest.C(t, func(t *dtest.T) {
		srcPath := ddebug.TestDataPath("zip", "path1", "1.txt")
		dstPath := ddebug.TestDataPath("zip", "zip.zip")

		t.Assert(dfile.Exists(dstPath), false)
		t.Assert(dcompress.ZipPath(srcPath, dstPath), nil)
		t.Assert(dfile.Exists(dstPath), true)
		defer dfile.Remove(dstPath)

		// unzip to temporary dir.
		tempDirPath := dfile.TempDir(dtime.TimestampNanoStr())
		t.Assert(dfile.Mkdir(tempDirPath), nil)
		t.Assert(dcompress.UnZipFile(dstPath, tempDirPath), nil)
		defer dfile.Remove(tempDirPath)

		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "1.txt")),
			dfile.GetContents(srcPath),
		)
	})
	// multiple files
	dtest.C(t, func(t *dtest.T) {
		var (
			srcPath1 = ddebug.TestDataPath("zip", "path1", "1.txt")
			srcPath2 = ddebug.TestDataPath("zip", "path2", "2.txt")
			dstPath  = dfile.TempDir(dtime.TimestampNanoStr(), "zip.zip")
		)
		if p := dfile.Dir(dstPath); !dfile.Exists(p) {
			t.Assert(dfile.Mkdir(p), nil)
		}

		t.Assert(dfile.Exists(dstPath), false)
		err := dcompress.ZipPath(srcPath1+","+srcPath2, dstPath)
		t.Assert(err, nil)
		t.Assert(dfile.Exists(dstPath), true)
		defer dfile.Remove(dstPath)

		// unzip to another temporary dir.
		tempDirPath := dfile.TempDir(dtime.TimestampNanoStr())
		t.Assert(dfile.Mkdir(tempDirPath), nil)
		err = dcompress.UnZipFile(dstPath, tempDirPath)
		t.Assert(err, nil)
		defer dfile.Remove(tempDirPath)

		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "1.txt")),
			dfile.GetContents(srcPath1),
		)
		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "2.txt")),
			dfile.GetContents(srcPath2),
		)
	})
	// one dir and one file.
	dtest.C(t, func(t *dtest.T) {
		var (
			srcPath1 = ddebug.TestDataPath("zip", "path1")
			srcPath2 = ddebug.TestDataPath("zip", "path2", "2.txt")
			dstPath  = dfile.TempDir(dtime.TimestampNanoStr(), "zip.zip")
		)
		if p := dfile.Dir(dstPath); !dfile.Exists(p) {
			t.Assert(dfile.Mkdir(p), nil)
		}

		t.Assert(dfile.Exists(dstPath), false)
		err := dcompress.ZipPath(srcPath1+","+srcPath2, dstPath)
		t.Assert(err, nil)
		t.Assert(dfile.Exists(dstPath), true)
		defer dfile.Remove(dstPath)

		// unzip to another temporary dir.
		tempDirPath := dfile.TempDir(dtime.TimestampNanoStr())
		t.Assert(dfile.Mkdir(tempDirPath), nil)
		err = dcompress.UnZipFile(dstPath, tempDirPath)
		t.Assert(err, nil)
		defer dfile.Remove(tempDirPath)

		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "path1", "1.txt")),
			dfile.GetContents(dfile.Join(srcPath1, "1.txt")),
		)
		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "2.txt")),
			dfile.GetContents(srcPath2),
		)
	})
	// directory.
	dtest.C(t, func(t *dtest.T) {
		srcPath := ddebug.TestDataPath("zip")
		dstPath := ddebug.TestDataPath("zip", "zip.zip")

		pwd := dfile.Pwd()
		err := dfile.Chdir(srcPath)
		defer dfile.Chdir(pwd)
		t.Assert(err, nil)

		t.Assert(dfile.Exists(dstPath), false)
		err = dcompress.ZipPath(srcPath, dstPath)
		t.Assert(err, nil)
		t.Assert(dfile.Exists(dstPath), true)
		defer dfile.Remove(dstPath)

		tempDirPath := dfile.TempDir(dtime.TimestampNanoStr())
		err = dfile.Mkdir(tempDirPath)
		t.Assert(err, nil)

		err = dcompress.UnZipFile(dstPath, tempDirPath)
		t.Assert(err, nil)
		defer dfile.Remove(tempDirPath)

		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "zip", "path1", "1.txt")),
			dfile.GetContents(dfile.Join(srcPath, "path1", "1.txt")),
		)
		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "zip", "path2", "2.txt")),
			dfile.GetContents(dfile.Join(srcPath, "path2", "2.txt")),
		)
	})
	// multiple directory paths joined using char ','.
	dtest.C(t, func(t *dtest.T) {
		var (
			srcPath  = ddebug.TestDataPath("zip")
			srcPath1 = ddebug.TestDataPath("zip", "path1")
			srcPath2 = ddebug.TestDataPath("zip", "path2")
			dstPath  = ddebug.TestDataPath("zip", "zip.zip")
		)
		pwd := dfile.Pwd()
		err := dfile.Chdir(srcPath)
		defer dfile.Chdir(pwd)
		t.Assert(err, nil)

		t.Assert(dfile.Exists(dstPath), false)
		err = dcompress.ZipPath(srcPath1+", "+srcPath2, dstPath)
		t.Assert(err, nil)
		t.Assert(dfile.Exists(dstPath), true)
		defer dfile.Remove(dstPath)

		tempDirPath := dfile.TempDir(dtime.TimestampNanoStr())
		err = dfile.Mkdir(tempDirPath)
		t.Assert(err, nil)

		zipContent := dfile.GetBytes(dstPath)
		t.AssertGT(len(zipContent), 0)
		err = dcompress.UnZipContent(zipContent, tempDirPath)
		t.Assert(err, nil)
		defer dfile.Remove(tempDirPath)

		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "path1", "1.txt")),
			dfile.GetContents(dfile.Join(srcPath, "path1", "1.txt")),
		)
		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "path2", "2.txt")),
			dfile.GetContents(dfile.Join(srcPath, "path2", "2.txt")),
		)
	})
}

func Test_ZipPathWriter(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			srcPath  = ddebug.TestDataPath("zip")
			srcPath1 = ddebug.TestDataPath("zip", "path1")
			srcPath2 = ddebug.TestDataPath("zip", "path2")
		)
		pwd := dfile.Pwd()
		err := dfile.Chdir(srcPath)
		defer dfile.Chdir(pwd)
		t.Assert(err, nil)

		writer := bytes.NewBuffer(nil)
		t.Assert(writer.Len(), 0)
		err = dcompress.ZipPathWriter(srcPath1+", "+srcPath2, writer)
		t.Assert(err, nil)
		t.AssertGT(writer.Len(), 0)

		tempDirPath := dfile.TempDir(dtime.TimestampNanoStr())
		err = dfile.Mkdir(tempDirPath)
		t.Assert(err, nil)

		zipContent := writer.Bytes()
		t.AssertGT(len(zipContent), 0)
		err = dcompress.UnZipContent(zipContent, tempDirPath)
		t.Assert(err, nil)
		defer dfile.Remove(tempDirPath)

		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "path1", "1.txt")),
			dfile.GetContents(dfile.Join(srcPath, "path1", "1.txt")),
		)
		t.Assert(
			dfile.GetContents(dfile.Join(tempDirPath, "path2", "2.txt")),
			dfile.GetContents(dfile.Join(srcPath, "path2", "2.txt")),
		)
	})
}
