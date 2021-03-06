// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar

import (
	"github.com/osgochina/donkeygo/util/dconv"
)

// Struct maps value of <v> to <pointer>.
// The parameter <pointer> should be a pointer to a struct instance.
// The parameter <mapping> is used to specify the key-to-attribute mapping rules.
func (that *Var) Struct(pointer interface{}, mapping ...map[string]string) error {
	return dconv.Struct(that.Val(), pointer, mapping...)
}

// StructDeep maps value of <v> to <pointer> recursively.
// The parameter <pointer> should be a pointer to a struct instance.
// The parameter <mapping> is used to specify the key-to-attribute mapping rules.
// Deprecated, use Struct instead.
func (that *Var) StructDeep(pointer interface{}, mapping ...map[string]string) error {
	return dconv.StructDeep(that.Val(), pointer, mapping...)
}

// Structs converts and returns <v> as given struct slice.
func (that *Var) Structs(pointer interface{}, mapping ...map[string]string) error {
	return dconv.Structs(that.Val(), pointer, mapping...)
}

// StructsDeep converts and returns <v> as given struct slice recursively.
// Deprecated, use Struct instead.
func (that *Var) StructsDeep(pointer interface{}, mapping ...map[string]string) error {
	return dconv.StructsDeep(that.Val(), pointer, mapping...)
}

// Scan automatically calls Struct or Structs function according to the type of parameter
// <pointer> to implement the converting.
// It calls function Struct if <pointer> is type of *struct/**struct to do the converting.
// It calls function Structs if <pointer> is type of *[]struct/*[]*struct to do the converting.
func (that *Var) Scan(pointer interface{}, mapping ...map[string]string) error {
	return dconv.Scan(that.Val(), pointer, mapping...)
}

// ScanDeep automatically calls StructDeep or StructsDeep function according to the type of
// parameter <pointer> to implement the converting.
// It calls function StructDeep if <pointer> is type of *struct/**struct to do the converting.
// It calls function StructsDeep if <pointer> is type of *[]struct/*[]*struct to do the converting.
func (that *Var) ScanDeep(pointer interface{}, mapping ...map[string]string) error {
	return dconv.ScanDeep(that.Val(), pointer, mapping...)
}
