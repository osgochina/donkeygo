package dlog

import "bytes"

// 写入日志
func (that *Logger) Write(p []byte) (n int, err error) {
	that.Header(false).Print(string(bytes.TrimRight(p, "\r\n")))
	return len(p), nil
}
