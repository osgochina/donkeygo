package dcfg

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/os/denv"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_Instance_Basic(t *testing.T) {
	config := `
array = [1.0, 2.0, 3.0]
v1 = 1.0
v2 = "true"
v3 = "off"
v4 = "1.234"

[redis]
  cache = "127.0.0.1:6379,1"
  disk = "127.0.0.1:6379,0"

`
	dtest.C(t, func(t *dtest.T) {
		path := DefaultConfigFile
		err := dfile.PutContents(path, config)
		t.Assert(err, nil)
		defer func() {
			t.Assert(dfile.Remove(path), nil)
		}()

		c := Instance()
		t.Assert(c.Get("v1"), 1)
		t.AssertEQ(c.GetInt("v1"), 1)
		t.AssertEQ(c.GetInt8("v1"), int8(1))
		t.AssertEQ(c.GetInt16("v1"), int16(1))
		t.AssertEQ(c.GetInt32("v1"), int32(1))
		t.AssertEQ(c.GetInt64("v1"), int64(1))
		t.AssertEQ(c.GetUint("v1"), uint(1))
		t.AssertEQ(c.GetUint8("v1"), uint8(1))
		t.AssertEQ(c.GetUint16("v1"), uint16(1))
		t.AssertEQ(c.GetUint32("v1"), uint32(1))
		t.AssertEQ(c.GetUint64("v1"), uint64(1))

		t.AssertEQ(c.GetVar("v1").String(), "1")
		t.AssertEQ(c.GetVar("v1").Bool(), true)
		t.AssertEQ(c.GetVar("v2").String(), "true")
		t.AssertEQ(c.GetVar("v2").Bool(), true)

		t.AssertEQ(c.GetString("v1"), "1")
		t.AssertEQ(c.GetFloat32("v4"), float32(1.234))
		t.AssertEQ(c.GetFloat64("v4"), float64(1.234))
		t.AssertEQ(c.GetString("v2"), "true")
		t.AssertEQ(c.GetBool("v2"), true)
		t.AssertEQ(c.GetBool("v3"), false)

		t.AssertEQ(c.Contains("v1"), true)
		t.AssertEQ(c.Contains("v2"), true)
		t.AssertEQ(c.Contains("v3"), true)
		t.AssertEQ(c.Contains("v4"), true)
		t.AssertEQ(c.Contains("v5"), false)

		t.AssertEQ(c.GetInts("array"), []int{1, 2, 3})
		t.AssertEQ(c.GetStrings("array"), []string{"1", "2", "3"})
		t.AssertEQ(c.GetArray("array"), []interface{}{1, 2, 3})
		t.AssertEQ(c.GetInterfaces("array"), []interface{}{1, 2, 3})
		t.AssertEQ(c.GetMap("redis"), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		filepath, _ := c.GetFilePath()
		t.AssertEQ(filepath, dfile.Pwd()+dfile.Separator+path)
	})
}

func Test_Instance_AutoLocateConfigFile(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(Instance("gf") != nil, true)
	})
	// Automatically locate the configuration file with supported file extensions.
	dtest.C(t, func(t *dtest.T) {
		pwd := dfile.Pwd()
		t.Assert(dfile.Chdir(ddebug.TestDataPath()), nil)
		defer dfile.Chdir(pwd)
		t.Assert(Instance("c1") != nil, true)
		t.Assert(Instance("c1").Get("my-config"), "1")
		t.Assert(Instance("folder1/c1").Get("my-config"), "2")
	})
	// Automatically locate the configuration file with supported file extensions.
	dtest.C(t, func(t *dtest.T) {
		pwd := dfile.Pwd()
		t.Assert(dfile.Chdir(ddebug.TestDataPath("folder1")), nil)
		defer dfile.Chdir(pwd)
		t.Assert(Instance("c2").Get("my-config"), 2)
	})
	// Default configuration file.
	dtest.C(t, func(t *dtest.T) {
		instances.Clear()
		pwd := dfile.Pwd()
		t.Assert(dfile.Chdir(ddebug.TestDataPath("default")), nil)
		defer dfile.Chdir(pwd)
		t.Assert(Instance().Get("my-config"), 1)

		instances.Clear()
		t.Assert(denv.Set("DK_DCFG_FILE", "config.json"), nil)
		defer denv.Set("DK_DCFG_FILE", "")
		t.Assert(Instance().Get("my-config"), 2)
	})
}

func Test_Instance_EnvPath(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		denv.Set("DK_DCFG_PATH", ddebug.TestDataPath("envpath"))
		defer denv.Set("DK_DCFG_PATH", "")
		t.Assert(Instance("c3") != nil, true)
		t.Assert(Instance("c3").Get("my-config"), "3")
		t.Assert(Instance("c4").Get("my-config"), "4")
		instances = dmap.NewStrAnyMap(true)
	})
}

func Test_Instance_EnvFile(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		file := ddebug.TestDataPath("envfile")
		denv.Set("DK_DCFG_PATH", file)
		defer denv.Set("DK_DCFG_PATH", "")
		denv.Set("DK_DCFG_FILE", "c6.json")
		defer denv.Set("DK_DCFG_FILE", "")
		t.Assert(Instance().Get("my-config"), "6")
		instances = dmap.NewStrAnyMap(true)
	})
}
