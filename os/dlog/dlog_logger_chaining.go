package dlog

import (
	"context"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dfile"
	"io"
)

// Ctx 日志的上下文
func (that *Logger) Ctx(ctx context.Context, keys ...interface{}) *Logger {
	if ctx == nil {
		return that
	}
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	logger.ctx = ctx
	if len(keys) > 0 {
		logger.SetCtxKeys(keys...)
	}
	return logger
}

// To 发送到自定义的writer
func (that *Logger) To(writer io.Writer) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	logger.SetWriter(writer)
	return logger
}

// Path 设置日志的路径
func (that *Logger) Path(path string) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	if path != "" {
		if err := logger.SetPath(path); err != nil {
			// panic(err)
			intlog.Error(context.TODO(), err)
		}
	}
	return logger
}

// Cat 设置日志的指定输出层
func (that *Logger) Cat(category string) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	if logger.config.Path != "" {
		if err := logger.SetPath(dfile.Join(logger.config.Path, category)); err != nil {
			// panic(err)
			intlog.Error(context.TODO(), err)
		}
	}
	return logger
}

// File 设置日志的输出文件
func (that *Logger) File(file string) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	logger.SetFile(file)
	return logger
}

// Level 日志等级
func (that *Logger) Level(level int) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	logger.SetLevel(level)
	return logger
}

// LevelStr 通过可读的日志级别设置日志级别
func (that *Logger) LevelStr(levelStr string) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	if err := logger.SetLevelStr(levelStr); err != nil {
		// panic(err)
		intlog.Error(context.TODO(), err)
	}
	return logger
}

// Skip 设置打印栈堆跳过的层级
func (that *Logger) Skip(skip int) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	logger.SetStackSkip(skip)
	return logger
}

// Stack 设置栈堆参数
func (that *Logger) Stack(enabled bool, skip ...int) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	logger.SetStack(enabled)
	if len(skip) > 0 {
		logger.SetStackSkip(skip[0])
	}
	return logger
}

// StackWithFilter 设置栈堆过滤字符串
func (that *Logger) StackWithFilter(filter string) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	logger.SetStack(true)
	logger.SetStackFilter(filter)
	return logger
}

// Stdout 设置日志打印到标准输出，默认true
func (that *Logger) Stdout(enabled ...bool) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	// stdout printing is enabled if <enabled> is not passed.
	if len(enabled) > 0 && !enabled[0] {
		logger.config.StdoutPrint = false
	} else {
		logger.config.StdoutPrint = true
	}
	return logger
}

// Header 打印日志头
func (that *Logger) Header(enabled ...bool) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	// header is enabled if <enabled> is not passed.
	if len(enabled) > 0 && !enabled[0] {
		logger.SetHeaderPrint(false)
	} else {
		logger.SetHeaderPrint(true)
	}
	return logger
}

// Line 设置日志打印时候是显示详细的文件信息和行号还是显示简短的文件路径和行号
func (that *Logger) Line(long ...bool) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	if len(long) > 0 && long[0] {
		logger.config.Flags |= FFileLong
	} else {
		logger.config.Flags |= FFileShort
	}
	return logger
}

// Async 设置是异步写文件还是同步写文件
func (that *Logger) Async(enabled ...bool) *Logger {
	logger := (*Logger)(nil)
	if that.parent == nil {
		logger = that.Clone()
	} else {
		logger = that
	}
	// async feature is enabled if <enabled> is not passed.
	if len(enabled) > 0 && !enabled[0] {
		logger.SetAsync(false)
	} else {
		logger.SetAsync(true)
	}
	return logger
}
