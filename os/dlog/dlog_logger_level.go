package dlog

import (
	"errors"
	"fmt"
	"strings"
)

//日志级别
const (
	// LevelAll 所有级别都匹配
	LevelAll = LevelDebug | LevelInfo | LevelNotice | LevelWarning | LevelError | LevelCritical
	// LevelDev 开发状态的日志级别，所有级别都匹配
	LevelDev = LevelAll
	// LevelProd 线上生产环节的日志级别
	LevelProd = LevelWarning | LevelError | LevelCritical
	// LevelDebug debug模式的日志级别
	LevelDebug = 1 << iota // 8
	// LevelInfo 打印详情
	LevelInfo // 16
	// LevelNotice 注意提醒级别
	LevelNotice // 32
	// LevelWarning  警告级别
	LevelWarning // 64
	// LevelError 错误
	LevelError // 128
	// LevelCritical 至关重要的日志
	LevelCritical // 256
	// LevelPanic Panic 错误
	LevelPanic // 512
	// LevelFatal 致命错误
	LevelFatal // 1024
)

// 默认的日志级别前缀
var defaultLevelPrefixes = map[int]string{
	LevelDebug:    "DEBU",
	LevelInfo:     "INFO",
	LevelNotice:   "NOTI",
	LevelWarning:  "WARN",
	LevelError:    "ERRO",
	LevelCritical: "CRIT",
	LevelPanic:    "PANI",
	LevelFatal:    "FATA",
}

var levelStringMap = map[string]int{
	"ALL":      LevelDebug | LevelInfo | LevelNotice | LevelWarning | LevelError | LevelCritical,
	"DEV":      LevelDebug | LevelInfo | LevelNotice | LevelWarning | LevelError | LevelCritical,
	"DEVELOP":  LevelDebug | LevelInfo | LevelNotice | LevelWarning | LevelError | LevelCritical,
	"PROD":     LevelWarning | LevelError | LevelCritical,
	"PRODUCT":  LevelWarning | LevelError | LevelCritical,
	"DEBU":     LevelDebug | LevelInfo | LevelNotice | LevelWarning | LevelError | LevelCritical,
	"DEBUG":    LevelDebug | LevelInfo | LevelNotice | LevelWarning | LevelError | LevelCritical,
	"INFO":     LevelInfo | LevelNotice | LevelWarning | LevelError | LevelCritical,
	"NOTI":     LevelNotice | LevelWarning | LevelError | LevelCritical,
	"NOTICE":   LevelNotice | LevelWarning | LevelError | LevelCritical,
	"WARN":     LevelWarning | LevelError | LevelCritical,
	"WARNING":  LevelWarning | LevelError | LevelCritical,
	"ERRO":     LevelError | LevelCritical,
	"ERROR":    LevelError | LevelCritical,
	"CRIT":     LevelCritical,
	"CRITICAL": LevelCritical,
}

// SetLevel 设置日志级别
func (that *Logger) SetLevel(level int) {
	that.config.Level = level
}

// GetLevel 获取日志级别
func (that *Logger) GetLevel() int {
	return that.config.Level
}

// SetLevelStr 通过可读字符串设置日志级别
func (that *Logger) SetLevelStr(levelStr string) error {
	if level, ok := levelStringMap[strings.ToUpper(levelStr)]; ok {
		that.config.Level = level
	} else {
		return errors.New(fmt.Sprintf(`invalid level string: %s`, levelStr))
	}
	return nil
}

// SetLevelPrefix 设置各种日志级别对应的前缀
func (that *Logger) SetLevelPrefix(level int, prefix string) {
	that.config.LevelPrefixes[level] = prefix
}

// SetLevelPrefixes 设置各种日志级别对应的前缀
func (that *Logger) SetLevelPrefixes(prefixes map[int]string) {
	for k, v := range prefixes {
		that.config.LevelPrefixes[k] = v
	}
}

// GetLevelPrefix 获取指定日志级别的前缀
func (that *Logger) GetLevelPrefix(level int) string {
	return that.config.LevelPrefixes[level]
}

// 获取指定日志级别的打印前缀
func (that *Logger) getLevelPrefixWithBrackets(level int) string {
	if s, ok := that.config.LevelPrefixes[level]; ok {
		return "[" + s + "]"
	}
	return ""
}
