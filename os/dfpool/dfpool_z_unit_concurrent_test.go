package dfpool_test

import (
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dfpool"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"os"
	"testing"
)

func Test_ConcurrentOS(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		defer dfile.Remove(path)
		f1, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f1.Close()

		f2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f2.Close()

		for i := 0; i < 100; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}
		for i := 0; i < 100; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}

		for i := 0; i < 1000; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}
		for i := 0; i < 1000; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}
		t.Assert(dstr.Count(dfile.GetContents(path), "@1234567890#"), 2200)
	})

	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		defer dfile.Remove(path)
		f1, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f1.Close()

		f2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f2.Close()

		for i := 0; i < 1000; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}
		for i := 0; i < 1000; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}
		t.Assert(dstr.Count(dfile.GetContents(path), "@1234567890#"), 2000)
	})
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		defer dfile.Remove(path)
		f1, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f1.Close()

		f2, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f2.Close()

		s1 := ""
		for i := 0; i < 1000; i++ {
			s1 += "@1234567890#"
		}
		_, err = f2.Write([]byte(s1))
		t.Assert(err, nil)

		s2 := ""
		for i := 0; i < 1000; i++ {
			s2 += "@1234567890#"
		}
		_, err = f2.Write([]byte(s2))
		t.Assert(err, nil)

		t.Assert(dstr.Count(dfile.GetContents(path), "@1234567890#"), 2000)
	})
}

func Test_ConcurrentGFPool(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		defer dfile.Remove(path)
		f1, err := dfpool.Open(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f1.Close()

		f2, err := dfpool.Open(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.Assert(err, nil)
		defer f2.Close()

		for i := 0; i < 1000; i++ {
			_, err = f1.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}
		for i := 0; i < 1000; i++ {
			_, err = f2.Write([]byte("@1234567890#"))
			t.Assert(err, nil)
		}
		t.Assert(dstr.Count(dfile.GetContents(path), "@1234567890#"), 2000)
	})
}
