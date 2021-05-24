package dlog

import (
	"github.com/osgochina/donkeygo/os/dcmd"
	"github.com/osgochina/donkeygo/os/dgpool"
)

var (

	//默认的日志
	logger = New()

	//默认是否开启debug
	defaultDebug = true

	//用于异步输出日志的goroutine池
	asyncPool = dgpool.New(1)
)

func init() {
	defaultDebug = dcmd.GetOptWithEnv("dk.dlog.debug", true).Bool()
	SetDebug(defaultDebug)
}

// DefaultLogger 获取默认的日志处理对象
func DefaultLogger() *Logger {
	return logger
}

// SetDefaultLogger 设置默认的日志处理对象
func SetDefaultLogger(l *Logger) {
	logger = l
}
