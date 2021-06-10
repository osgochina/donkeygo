package dcfg_test

import (
	"github.com/osgochina/donkeygo/os/dcfg"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"os"
	"testing"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
)

func init() {
	os.Setenv("DK_DCFG_ERRORPRINT", "false")
}

func Test_Basic1(t *testing.T) {
	config := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	dtest.C(t, func(t *dtest.T) {
		path := dcfg.DefaultConfigFile
		err := dfile.PutContents(path, config)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		c := dcfg.New()
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
		t.AssertEQ(c.GetFloat32("v4"), float32(1.23))
		t.AssertEQ(c.GetFloat64("v4"), float64(1.23))
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

func Test_Basic2(t *testing.T) {
	config := `log-path = "logs"`
	dtest.C(t, func(t *dtest.T) {
		path := dcfg.DefaultConfigFile
		err := dfile.PutContents(path, config)
		t.Assert(err, nil)
		defer func() {
			_ = dfile.Remove(path)
		}()

		c := dcfg.New()
		t.Assert(c.Get("log-path"), "logs")
	})
}

func Test_Content(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	dcfg.SetContent(content)
	defer dcfg.ClearContent()

	dtest.C(t, func(t *dtest.T) {
		c := dcfg.New()
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
		t.AssertEQ(c.GetFloat32("v4"), float32(1.23))
		t.AssertEQ(c.GetFloat64("v4"), float64(1.23))
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
	})
}

func Test_SetFileName(t *testing.T) {
	config := `
{
	"array": [
		1,
		2,
		3
	],
	"redis": {
		"cache": "127.0.0.1:6379,1",
		"disk": "127.0.0.1:6379,0"
	},
	"v1": 1,
	"v2": "true",
	"v3": "off",
	"v4": "1.234"
}
`
	dtest.C(t, func(t *dtest.T) {
		path := "config.json"
		err := dfile.PutContents(path, config)
		t.Assert(err, nil)
		defer func() {
			_ = dfile.Remove(path)
		}()

		c := dcfg.New()
		c.SetFileName(path)
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

func TestCfg_New(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		os.Setenv("DK_DCFG_PATH", "config")
		c := dcfg.New("config.yml")
		t.Assert(c.Get("name"), nil)
		t.Assert(c.GetFileName(), "config.yml")

		configPath := dfile.Pwd() + dfile.Separator + "config"
		_ = dfile.Mkdir(configPath)
		defer dfile.Remove(configPath)

		c = dcfg.New("config.yml")
		t.Assert(c.Get("name"), nil)

		_ = os.Unsetenv("DK_DCFG_PATH")
		c = dcfg.New("config.yml")
		t.Assert(c.Get("name"), nil)
	})
}

func TestCfg_SetPath(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		c := dcfg.New("config.yml")
		err := c.SetPath("tmp")
		t.AssertNE(err, nil)
		err = c.SetPath("gcfg.go")
		t.AssertNE(err, nil)
		t.Assert(c.Get("name"), nil)
	})
}

func TestCfg_SetViolenceCheck(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		c := dcfg.New("config.yml")
		c.SetViolenceCheck(true)
		t.Assert(c.Get("name"), nil)
	})
}

func TestCfg_AddPath(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		c := dcfg.New("config.yml")
		err := c.AddPath("tmp")
		t.AssertNE(err, nil)
		err = c.AddPath("gcfg.go")
		t.AssertNE(err, nil)
		t.Assert(c.Get("name"), nil)
	})
}

func TestCfg_FilePath(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		c := dcfg.New("config.yml")
		path, _ := c.GetFilePath("tmp")
		t.Assert(path, "")
		path, _ = c.GetFilePath("tmp")
		t.Assert(path, "")
	})
}

func TestCfg_et(t *testing.T) {
	config := `log-path = "logs"`
	dtest.C(t, func(t *dtest.T) {
		path := dcfg.DefaultConfigFile
		err := dfile.PutContents(path, config)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		c := dcfg.New()
		t.Assert(c.Get("log-path"), "logs")

		err = c.Set("log-path", "custom-logs")
		t.Assert(err, nil)
		t.Assert(c.Get("log-path"), "custom-logs")
	})
}

