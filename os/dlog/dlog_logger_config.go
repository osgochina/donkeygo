package dlog

import "io"

type Config struct {
	Writer io.Writer `json:"-"`      // 定制 io.Writer.
	Path   string    `json:"path"`   // 日志目录地址
	File   string    `json:"file"`   // 当前日志需要保存的文件
	Level  int       `json:"level"`  // 日志的输出等级
	Prefix string    `json:"prefix"` // 每条日志的前缀
}

func DefaultConfig() Config {
	c := Config{
		File: defaultFileFormat,
	}

	return c
}
