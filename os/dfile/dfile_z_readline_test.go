// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dfile_test

import (
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_NotFound(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		teatFile := dfile.Dir(ddebug.CallerFilePath()) + dfile.Separator + "testdata/readline/error.log"
		callback := func(line string) error {
			return nil
		}
		err := dfile.ReadLines(teatFile, callback)
		t.AssertNE(err, nil)
	})
}

func Test_ReadLines(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			expectList = []string{"a", "b", "c", "d", "e"}
			getList    = make([]string, 0)
			callback   = func(line string) error {
				getList = append(getList, line)
				return nil
			}
			teatFile = dfile.Dir(ddebug.CallerFilePath()) + dfile.Separator + "testdata/readline/file.log"
		)
		err := dfile.ReadLines(teatFile, callback)
		t.AssertEQ(getList, expectList)
		t.AssertEQ(err, nil)
	})
}

func Test_ReadLines_Error(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			callback = func(line string) error {
				return derror.New("custom error")
			}
			teatFile = dfile.Dir(ddebug.CallerFilePath()) + dfile.Separator + "testdata/readline/file.log"
		)
		err := dfile.ReadLines(teatFile, callback)
		t.AssertEQ(err.Error(), "custom error")
	})
}

func Test_ReadLinesBytes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			expectList = [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e")}
			getList    = make([][]byte, 0)
			callback   = func(line []byte) error {
				getList = append(getList, line)
				return nil
			}
			teatFile = dfile.Dir(ddebug.CallerFilePath()) + dfile.Separator + "testdata/readline/file.log"
		)
		err := dfile.ReadLinesBytes(teatFile, callback)
		t.AssertEQ(getList, expectList)
		t.AssertEQ(err, nil)
	})
}

func Test_ReadLinesBytes_Error(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			callback = func(line []byte) error {
				return derror.New("custom error")
			}
			teatFile = dfile.Dir(ddebug.CallerFilePath()) + dfile.Separator + "testdata/readline/file.log"
		)
		err := dfile.ReadLinesBytes(teatFile, callback)
		t.AssertEQ(err.Error(), "custom error")
	})
}
