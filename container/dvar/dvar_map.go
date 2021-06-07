package dvar

import (
	"github.com/osgochina/donkeygo/util/dconv"
)

// Map 返回一个map
func (that *Var) Map(tags ...string) map[string]interface{} {
	return dconv.Map(that.Val(), tags...)
}

// MapStrAny 返回map[string]interface{}格式
func (that *Var) MapStrAny() map[string]interface{} {
	return that.Map()
}

// MapStrStr 返回字符串map
func (that *Var) MapStrStr(tags ...string) map[string]string {
	return dconv.MapStrStr(that.Val(), tags...)
}

// MapStrVar 返回一个字符串为key，Var为值的map
func (that *Var) MapStrVar(tags ...string) map[string]*Var {
	m := that.Map(tags...)
	if len(m) > 0 {
		vMap := make(map[string]*Var, len(m))
		for k, v := range m {
			vMap[k] = New(v)
		}
		return vMap
	}
	return nil
}

// MapDeep 返回var转换成map[string]interface{}
func (that *Var) MapDeep(tags ...string) map[string]interface{} {
	return dconv.MapDeep(that.Val(), tags...)
}

// MapStrStrDeep 返回var转换成map[string]string
func (that *Var) MapStrStrDeep(tags ...string) map[string]string {
	return dconv.MapStrStrDeep(that.Val(), tags...)
}

// MapStrVarDeep 返回var转换成 map[string]*Var
func (that *Var) MapStrVarDeep(tags ...string) map[string]*Var {
	m := that.MapDeep(tags...)
	if len(m) > 0 {
		vMap := make(map[string]*Var, len(m))
		for k, v := range m {
			vMap[k] = New(v)
		}
		return vMap
	}
	return nil
}

// Maps 以切片的形式返回多个map，
func (that *Var) Maps(tags ...string) []map[string]interface{} {
	return dconv.Maps(that.Val(), tags...)
}

// MapToMap converts any map type variable <params> to another map type variable <pointer>.
// See gconv.MapToMap.
func (that *Var) MapToMap(pointer interface{}, mapping ...map[string]string) (err error) {
	return dconv.MapToMap(that.Val(), pointer, mapping...)
}

// MapToMaps converts any map type variable <params> to another map type variable <pointer>.
// See gconv.MapToMaps.
func (that *Var) MapToMaps(pointer interface{}, mapping ...map[string]string) (err error) {
	return dconv.MapToMaps(that.Val(), pointer, mapping...)
}

// MapToMapsDeep converts any map type variable <params> to another map type variable
// <pointer> recursively.
// See gconv.MapToMapsDeep.
func (that *Var) MapToMapsDeep(pointer interface{}, mapping ...map[string]string) (err error) {
	return dconv.MapToMapsDeep(that.Val(), pointer, mapping...)
}
