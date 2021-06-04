package dfile

import (
	"os"
	"time"
)

// MTime 返回文件或目录路径的最后修改时间
func MTime(path string) time.Time {
	s, e := os.Stat(path)
	if e != nil {
		return time.Time{}
	}
	return s.ModTime()
}

// MTimestamp 返回文件或目录路径的最后修改时间戳
func MTimestamp(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.Unix()
}

// MTimestampMilli 返回文件或目录路径的最后修改毫秒时间戳
func MTimestampMilli(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.UnixNano() / 1000000
}
