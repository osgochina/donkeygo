package dfpool_test

import (
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dfpool"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/test/dtest"
	"os"
	"testing"
	"time"
)

// TestOpen test open file cache
func TestOpen(t *testing.T) {
	testFile := start("TestOpen.txt")

	dtest.C(t, func(t *dtest.T) {
		f, err := dfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertEQ(err, nil)
		f.Close()

		f2, err1 := dfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertEQ(err1, nil)
		t.AssertEQ(f, f2)
		f2.Close()
	})

	stop(testFile)
}

// TestOpenErr test open file error
func TestOpenErr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		testErrFile := "errorPath"
		_, err := dfpool.Open(testErrFile, os.O_RDWR, 0666)
		t.AssertNE(err, nil)

		// delete file error
		testFile := start("TestOpenDeleteErr.txt")
		pool := dfpool.New(testFile, os.O_RDWR, 0666)
		_, err1 := pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertNE(err1, nil)

		// append mode delete file error and create again
		testFile = start("TestOpenCreateErr.txt")
		pool = dfpool.New(testFile, os.O_CREATE, 0666)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)

		// append mode delete file error
		testFile = start("TestOpenAppendErr.txt")
		pool = dfpool.New(testFile, os.O_APPEND, 0666)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertNE(err1, nil)

		// trunc mode delete file error
		testFile = start("TestOpenTruncErr.txt")
		pool = dfpool.New(testFile, os.O_TRUNC, 0666)
		_, err1 = pool.File()
		t.AssertEQ(err1, nil)
		stop(testFile)
		_, err1 = pool.File()
		t.AssertNE(err1, nil)
	})
}

// TestOpenExpire test open file cache expire
func TestOpenExpire(t *testing.T) {
	testFile := start("TestOpenExpire.txt")

	dtest.C(t, func(t *dtest.T) {
		f, err := dfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100*time.Millisecond)
		t.AssertEQ(err, nil)
		f.Close()

		time.Sleep(150 * time.Millisecond)
		f2, err1 := dfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666, 100*time.Millisecond)
		t.AssertEQ(err1, nil)
		//t.AssertNE(f, f2)
		f2.Close()
	})

	stop(testFile)
}

// TestNewPool test gfpool new function
func TestNewPool(t *testing.T) {
	testFile := start("TestNewPool.txt")

	dtest.C(t, func(t *dtest.T) {
		f, err := dfpool.Open(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		t.AssertEQ(err, nil)
		f.Close()

		pool := dfpool.New(testFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
		f2, err1 := pool.File()
		// pool not equal
		t.AssertEQ(err1, nil)
		//t.AssertNE(f, f2)
		f2.Close()

		pool.Close()
	})

	stop(testFile)
}

// test before
func start(name string) string {
	testFile := os.TempDir() + string(os.PathSeparator) + name
	if dfile.Exists(testFile) {
		dfile.Remove(testFile)
	}
	content := "123"
	dfile.PutContents(testFile, content)
	return testFile
}

// test after
func stop(testFile string) {
	if dfile.Exists(testFile) {
		err := dfile.Remove(testFile)
		if err != nil {
			dlog.Error(err)
		}
	}
}
