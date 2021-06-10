package dcfg

import (
	"errors"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gtime"
	"github.com/osgochina/donkeygo/container/dvar"
	"time"
)

// Set 修改内存中的配置内容
func (that *Config) Set(pattern string, value interface{}) error {
	if j := that.getJson(); j != nil {
		return j.Set(pattern, value)
	}
	return nil
}

// Get 获取配置文件的选项
func (that *Config) Get(pattern string, def ...interface{}) interface{} {
	if j := that.getJson(); j != nil {
		return j.Get(pattern, def...)
	}
	return nil
}

// GetVar 返回`pattern`对应的Var对象
func (that *Config) GetVar(pattern string, def ...interface{}) *dvar.Var {
	if j := that.getJson(); j != nil {
		return dvar.New(j.GetVar(pattern, def...).Val())
	}
	return dvar.New(nil)
}

// Contains 判断配置项是否存在
func (that *Config) Contains(pattern string) bool {
	if j := that.getJson(); j != nil {
		return j.Contains(pattern)
	}
	return false
}

// GetMap 获取map格式的配置项
func (that *Config) GetMap(pattern string, def ...interface{}) map[string]interface{} {
	if j := that.getJson(); j != nil {
		return j.GetMap(pattern, def...)
	}
	return nil
}

// GetMapStrStr 获取 map[string]string格式的配置项
func (that *Config) GetMapStrStr(pattern string, def ...interface{}) map[string]string {
	if j := that.getJson(); j != nil {
		return j.GetMapStrStr(pattern, def...)
	}
	return nil
}

// GetArray 获取数组格式的配置项
func (that *Config) GetArray(pattern string, def ...interface{}) []interface{} {
	if j := that.getJson(); j != nil {
		return j.GetArray(pattern, def...)
	}
	return nil
}

// GetBytes 获取配置项的字节数组类型
func (that *Config) GetBytes(pattern string, def ...interface{}) []byte {
	if j := that.getJson(); j != nil {
		return j.GetBytes(pattern, def...)
	}
	return nil
}

// GetString 获取字符串类型的配置项
func (that *Config) GetString(pattern string, def ...interface{}) string {
	if j := that.getJson(); j != nil {
		return j.GetString(pattern, def...)
	}
	return ""
}

// GetStrings 获取字符串数组配置项
func (that *Config) GetStrings(pattern string, def ...interface{}) []string {
	if j := that.getJson(); j != nil {
		return j.GetStrings(pattern, def...)
	}
	return nil
}

// GetInterfaces 获取数组格式的配置项
func (that *Config) GetInterfaces(pattern string, def ...interface{}) []interface{} {
	if j := that.getJson(); j != nil {
		return j.GetInterfaces(pattern, def...)
	}
	return nil
}

func (that *Config) GetBool(pattern string, def ...interface{}) bool {
	if j := that.getJson(); j != nil {
		return j.GetBool(pattern, def...)
	}
	return false
}

// GetFloat32 retrieves the value by specified `pattern` and converts it to float32.
func (that *Config) GetFloat32(pattern string, def ...interface{}) float32 {
	if j := that.getJson(); j != nil {
		return j.GetFloat32(pattern, def...)
	}
	return 0
}

// GetFloat64 retrieves the value by specified `pattern` and converts it to float64.
func (that *Config) GetFloat64(pattern string, def ...interface{}) float64 {
	if j := that.getJson(); j != nil {
		return j.GetFloat64(pattern, def...)
	}
	return 0
}

// GetFloats retrieves the value by specified `pattern` and converts it to []float64.
func (that *Config) GetFloats(pattern string, def ...interface{}) []float64 {
	if j := that.getJson(); j != nil {
		return j.GetFloats(pattern, def...)
	}
	return nil
}

// GetInt retrieves the value by specified `pattern` and converts it to int.
func (that *Config) GetInt(pattern string, def ...interface{}) int {
	if j := that.getJson(); j != nil {
		return j.GetInt(pattern, def...)
	}
	return 0
}

