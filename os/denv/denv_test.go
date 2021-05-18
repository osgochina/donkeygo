package denv_test

import (
	"github.com/osgochina/donkeygo/os/denv"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"os"
	"testing"
	"time"
)

func Test_DEnv_All(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(os.Environ(), denv.All())
	})
}

func Test_DEnv_Map(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := dconv.String(time.Nanosecond)
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.Assert(err, nil)
		t.Assert(denv.Map()[key], "TEST")
	})
}

func Test_DEnv_Get(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := dconv.String(time.Nanosecond)
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.Assert(err, nil)
		t.AssertEQ(denv.Get(key), "TEST")
	})
}

func Test_DEnv_Contains(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := dconv.String(time.Nanosecond)
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.Assert(err, nil)
		t.AssertEQ(denv.Contains(key), true)
		t.AssertEQ(denv.Contains("none"), false)
	})
}

func Test_DEnv_Set(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := dconv.String(time.Nanosecond)
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.Assert(err, nil)
		t.Assert(os.Getenv(key), "TEST")
	})
}

func Test_DEnv_SetMap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := denv.SetMap(map[string]string{
			"K1": "TEST1",
			"K2": "TEST2",
		})
		t.Assert(err, nil)
		t.AssertEQ(os.Getenv("K1"), "TEST1")
		t.AssertEQ(os.Getenv("K2"), "TEST2")
	})
}

func Test_DEnv_Build(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := denv.Build(map[string]string{
			"k1": "v1",
			"k2": "v2",
		})
		t.AssertIN("k1=v1", s)
		t.AssertIN("k2=v2", s)
	})
}

func Test_DEnv_Remove(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := dconv.String(time.Nanosecond)
		key := "TEST_ENV_" + value
		err := os.Setenv(key, "TEST")
		t.Assert(err, nil)
		err = denv.Remove(key)
		t.Assert(err, nil)
		t.AssertEQ(os.Getenv(key), "")
	})
}