func TestCfg_Get(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var err error
		configPath := dfile.TempDir(dtime.TimestampNanoStr())
		err = dfile.Mkdir(configPath)
		t.Assert(err, nil)
		defer dfile.Remove(configPath)

		defer dfile.Chdir(dfile.Pwd())
		err = dfile.Chdir(configPath)
		t.Assert(err, nil)

		err = dfile.PutContents(
			dfile.Join(configPath, "config.yml"),
			"wrong config",
		)
		t.Assert(err, nil)
		c := dcfg.New("config.yml")
		t.Assert(c.Get("name"), nil)
		t.Assert(c.GetVar("name").Val(), nil)
		t.Assert(c.Contains("name"), false)
		t.Assert(c.GetMap("name"), nil)
		t.Assert(c.GetArray("name"), nil)
		t.Assert(c.GetString("name"), "")
		t.Assert(c.GetStrings("name"), nil)
		t.Assert(c.GetInterfaces("name"), nil)
		t.Assert(c.GetBool("name"), false)
		t.Assert(c.GetFloat32("name"), 0)
		t.Assert(c.GetFloat64("name"), 0)
		t.Assert(c.GetFloats("name"), nil)
		t.Assert(c.GetInt("name"), 0)
		t.Assert(c.GetInt8("name"), 0)
		t.Assert(c.GetInt16("name"), 0)
		t.Assert(c.GetInt32("name"), 0)
		t.Assert(c.GetInt64("name"), 0)
		t.Assert(c.GetInts("name"), nil)
		t.Assert(c.GetUint("name"), 0)
		t.Assert(c.GetUint8("name"), 0)
		t.Assert(c.GetUint16("name"), 0)
		t.Assert(c.GetUint32("name"), 0)
		t.Assert(c.GetUint64("name"), 0)
		t.Assert(c.GetTime("name").Format("2006-01-02"), "0001-01-01")
		t.Assert(c.GetGTime("name"), nil)
		t.Assert(c.GetDuration("name").String(), "0s")
		name := struct {
			Name string
		}{}
		t.Assert(c.GetStruct("name", &name) == nil, false)

		c.Clear()

		arr, _ := gjson.Encode(
			g.Map{
				"name":   "gf",
				"time":   "2019-06-12",
				"person": g.Map{"name": "gf"},
				"floats": g.Slice{1, 2, 3},
			},
		)
		err = dfile.PutBytes(
			dfile.Join(configPath, "config.yml"),
			arr,
		)
		t.Assert(err, nil)
		t.Assert(c.GetTime("time").Format("2006-01-02"), "2019-06-12")
		t.Assert(c.GetGTime("time").Format("Y-m-d"), "2019-06-12")
		t.Assert(c.GetDuration("time").String(), "0s")

		err = c.GetStruct("person", &name)
		t.Assert(err, nil)
		t.Assert(name.Name, "gf")
		t.Assert(c.GetFloats("floats") == nil, false)
	})
}

func TestCfg_Config(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		dcfg.SetContent("gf", "config.yml")
		t.Assert(dcfg.GetContent("config.yml"), "gf")
		dcfg.SetContent("gf1", "config.yml")
		t.Assert(dcfg.GetContent("config.yml"), "gf1")
		dcfg.RemoveContent("config.yml")
		dcfg.ClearContent()
		t.Assert(dcfg.GetContent("name"), "")
	})
}

func TestCfg_With_UTF8_BOM(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cfg := dcfg.Instance("test-cfg-with-utf8-bom")
		t.Assert(cfg.SetPath("testdata"), nil)
		cfg.SetFileName("cfg-with-utf8-bom.toml")
		t.Assert(cfg.GetInt("test.testInt"), 1)
		t.Assert(cfg.GetString("test.testStr"), "test")
	})
}
