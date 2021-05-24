// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dlog

import (
	"context"
	"io"
)

// Expose 返回默认日志处理对象
func Expose() *Logger {
	return logger
}

// Ctx 设置默认日志处理对象的上下文
func Ctx(ctx context.Context, keys ...interface{}) *Logger {
	return logger.Ctx(ctx, keys...)
}

// To 自定义writer
func To(writer io.Writer) *Logger {
	return logger.To(writer)
}

// Path 设置日志存储目录
func Path(path string) *Logger {
	return logger.Path(path)
}

// Cat 设置日志分类
func Cat(category string) *Logger {
	return logger.Cat(category)
}

// File 设置日志文件名
func File(pattern string) *Logger {
	return logger.File(pattern)
}

// Level 设置日志等级
func Level(level int) *Logger {
	return logger.Level(level)
}

// LevelStr 自然语言设置日志等级
func LevelStr(levelStr string) *Logger {
	return logger.LevelStr(levelStr)
}

// Skip debug栈堆跳过的层
func Skip(skip int) *Logger {
	return logger.Skip(skip)
}

// Stack 打印错误栈堆的参数
func Stack(enabled bool, skip ...int) *Logger {
	return logger.Stack(enabled, skip...)
}

// StackWithFilter 打印错误栈堆的过滤字符串
func StackWithFilter(filter string) *Logger {
	return logger.StackWithFilter(filter)
}

// Stdout  是否打印到标准输出
func Stdout(enabled ...bool) *Logger {
	return logger.Stdout(enabled...)
}

// Header 是否开启日志头
func Header(enabled ...bool) *Logger {
	return logger.Header(enabled...)
}

// Line 日志是打印长文件信息还是短文件信息
func Line(long ...bool) *Logger {
	return logger.Line(long...)
}

// Async 是否异步写日志
func Async(enabled ...bool) *Logger {
	return logger.Async(enabled...)
}
