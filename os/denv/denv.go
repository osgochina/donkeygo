package denv

import (
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/os/dcmd"
	"os"
	"strings"
)

// All 返回所有的环境变量
func All() []string {
	return os.Environ()
}

// Map 把环境变量转换成map格式
func Map() map[string]string {
	m := make(map[string]string)
	i := 0
	for _, e := range os.Environ() {
		i = strings.IndexByte(e, '=')
		m[e[0:i]] = e[i+1:]
	}
	return m
}

// Get 获取指定key的环境变量值
func Get(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}

// GetVar 获取var类型的变量
func GetVar(key string, def ...interface{}) *dvar.Var {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return dvar.New(def[0])
	}
	return dvar.New(v)
}

// Set 设置环境变量
func Set(key string, value string) error {
	return os.Setenv(key, value)
}

// SetMap 批量设置环境变量
func SetMap(m map[string]string) error {
	for k, v := range m {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Contains 判断环境变量是否存在
func Contains(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Remove 移除环境变量，支持多个key
func Remove(key ...string) error {
	var err error
	for _, k := range key {
		err = os.Unsetenv(k)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetWithCmd 如果环境变量中没有，则去cmd中获取
func GetWithCmd(key string, def ...interface{}) *dvar.Var {
	value := interface{}(nil)
	if len(def) > 0 {
		value = def[0]
	}
	envKey := strings.ToUpper(strings.Replace(key, ".", "_", -1))
	if v := os.Getenv(envKey); v != "" {
		value = v
	} else {
		cmdKey := strings.ToLower(strings.Replace(key, "_", ".", -1))
		if v := dcmd.GetOpt(cmdKey); v != "" {
			value = v
		}
	}
	return dvar.New(value)
}

// Build 创建环境变量
func Build(m map[string]string) []string {
	array := make([]string, len(m))
	index := 0
	for k, v := range m {
		array[index] = k + "=" + v
		index++
	}
	return array
}
