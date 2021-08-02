package dlog

import (
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/util/dconv"
	"github.com/osgochina/donkeygo/util/dutil"
	"io"
	"strings"
	"time"
)

// Config 日志的配置信息
type Config struct {
	Writer               io.Writer      `json:"-"`                    // 定制 io.Writer.
	Flags                int            `json:"flags"`                // 用户记录日志的额外标识
	Path                 string         `json:"path"`                 // 日志目录地址
	File                 string         `json:"file"`                 // 当前日志需要保存的文件
	Level                int            `json:"level"`                // 日志的输出等级
	Prefix               string         `json:"prefix"`               // 每条日志的前缀
	StSkip               int            `json:"stSkip"`               // 跳过多少层栈堆不用打印
	StStatus             int            `json:"stStatus"`             // 是否打开栈堆信息 1：打开，0关闭，默认是打开
	StFilter             string         `json:"stFilter"`             // 栈堆字符串过滤
	CtxKeys              []interface{}  `json:"ctxKeys"`              // 用于日志记录的上下文键，用于从上下文检索值。
	StdoutPrint          bool           `json:"stdout"`               // 开启日志的标准输出，默认是true.
	HeaderPrint          bool           `json:"header"`               // 打印日志标题，默认是true
	LevelPrefixes        map[int]string `json:"levelPrefixes"`        // 日志级别的前缀映射
	RotateCheckInterval  time.Duration  `json:"rotateCheckInterval"`  // 每隔多久检查一下日志文件，看看日志文件是否需要压缩备份。默认是1小时
	RotateSize           int64          `json:"rotateSize"`           // 如果日志文件 size > RotateSize 则旋转日志.
	RotateExpire         time.Duration  `json:"rotateExpire"`         // 如果日志文件的最后更新日志大于RotateExpire则旋转日志。
	RotateBackupLimit    int            `json:"rotateBackupLimit"`    // 最大旋转备份的文件数量，默认是0，意味着没有备份.
	RotateBackupExpire   time.Duration  `json:"rotateBackupExpire"`   // 旋转文件的最大过期时间，默认是0，意味着不过期。
	RotateBackupCompress int            `json:"rotateBackupCompress"` // 使用gzip算法压缩旋转文件的级别。默认值是0，表示没有压缩。
}

// DefaultConfig 生成默认配置信息
func DefaultConfig() Config {
	c := Config{
		File:                defaultFileFormat,
		Flags:               FTimeStd,
		Level:               LevelAll,
		StStatus:            1,
		HeaderPrint:         true,
		StdoutPrint:         true,
		LevelPrefixes:       make(map[int]string, len(defaultLevelPrefixes)),
		RotateCheckInterval: time.Hour,
	}
	for k, v := range defaultLevelPrefixes {
		c.LevelPrefixes[k] = v
	}
	if !defaultDebug {
		c.Level = c.Level & ^LevelDebug
	}
	return c
}

// SetConfig 设置日志的配置信息
func (that *Logger) SetConfig(config Config) error {
	that.config = config
	if that.config.Path != "" {
		if err := that.SetPath(config.Path); err != nil {
			intlog.Error(err)
			return err
		}
	}
	intlog.Printf("SetConfig: %+v", that.config)
	return nil
}

// SetConfigWithMap 通过map传入配置信息
func (that *Logger) SetConfigWithMap(m map[string]interface{}) error {
	//传入的配置必须有数据
	if m == nil || len(m) == 0 {
		return errors.New("configuration cannot be empty")
	}

	m = dutil.MapCopy(m)
	// 查找map中的日志级别
	levelKey, levelValue := dutil.MapPossibleItemByKey(m, "Level")
	if levelValue != nil {
		if level, ok := levelStringMap[strings.ToUpper(dconv.String(levelValue))]; ok {
			m[levelKey] = level
		} else {
			return errors.New(fmt.Sprintf(`invalid level string: %v`, levelValue))
		}
	}
	// 文件大小限制，超过该大小则重新建一个文件写
	rotateSizeKey, rotateSizeValue := dutil.MapPossibleItemByKey(m, "RotateSize")
	if rotateSizeValue != nil {
		m[rotateSizeKey] = dfile.StrToSize(dconv.String(rotateSizeValue))
		if m[rotateSizeKey] == -1 {
			return errors.New(fmt.Sprintf(`invalid rotate size: %v`, rotateSizeValue))
		}
	}
	err := dconv.Struct(m, &that.config)
	if err != nil {
		return err
	}
	return that.SetConfig(that.config)
}

