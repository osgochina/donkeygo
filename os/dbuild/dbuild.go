package dbuild

import (
	"context"
	"github.com/osgochina/donkeygo"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/encoding/dbase64"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/util/dconv"
	"runtime"
)

var (
	builtInVarStr = ""                       // Raw variable base64 string.
	builtInVarMap = map[string]interface{}{} // Binary custom variable map decoded.
)

func init() {
	if builtInVarStr != "" {
		err := json.UnmarshalUseNumber(dbase64.MustDecodeString(builtInVarStr), &builtInVarMap)
		if err != nil {
			intlog.Error(context.TODO(), err)
		}
		builtInVarMap["dkVersion"] = donkeygo.VERSION
		builtInVarMap["goVersion"] = runtime.Version()
		intlog.Printf(context.TODO(), "build variables: %+v", builtInVarMap)
	} else {
		intlog.Print(context.TODO(), "no build variables")
	}
}

// Info returns the basic built information of the binary as map.
// Note that it should be used with gf-cli tool "gf build",
// which injects necessary information into the binary.
func Info() map[string]string {
	return map[string]string{
		"dk":   GetString("dkVersion"),
		"go":   GetString("goVersion"),
		"git":  GetString("builtGit"),
		"time": GetString("builtTime"),
	}
}

// Get retrieves and returns the build-in binary variable with given name.
func Get(name string, def ...interface{}) interface{} {
	if v, ok := builtInVarMap[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// GetVar retrieves and returns the build-in binary variable of given name as gvar.Var.
func GetVar(name string, def ...interface{}) *dvar.Var {
	return dvar.New(Get(name, def...))
}

// GetString retrieves and returns the build-in binary variable of given name as string.
func GetString(name string, def ...interface{}) string {
	return dconv.String(Get(name, def...))
}

// Map returns the custom build-in variable map.
func Map() map[string]interface{} {
	return builtInVarMap
}