// GetInt8 retrieves the value by specified `pattern` and converts it to int8.
func (that *Config) GetInt8(pattern string, def ...interface{}) int8 {
	if j := that.getJson(); j != nil {
		return j.GetInt8(pattern, def...)
	}
	return 0
}

// GetInt16 retrieves the value by specified `pattern` and converts it to int16.
func (that *Config) GetInt16(pattern string, def ...interface{}) int16 {
	if j := that.getJson(); j != nil {
		return j.GetInt16(pattern, def...)
	}
	return 0
}

// GetInt32 retrieves the value by specified `pattern` and converts it to int32.
func (that *Config) GetInt32(pattern string, def ...interface{}) int32 {
	if j := that.getJson(); j != nil {
		return j.GetInt32(pattern, def...)
	}
	return 0
}

// GetInt64 retrieves the value by specified `pattern` and converts it to int64.
func (that *Config) GetInt64(pattern string, def ...interface{}) int64 {
	if j := that.getJson(); j != nil {
		return j.GetInt64(pattern, def...)
	}
	return 0
}

// GetInts retrieves the value by specified `pattern` and converts it to []int.
func (that *Config) GetInts(pattern string, def ...interface{}) []int {
	if j := that.getJson(); j != nil {
		return j.GetInts(pattern, def...)
	}
	return nil
}

// GetUint retrieves the value by specified `pattern` and converts it to uint.
func (that *Config) GetUint(pattern string, def ...interface{}) uint {
	if j := that.getJson(); j != nil {
		return j.GetUint(pattern, def...)
	}
	return 0
}

// GetUint8 retrieves the value by specified `pattern` and converts it to uint8.
func (that *Config) GetUint8(pattern string, def ...interface{}) uint8 {
	if j := that.getJson(); j != nil {
		return j.GetUint8(pattern, def...)
	}
	return 0
}

// GetUint16 retrieves the value by specified `pattern` and converts it to uint16.
func (that *Config) GetUint16(pattern string, def ...interface{}) uint16 {
	if j := that.getJson(); j != nil {
		return j.GetUint16(pattern, def...)
	}
	return 0
}

// GetUint32 retrieves the value by specified `pattern` and converts it to uint32.
func (that *Config) GetUint32(pattern string, def ...interface{}) uint32 {
	if j := that.getJson(); j != nil {
		return j.GetUint32(pattern, def...)
	}
	return 0
}

// GetUint64 retrieves the value by specified `pattern` and converts it to uint64.
func (that *Config) GetUint64(pattern string, def ...interface{}) uint64 {
	if j := that.getJson(); j != nil {
		return j.GetUint64(pattern, def...)
	}
	return 0
}

// GetTime retrieves the value by specified `pattern` and converts it to time.Time.
func (that *Config) GetTime(pattern string, format ...string) time.Time {
	if j := that.getJson(); j != nil {
		return j.GetTime(pattern, format...)
	}
	return time.Time{}
}

// GetDuration retrieves the value by specified `pattern` and converts it to time.Duration.
func (that *Config) GetDuration(pattern string, def ...interface{}) time.Duration {
	if j := that.getJson(); j != nil {
		return j.GetDuration(pattern, def...)
	}
	return 0
}

// GetGTime retrieves the value by specified `pattern` and converts it to *gtime.Time.
func (that *Config) GetGTime(pattern string, format ...string) *gtime.Time {
	if j := that.getJson(); j != nil {
		return j.GetGTime(pattern, format...)
	}
	return nil
}

// GetJson gets the value by specified `pattern`,
// and converts it to a un-concurrent-safe Json object.
func (that *Config) GetJson(pattern string, def ...interface{}) *gjson.Json {
	if j := that.getJson(); j != nil {
		return j.GetJson(pattern, def...)
	}
	return nil
}

// GetJsons gets the value by specified `pattern`,
// and converts it to a slice of un-concurrent-safe Json object.
func (that *Config) GetJsons(pattern string, def ...interface{}) []*gjson.Json {
	if j := that.getJson(); j != nil {
		return j.GetJsons(pattern, def...)
	}
	return nil
}

