// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dfile_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_ScanDir(t *testing.T) {
	teatPath := ddebug.TestDataPath()
	dtest.C(t, func(t *dtest.T) {
		files, err := dfile.ScanDir(teatPath, "*", false)
		t.Assert(err, nil)
		t.AssertIN(teatPath+dfile.Separator+"dir1", files)
		t.AssertIN(teatPath+dfile.Separator+"dir2", files)
		t.AssertNE(teatPath+dfile.Separator+"dir1"+dfile.Separator+"file1", files)
	})
	dtest.C(t, func(t *dtest.T) {
		files, err := dfile.ScanDir(teatPath, "*", true)
		t.Assert(err, nil)
		t.AssertIN(teatPath+dfile.Separator+"dir1", files)
		t.AssertIN(teatPath+dfile.Separator+"dir2", files)
		t.AssertIN(teatPath+dfile.Separator+"dir1"+dfile.Separator+"file1", files)
		t.AssertIN(teatPath+dfile.Separator+"dir2"+dfile.Separator+"file2", files)
	})
}

func Test_ScanDirFunc(t *testing.T) {
	teatPath := ddebug.TestDataPath()
	dtest.C(t, func(t *dtest.T) {
		files, err := dfile.ScanDirFunc(teatPath, "*", true, func(path string) string {
			if dfile.Name(path) != "file1" {
				return ""
			}
			return path
		})
		t.Assert(err, nil)
		t.Assert(len(files), 1)
		t.Assert(dfile.Name(files[0]), "file1")
	})
}

func Test_ScanDirFile(t *testing.T) {
	teatPath := ddebug.TestDataPath()
	dtest.C(t, func(t *dtest.T) {
		files, err := dfile.ScanDirFile(teatPath, "*", false)
		t.Assert(err, nil)
		t.Assert(len(files), 0)
	})
	dtest.C(t, func(t *dtest.T) {
		files, err := dfile.ScanDirFile(teatPath, "*", true)
		t.Assert(err, nil)
		t.AssertNI(teatPath+dfile.Separator+"dir1", files)
		t.AssertNI(teatPath+dfile.Separator+"dir2", files)
		t.AssertIN(teatPath+dfile.Separator+"dir1"+dfile.Separator+"file1", files)
		t.AssertIN(teatPath+dfile.Separator+"dir2"+dfile.Separator+"file2", files)
	})
}

func Test_ScanDirFileFunc(t *testing.T) {
	teatPath := ddebug.TestDataPath()
	dtest.C(t, func(t *dtest.T) {
		array := darray.New()
		files, err := dfile.ScanDirFileFunc(teatPath, "*", false, func(path string) string {
			array.Append(1)
			return path
		})
		t.Assert(err, nil)
		t.Assert(len(files), 0)
		t.Assert(array.Len(), 0)
	})
	dtest.C(t, func(t *dtest.T) {
		array := darray.New()
		files, err := dfile.ScanDirFileFunc(teatPath, "*", true, func(path string) string {
			array.Append(1)
			if dfile.Basename(path) == "file1" {
				return path
			}
			return ""
		})
		t.Assert(err, nil)
		t.Assert(len(files), 1)
		t.Assert(array.Len(), 3)
	})
}
