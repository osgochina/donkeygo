package dlog

import (
	"github.com/osgochina/donkeygo/container/dmap"
)

const (
	// DefaultName Default group name for instance usage.
	DefaultName = "default"
)

var (
	// Instances map.
	instances = dmap.NewStrAnyMap(true)
)

// Instance 获取单例对象
func Instance(name ...string) *Logger {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		return New()
	}).(*Logger)
}
