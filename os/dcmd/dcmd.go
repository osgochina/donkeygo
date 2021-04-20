package dcmd

import (
	"donkeygo/container/dvar"
	"donkeygo/internal/command"
	"os"
	"strings"
)

var (
	defaultCommandFuncMap = make(map[string]func())
)

// Init 初始化参数
func Init(args ...string) {
	command.Init(args...)
}

// GetOpt 获取option中的值
func GetOpt(name string, def ...string) string {
	Init()
	return command.GetOpt(name, def...)
}

// GetOptVar 获取var类型的options
func GetOptVar(name string, def ...string) *dvar.Var {
	Init()
	return dvar.New(GetOpt(name, def...))
}

// GetOptAll 获取全部options
func GetOptAll() map[string]string {
	Init()
	return command.GetOptAll()
}

// ContainsOpt 判断指定name的option是否存在
func ContainsOpt(name string, def ...string) bool {
	Init()
	return command.ContainsOpt(name)
}

// GetArg 获取参数
func GetArg(index int, def ...string) string {
	Init()
	return command.GetArg(index, def...)
}

// GetArgVar 获取var类型的参数
func GetArgVar(index int, def ...string) *dvar.Var {
	Init()
	return dvar.New(GetArg(index, def...))
}

// GetArgAll 获取全部参数
func GetArgAll() []string {
	Init()
	return command.GetArgAll()
}

// GetWithEnv 获取参数，如果命令行不存在则从环境变量中获取
func GetWithEnv(key string, def ...interface{}) *dvar.Var {
	return GetOptWithEnv(key, def...)
}

// GetOptWithEnv returns the command line argument of the specified <key>.
// If the argument does not exist, then it returns the environment variable with specified <key>.
// It returns the default value <def> if none of them exists.
//
// Fetching Rules:
// 1. Command line arguments are in lowercase format, eg: gf.<package name>.<variable name>;
// 2. Environment arguments are in uppercase format, eg: GF_<package name>_<variable name>；
func GetOptWithEnv(key string, def ...interface{}) *dvar.Var {
	cmdKey := strings.ToLower(strings.Replace(key, "_", ".", -1))
	if ContainsOpt(cmdKey) {
		return dvar.New(GetOpt(cmdKey))
	} else {
		envKey := strings.ToUpper(strings.Replace(key, ".", "_", -1))
		if r, ok := os.LookupEnv(envKey); ok {
			return dvar.New(r)
		} else {
			if len(def) > 0 {
				return dvar.New(def[0])
			}
		}
	}
	return dvar.New(nil)
}

// BuildOptions 通过map构造出options
func BuildOptions(m map[string]string, prefix ...string) string {
	options := ""
	leadStr := "-"
	if len(prefix) > 0 {
		leadStr = prefix[0]
	}
	for k, v := range m {
		if len(options) > 0 {
			options += " "
		}
		options += leadStr + k
		if v != "" {
			options += "=" + v
		}
	}
	return options
}
