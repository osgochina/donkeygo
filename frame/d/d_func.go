// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package d

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gproc"
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/util/dutil"
)

// NewVar returns a dvar.Var.
func NewVar(i interface{}, safe ...bool) *Var {
	return dvar.New(i, safe...)
}

// Wait is an alias of ghttp.Wait, which blocks until all the web servers shutdown.
// It's commonly used in multiple servers situation.
func Wait() {
	ghttp.Wait()
}

// Listen is an alias of gproc.Listen, which handles the signals received and automatically
// calls registered signal handler functions.
// It blocks until shutdown signals received and all registered shutdown handlers done.
func Listen() {
	gproc.Listen()
}

// Dump dumps a variable to stdout with more manually readable.
func Dump(i ...interface{}) {
	dutil.Dump(i...)
}

// Export exports a variable to string with more manually readable.
func Export(i ...interface{}) string {
	return dutil.Export(i...)
}

// Throw throws a exception, which can be caught by TryCatch function.
// It always be used in TryCatch function.
func Throw(exception interface{}) {
	dutil.Throw(exception)
}

// Try implements try... logistics using internal panic...recover.
// It returns error if any exception occurs, or else it returns nil.
func Try(try func()) (err error) {
	return dutil.Try(try)
}

// TryCatch implements try...catch... logistics using internal panic...recover.
// It automatically calls function <catch> if any exception occurs ans passes the exception as an error.
func TryCatch(try func(), catch ...func(exception error)) {
	dutil.TryCatch(try, catch...)
}

// IsNil checks whether given <value> is nil.
// Parameter <traceSource> is used for tracing to the source variable if given <value> is type
// of a pinter that also points to a pointer. It returns nil if the source is nil when <traceSource>
// is true.
// Note that it might use reflect feature which affects performance a little bit.
func IsNil(value interface{}, traceSource ...bool) bool {
	return empty.IsNil(value, traceSource...)
}

// IsEmpty checks whether given <value> empty.
// It returns true if <value> is in: 0, nil, false, "", len(slice/map/chan) == 0.
// Or else it returns true.
func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}
