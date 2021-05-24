package ddebug

import (
	"runtime"
	"strings"
)

const (
	maxCallerDepth = 1000
	stackFilterKey = "/debug/ddebug/ddebug"
)

var (
	goRootForFilter  = runtime.GOROOT() // goRootForFilter is used for stack filtering purpose.
	binaryVersion    = ""               // The version of current running binary(uint64 hex).
	binaryVersionMd5 = ""               // The version of current running binary(MD5).
	selfPath         = ""               // Current running binary absolute path.
)

// CallerWithFilter  returns the function name and the absolute file path along with
// its line number of the caller.
//
// The parameter <filter> is used to filter the path of the caller.
func CallerWithFilter(filter string, skip ...int) (function string, path string, line int) {
	number := 0
	if len(skip) > 0 {
		number = skip[0]
	}
	ok := true
	pc, file, line, start := callerFromIndex([]string{filter})
	if start != -1 {
		for i := start + number; i < maxCallerDepth; i++ {
			if i != start {
				pc, file, line, ok = runtime.Caller(i)
			}
			if ok {
				if filter != "" && strings.Contains(file, filter) {
					continue
				}
				if strings.Contains(file, stackFilterKey) {
					continue
				}
				function := ""
				if fn := runtime.FuncForPC(pc); fn == nil {
					function = "unknown"
				} else {
					function = fn.Name()
				}
				return function, file, line
			} else {
				break
			}
		}
	}
	return "", "", -1
}

// callerFromIndex returns the caller position and according information exclusive of the
// debug package.
//
// VERY NOTE THAT, the returned index value should be <index - 1> as the caller's start point.
func callerFromIndex(filters []string) (pc uintptr, file string, line int, index int) {
	var filtered, ok bool
	for index = 0; index < maxCallerDepth; index++ {
		if pc, file, line, ok = runtime.Caller(index); ok {
			filtered = false
			for _, filter := range filters {
				if filter != "" && strings.Contains(file, filter) {
					filtered = true
					break
				}
			}
			if filtered {
				continue
			}
			if strings.Contains(file, stackFilterKey) {
				continue
			}
			if index > 0 {
				index--
			}
			return
		}
	}
	return 0, "", -1, -1
}
