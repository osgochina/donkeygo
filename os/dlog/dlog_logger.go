package dlog

import (
	"context"
	"donkeygo/container/dtype"
	"donkeygo/text/dregex"
	"github.com/gogf/gf/os/gfile"
	"io"
	"strings"
	"time"
)

const (
	defaultFileFormat = `{Y-m-d}.log`
)

type Logger struct {
	ctx    context.Context
	init   *dtype.Bool
	parent *Logger
	config Config
}

func New() *Logger {
	return &Logger{
		init:   dtype.NewBool(),
		config: DefaultConfig(),
	}
}

// 获取将要写入的日志文件名
func (that *Logger) getFilePath(now time.Time) string {
	file, _ := dregex.ReplaceStringFunc(`{.+?}`, that.config.Path, func(s string) string {
		return now.Format(strings.Trim(s, "{}"))
		//return gtime.New(now).Format(strings.Trim(s, "{}"))
	})
	file = gfile.Join(that.config.Path, file)
	return file
}

func (that *Logger) print(std io.Writer, lead string, values ...interface{}) {

}
