// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dcompress_test

import (
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/encoding/dcompress"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_Gzip_UnGzip(t *testing.T) {
	src := "Hello World!!"

	gzip := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xff,
		0xf2, 0x48, 0xcd, 0xc9, 0xc9,
		0x57, 0x08, 0xcf, 0x2f, 0xca,
		0x49, 0x51, 0x54, 0x04, 0x04,
		0x00, 0x00, 0xff, 0xff, 0x9d,
		0x24, 0xa8, 0xd1, 0x0d, 0x00,
		0x00, 0x00,
	}
	dtest.C(t, func(t *dtest.T) {
		arr := []byte(src)
		data, _ := dcompress.Gzip(arr)
		t.Assert(data, gzip)

		data, _ = dcompress.UnGzip(gzip)
		t.Assert(data, arr)

		data, _ = dcompress.UnGzip(gzip[1:])
		t.Assert(data, nil)
	})
}

func Test_Gzip_UnGzip_File(t *testing.T) {
	srcPath := ddebug.TestDataPath("gzip", "file.txt")
	dstPath1 := dfile.TempDir(dtime.TimestampNanoStr(), "gzip.zip")
	dstPath2 := dfile.TempDir(dtime.TimestampNanoStr(), "file.txt")

	// Compress.
	dtest.C(t, func(t *dtest.T) {
		err := dcompress.GzipFile(srcPath, dstPath1, 9)
		t.Assert(err, nil)
		defer dfile.Remove(dstPath1)
		t.Assert(dfile.Exists(dstPath1), true)

		// Decompress.
		err = dcompress.UnGzipFile(dstPath1, dstPath2)
		t.Assert(err, nil)
		defer dfile.Remove(dstPath2)
		t.Assert(dfile.Exists(dstPath2), true)

		t.Assert(dfile.GetContents(srcPath), dfile.GetContents(dstPath2))
	})
}
