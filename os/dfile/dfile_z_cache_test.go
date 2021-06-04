package dfile_test

import (
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/test/dtest"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_GetContentsWithCache(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var f *os.File
		var err error
		fileName := "test"
		strTest := "123"

		if !dfile.Exists(fileName) {
			f, err = ioutil.TempFile("", fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		defer f.Close()
		defer os.Remove(f.Name())

		if dfile.Exists(f.Name()) {
			f, err = dfile.OpenFile(f.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("file open fail", err)
			}

			err = dfile.PutContents(f.Name(), strTest)
			if err != nil {
				t.Error("write error", err)
			}

			cache := dfile.GetContentsWithCache(f.Name(), 1)
			t.Assert(cache, strTest)
		}
	})

	dtest.C(t, func(t *dtest.T) {

		var f *os.File
		var err error
		fileName := "test2"
		strTest := "123"

		if !dfile.Exists(fileName) {
			f, err = ioutil.TempFile("", fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		defer f.Close()
		defer os.Remove(f.Name())

		if dfile.Exists(f.Name()) {
			cache := dfile.GetContentsWithCache(f.Name())

			f, err = dfile.OpenFile(f.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("file open fail", err)
			}

			err = dfile.PutContents(f.Name(), strTest)
			if err != nil {
				t.Error("write error", err)
			}

			t.Assert(cache, "")

			time.Sleep(100 * time.Millisecond)
		}
	})
}
