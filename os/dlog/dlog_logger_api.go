package dlog

import (
	"fmt"
	"os"
)

// Print 标准打印
func (that *Logger) Print(v ...interface{}) {
	that.printStd("", v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func (that *Logger) Printf(format string, v ...interface{}) {
	that.printStd("", that.format(format, v...))
}

// Println is alias of Print.
// See Print.
func (that *Logger) Println(v ...interface{}) {
	that.Print(v...)
}

// Fatal 打印致命错误，并结束当前进程
func (that *Logger) Fatal(v ...interface{}) {
	that.printErr(that.getLevelPrefixWithBrackets(LevelFatal), v...)
	os.Exit(1)
}

// Fatalf 格式化打印致命错误，并结束当前进程
func (that *Logger) Fatalf(format string, v ...interface{}) {
	that.printErr(that.getLevelPrefixWithBrackets(LevelFatal), that.format(format, v...))
	os.Exit(1)
}

// Panic 打印日志并且 panic
func (that *Logger) Panic(v ...interface{}) {
	that.printErr(that.getLevelPrefixWithBrackets(LevelPanic), v...)
	panic(fmt.Sprint(v...))
}

// Panicf 打印日志并且 panic
func (that *Logger) Panicf(format string, v ...interface{}) {
	that.printErr(that.getLevelPrefixWithBrackets(LevelPanic), that.format(format, v...))
	panic(that.format(format, v...))
}

// Info 打印info日志
func (that *Logger) Info(v ...interface{}) {
	if that.checkLevel(LevelInfo) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelInfo), v...)
	}
}

// Infof 打印info日志
func (that *Logger) Infof(format string, v ...interface{}) {
	if that.checkLevel(LevelInfo) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelInfo), that.format(format, v...))
	}
}

// Debug 打印debug日志
func (that *Logger) Debug(v ...interface{}) {
	if that.checkLevel(LevelDebug) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelDebug), v...)
	}
}

// Debugf 打印debug日志
func (that *Logger) Debugf(format string, v ...interface{}) {
	if that.checkLevel(LevelDebug) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelDebug), that.format(format, v...))
	}
}

// Notice 打印提醒日志
func (that *Logger) Notice(v ...interface{}) {
	if that.checkLevel(LevelNotice) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelNotice), v...)
	}
}

// Noticef 打印提醒日志
func (that *Logger) Noticef(format string, v ...interface{}) {
	if that.checkLevel(LevelNotice) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelNotice), that.format(format, v...))
	}
}

// Warning 打印警告日志
func (that *Logger) Warning(v ...interface{}) {
	if that.checkLevel(LevelWarning) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelWarning), v...)
	}
}

// Warningf 打印警告日志
func (that *Logger) Warningf(format string, v ...interface{}) {
	if that.checkLevel(LevelWarning) {
		that.printStd(that.getLevelPrefixWithBrackets(LevelWarning), that.format(format, v...))
	}
}

// Error 打印错误日志
func (that *Logger) Error(v ...interface{}) {
	if that.checkLevel(LevelError) {
		that.printErr(that.getLevelPrefixWithBrackets(LevelError), v...)
	}
}

// Errorf 打印错误日志
func (that *Logger) Errorf(format string, v ...interface{}) {
	if that.checkLevel(LevelError) {
		that.printErr(that.getLevelPrefixWithBrackets(LevelError), that.format(format, v...))
	}
}

// Critical 打印重要的日志
func (that *Logger) Critical(v ...interface{}) {
	if that.checkLevel(LevelCritical) {
		that.printErr(that.getLevelPrefixWithBrackets(LevelCritical), v...)
	}
}

// Criticalf 打印重要的日志
func (that *Logger) Criticalf(format string, v ...interface{}) {
	if that.checkLevel(LevelCritical) {
		that.printErr(that.getLevelPrefixWithBrackets(LevelCritical), that.format(format, v...))
	}
}

//检查日志等级是否可以输出
func (that *Logger) checkLevel(level int) bool {
	return that.config.Level&level > 0
}