// SetDebug 开启日志级别为debug模式
func (that *Logger) SetDebug(debug bool) {
	if debug {
		that.config.Level = that.config.Level | LevelDebug
	} else {
		that.config.Level = that.config.Level & ^LevelDebug
	}
}

// SetAsync 设置写入日志到文件为异步模式
func (that *Logger) SetAsync(enabled bool) {
	if enabled {
		that.config.Flags = that.config.Flags | FAsync
	} else {
		that.config.Flags = that.config.Flags & ^FAsync
	}
}

// SetFlags 设置写日志信息的额外标识
func (that *Logger) SetFlags(flags int) {
	that.config.Flags = flags
}

// GetFlags 获取额外标识
func (that *Logger) GetFlags() int {
	return that.config.Flags
}

// SetStack 启动/禁用 日志输出的时候是否打印栈堆信息
func (that *Logger) SetStack(enabled bool) {
	if enabled {
		that.config.StStatus = 1
	} else {
		that.config.StStatus = 0
	}
}

// SetStackSkip 设置打印栈堆信息的偏移量，从错误触发点开始算第几层
func (that *Logger) SetStackSkip(skip int) {
	that.config.StSkip = skip
}

// SetStackFilter 设置栈堆字符串过滤器
func (that *Logger) SetStackFilter(filter string) {
	that.config.StFilter = filter
}

// SetCtxKeys 设置日志的上下文键。键用于检索值 从上下文和打印到日志内容。
func (that *Logger) SetCtxKeys(keys ...interface{}) {
	that.config.CtxKeys = keys
}

// GetCtxKeys 获取上下文记录的key值
func (that *Logger) GetCtxKeys() []interface{} {
	return that.config.CtxKeys
}

//SetWriter 设置自定义日志用于日志记录。
//对象应该实现io.writer接口。
//开发者可以使用自定义日志将日志输出重定向到另一个服务，
//例如:kafka, mysql, mongodb等。
func (that *Logger) SetWriter(writer io.Writer) {
	that.config.Writer = writer
}

// GetWriter 获取当前的日志写接口
func (that *Logger) GetWriter() io.Writer {
	return that.config.Writer
}

// SetPath 设置日志的路径
func (that *Logger) SetPath(path string) error {
	if path == "" {
		return errors.New("logging path is empty")
	}
	if !dfile.Exists(path) {
		if err := dfile.Mkdir(path); err != nil {
			return derror.Wrapf(err, `Mkdir "%s" failed in PWD "%s"`, path, dfile.Pwd())
		}
	}
	that.config.Path = strings.TrimRight(path, dfile.Separator)
	return nil
}

// GetPath 获取当前日志的写入路径
func (that *Logger) GetPath() string {
	return that.config.Path
}

// SetFile 日志的文件名
func (that *Logger) SetFile(pattern string) {
	that.config.File = pattern
}

// SetStdoutPrint 开启关闭日志的标准输出
func (that *Logger) SetStdoutPrint(enabled bool) {
	that.config.StdoutPrint = enabled
}

// SetHeaderPrint 是否打印日志头
func (that *Logger) SetHeaderPrint(enabled bool) {
	that.config.HeaderPrint = enabled
}

// SetPrefix 设置日志的前缀
func (that *Logger) SetPrefix(prefix string) {
	that.config.Prefix = prefix
}
