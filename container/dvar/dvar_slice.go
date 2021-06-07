// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/osgochina/donkeygo/util/dconv"
)

// Ints 转换并返回int数组
func (that *Var) Ints() []int {
	return dconv.Ints(that.Val())
}

// Int64s  转换并返回Int64s数组
func (that *Var) Int64s() []int64 {
	return dconv.Int64s(that.Val())
}

// Uints 转换并返回uint数组
func (that *Var) Uints() []uint {
	return dconv.Uints(that.Val())
}

// Uint64s converts and returns <v> as []uint64.
func (that *Var) Uint64s() []uint64 {
	return dconv.Uint64s(that.Val())
}

// Floats is alias of Float64s.
func (that *Var) Floats() []float64 {
	return dconv.Floats(that.Val())
}

// Float32s converts and returns <v> as []float32.
func (that *Var) Float32s() []float32 {
	return gconv.Float32s(that.Val())
}

// Float64s converts and returns <v> as []float64.
func (that *Var) Float64s() []float64 {
	return dconv.Float64s(that.Val())
}

// Strings converts and returns <v> as []string.
func (that *Var) Strings() []string {
	return dconv.Strings(that.Val())
}

// Interfaces converts and returns <v> as []interfaces{}.
func (that *Var) Interfaces() []interface{} {
	return dconv.Interfaces(that.Val())
}

// Slice is alias of Interfaces.
func (that *Var) Slice() []interface{} {
	return that.Interfaces()
}

// Array is alias of Interfaces.
func (that *Var) Array() []interface{} {
	return that.Interfaces()
}

// Vars converts and returns <v> as []Var.
func (that *Var) Vars() []*Var {
	array := dconv.Interfaces(that.Val())
	if len(array) == 0 {
		return nil
	}
	vars := make([]*Var, len(array))
	for k, v := range array {
		vars[k] = New(v)
	}
	return vars
}