// GetJsonMap gets the value by specified `pattern`,
// and converts it to a map of un-concurrent-safe Json object.
func (that *Config) GetJsonMap(pattern string, def ...interface{}) map[string]*gjson.Json {
	if j := that.getJson(); j != nil {
		return j.GetJsonMap(pattern, def...)
	}
	return nil
}

// GetStruct retrieves the value by specified `pattern` and converts it to specified object
// `pointer`. The `pointer` should be the pointer to an object.
func (that *Config) GetStruct(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.GetStruct(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// GetStructDeep does GetStruct recursively.
// Deprecated, use GetStruct instead.
func (that *Config) GetStructDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.GetStructDeep(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// GetStructs 将任何片转换为给定的结构片。
func (that *Config) GetStructs(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.GetStructs(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// GetStructsDeep 递归地将任何片转换为给定的结构片。
// Deprecated, use GetStructs instead.
func (that *Config) GetStructsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.GetStructsDeep(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// GetMapToMap 根据指定的“模式”检索值并将其转换为指定的映射变量。
// See gconv.MapToMap.
func (that *Config) GetMapToMap(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.GetMapToMap(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// GetMapToMaps 根据指定的“模式”检索值，并将其转换为指定的映射片变量。
// See gconv.MapToMaps.
func (that *Config) GetMapToMaps(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.GetMapToMaps(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// GetMapToMapsDeep 根据指定的“模式”检索值，并将其递归转换为指定的映射片变量。
// See gconv.MapToMapsDeep.
func (that *Config) GetMapToMapsDeep(pattern string, pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.GetMapToMapsDeep(pattern, pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// ToMap 将当前Json对象转换为 map[string]interface{}.
// 如果失败返回nil。
func (that *Config) ToMap() map[string]interface{} {
	if j := that.getJson(); j != nil {
		return j.ToMap()
	}
	return nil
}

// ToArray converts current Json object to []interface{}.
// It returns nil if fails.
func (that *Config) ToArray() []interface{} {
	if j := that.getJson(); j != nil {
		return j.ToArray()
	}
	return nil
}

// ToStruct 将当前Json对象转换为指定对象。
// 指针应该是*struct类型的指针。
func (that *Config) ToStruct(pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.ToStruct(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// ToStructDeep 将当前Json对象递归转换为指定对象。
//指针应该是*struct类型的指针。
func (that *Config) ToStructDeep(pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.ToStructDeep(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// ToStructs 将当前Json对象转换为指定的对象切片。
//指针类型应该是[]struct/*struct。
func (that *Config) ToStructs(pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.ToStructs(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// ToStructsDeep 将当前Json对象递归转换为指定的对象切片。
// 指针类型应该是[]struct/*struct。
func (that *Config) ToStructsDeep(pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.ToStructsDeep(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// ToMapToMap 将当前Json对象转换为指定的映射变量。
// 指针的形参应该是*map类型。
func (that *Config) ToMapToMap(pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.ToMapToMap(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// ToMapToMaps 将当前Json对象转换为指定的映射变量切片。
// 指针的参数应该是[]map/*map类型。
func (that *Config) ToMapToMaps(pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.ToMapToMaps(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// ToMapToMapsDeep
// 将当前Json对象递归转换为指定的映射变量切片。
// 指针的参数应该是[]map/*map类型。
func (that *Config) ToMapToMapsDeep(pointer interface{}, mapping ...map[string]string) error {
	if j := that.getJson(); j != nil {
		return j.ToMapToMapsDeep(pointer, mapping...)
	}
	return errors.New("configuration not found")
}

// Clear 删除所有已解析的配置文件内容缓存，
//这将强制从文件重新加载配置内容。
func (that *Config) Clear() {
	that.jsonMap.Clear()
}

// Dump 打印当前Json对象，更具手动可读性。
func (that *Config) Dump() {
	if j := that.getJson(); j != nil {
		j.Dump()
	}
}
