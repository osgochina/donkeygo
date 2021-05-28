package dvar

import (
	"github.com/osgochina/donkeygo/internal/empty"
	"reflect"
)

// IsNil 检查只是否为nil
func (that *Var) IsNil() bool {
	return that.Val() == nil
}

// IsEmpty 是否为空
func (that *Var) IsEmpty() bool {
	return empty.IsEmpty(that.Val())
}

// IsInt 是否是int类型
func (that *Var) IsInt() bool {
	switch that.Val().(type) {
	case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64:
		return true
	}
	return false
}

// IsUint 是否是uint类型
func (that *Var) IsUint() bool {
	switch that.Val().(type) {
	case uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		return true
	}
	return false
}

// IsFloat 是否为float类型
func (that *Var) IsFloat() bool {
	switch that.Val().(type) {
	case float32, *float32, float64, *float64:
		return true
	}
	return false
}

// IsSlice 是否为slice类型
func (that *Var) IsSlice() bool {
	var (
		reflectValue = reflect.ValueOf(that.Val())
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// IsMap 是否为map类型
func (that *Var) IsMap() bool {
	var (
		reflectValue = reflect.ValueOf(that.Val())
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectKind {
	case reflect.Map:
		return true
	}
	return false
}

// IsStruct 是否为struct类型
func (that *Var) IsStruct() bool {
	var (
		reflectValue = reflect.ValueOf(that.Val())
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		return true
	}
	return false
}
