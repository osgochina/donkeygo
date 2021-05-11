package derror

import "runtime"

type stack []uintptr

const (
	// maxStackDepth 最深追踪的栈堆层级
	maxStackDepth = 32
)

// 追踪栈堆，skip是指定跳过的层级
func callers(skip ...int) stack {
	var (
		pcs [maxStackDepth]uintptr
		n   = 3
	)
	if len(skip) > 0 {
		n += skip[0]
	}
	return pcs[:runtime.Callers(n, pcs[:])]
}
