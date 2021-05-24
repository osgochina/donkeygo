package intlog

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/osgochina/donkeygo/internal/utils"
	"path/filepath"
	"time"
)

const (
	stackFilterKey = "/internal/intlog"
)

var isDKDebug = false

func init() {
	isDKDebug = utils.IsDebugEnabled()
}

// SetEnabled 开启关闭内部日志，非并发安全
func SetEnabled(enabled bool) {
	if isDKDebug != enabled {
		isDKDebug = enabled
	}
}

// Print 打印日志
func Print(v ...interface{}) {
	if !isDKDebug {
		return
	}
	fmt.Println(append([]interface{}{now(), "[INTE]", file()}, v...)...)
}

// Printf 打印日志
func Printf(format string, v ...interface{}) {
	if !isDKDebug {
		return
	}
	fmt.Printf(now()+" [INTE] "+file()+" "+format+"\n", v...)
}

func Error(v ...interface{}) {
	if !isDKDebug {
		return
	}
	array := append([]interface{}{now(), "[INTE]", file()}, v...)
	array = append(array, "\n"+gdebug.StackWithFilter(stackFilterKey))
	fmt.Println(array...)
}

func Errorf(format string, v ...interface{}) {
	if !isDKDebug {
		return
	}
	fmt.Printf(
		now()+" [INTE] "+file()+" "+format+"\n%s\n",
		append(v, gdebug.StackWithFilter(stackFilterKey))...,
	)
}

// 当前时间
func now() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// 返回调用者的文件名和行号
func file() string {
	_, p, l := gdebug.CallerWithFilter(stackFilterKey)
	return fmt.Sprintf(`%s:%d`, filepath.Base(p), l)
}
