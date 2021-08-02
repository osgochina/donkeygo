package dlog

import (
	"bytes"
	"context"
	"fmt"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dfpool"
	"github.com/osgochina/donkeygo/os/dmlock"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/text/dregex"
	"github.com/osgochina/donkeygo/util/dconv"
	"io"
	"os"
	"strings"
	"time"
)

const (
	defaultFileFormat = `{Y-m-d}.log`   //默认的日志文件名格式
	pathFilterKey     = "/os/dlog/dlog" //
	defaultFileFlags  = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	defaultFilePerm   = os.FileMode(0666)
	defaultFileExpire = time.Minute
)

const (
	FAsync     = 1 << iota              // 异步打印日志内容
	FFileLong                           // 打印完整的文件名和行号: /a/b/c/d.go:23.
	FFileShort                          // 打印最终的文件名和行号: d.go:23. 覆盖 FFileLong.
	FTimeDate                           // 打印当地时区的日期信息: 2009-01-23.
	FTimeTime                           // 打印当地时区的时间信息: 01:23:23.
	FTimeMilli                          // 以当地时间的毫秒为单位打印信息: 01:23:23.675.
	FCallerFn                           // 打印调用这的函数名及包名: main.main
	FTimeStd   = FTimeDate | FTimeMilli //标准时间格式
)

// Logger 日志对象
type Logger struct {
	ctx    context.Context //日志对象上下文
	init   *dtype.Bool     // 是否初始化
	parent *Logger         //
	config Config          // 日志配置
}

// New 创建日志对象
func New() *Logger {
	return &Logger{
		init:   dtype.NewBool(),
		config: DefaultConfig(),
	}
}

// NewWithWriter 创建一个日志对象，使用传入的writer
func NewWithWriter(writer io.Writer) *Logger {
	log := New()
	log.SetWriter(writer)
	return log
}

// Clone clone一个日志对象
func (that *Logger) Clone() *Logger {
	logger := New()
	logger.ctx = that.ctx
	logger.config = that.config
	logger.parent = that
	return logger
}

// 获取将要写入的日志文件名
func (that *Logger) getFilePath(now time.Time) string {
	file, _ := dregex.ReplaceStringFunc(`{.+?}`, that.config.File, func(s string) string {
		return dtime.New(now).Format(strings.Trim(s, "{}"))
	})
	file = dfile.Join(that.config.Path, file)
	return file
}

// 打印日志
func (that *Logger) print(std io.Writer, lead string, values ...interface{}) {

	p := that
	if p.parent != nil {
		p = p.parent
	}

	// 判断日志处理器是否初始化
	if !p.init.Val() && p.init.Cas(false, true) {
		//判断是否需要旋转日志
		if p.config.RotateSize > 0 || p.config.RotateExpire > 0 {
			dtimer.AddOnce(p.config.RotateCheckInterval, p.rotateChecksTimely)
			intlog.Printf("logger rotation initialized: every %s", p.config.RotateCheckInterval.String())
		}
	}
	var (
		now    = time.Now()
		buffer = bytes.NewBuffer(nil)
	)
	//是否需要打印头信息
	if that.config.HeaderPrint {
		timeFormat := ""
		if that.config.Flags&FTimeDate > 0 {
			timeFormat += "2006-01-02 "
		}
		if that.config.Flags&FTimeTime > 0 {
			timeFormat += "15:04:05 "
		}
		if that.config.Flags&FTimeMilli > 0 {
			timeFormat += "15:04:05.000 "
		}
		if len(timeFormat) > 0 {
			buffer.WriteString(now.Format(timeFormat))
		}
		//日志级别的标识
		if len(lead) > 0 {
			buffer.WriteString(lead)
			if len(values) > 0 {
				buffer.WriteByte(' ')
			}
		}
		// 判断日志打印信息
		if that.config.Flags&(FFileLong|FFileShort|FCallerFn) > 0 {
			callerPath := ""
			callerFnName, path, line := ddebug.CallerWithFilter(pathFilterKey, that.config.StSkip)
			//打印栈堆
			if that.config.Flags&FCallerFn > 0 {
				buffer.WriteString(fmt.Sprintf(`[%s] `, callerFnName))
			}
			//打印详细文件名和行号
			if that.config.Flags&FFileLong > 0 {
				callerPath = fmt.Sprintf(`%s:%d: `, path, line)
			}
			//打印端文件名和行号
			if that.config.Flags&FFileShort > 0 {
				callerPath = fmt.Sprintf(`%s:%d: `, dfile.Basename(path), line)
			}
			buffer.WriteString(callerPath)
		}

		//日志前缀
		if len(that.config.Prefix) > 0 {
			buffer.WriteString(that.config.Prefix + " ")
		}
	}

	if that.ctx != nil {
		// Tracing values.
		//spanCtx := trace.SpanContextFromContext(that.ctx)
		//if traceId := spanCtx.TraceID(); traceId.IsValid() {
		//	buffer.WriteString(fmt.Sprintf("{TraceID:%s} ", traceId.String()))
		//}
		// Context values.
		if len(that.config.CtxKeys) > 0 {
			ctxStr := ""
			for _, key := range that.config.CtxKeys {
				if v := that.ctx.Value(key); v != nil {
					if ctxStr != "" {
						ctxStr += ", "
					}
					ctxStr += fmt.Sprintf("%s: %+v", key, v)
				}
			}
			if ctxStr != "" {
				buffer.WriteString(fmt.Sprintf("{%s} ", ctxStr))
			}
		}
	}
	// 将值转换为字符串
	var (
		tempStr  = ""
		valueStr = ""
	)
	for _, v := range values {
		tempStr = dconv.String(v)
		if len(valueStr) > 0 {
			if valueStr[len(valueStr)-1] == '\n' {
				// Remove one blank line(\n\n).
				if tempStr[0] == '\n' {
					valueStr += tempStr[1:]
				} else {
					valueStr += tempStr
				}
			} else {
				valueStr += " " + tempStr
			}
		} else {
			valueStr = tempStr
		}
	}
	buffer.WriteString(valueStr + "\n")

	if that.config.Flags&FAsync > 0 {
		err := asyncPool.Add(func() {
			that.printToWriter(now, std, buffer)
		})
		if err != nil {
			intlog.Error(err)
		}
	} else {
		that.printToWriter(now, std, buffer)
	}
}

//打印到writer
func (that *Logger) printToWriter(now time.Time, std io.Writer, buffer *bytes.Buffer) {

	if that.config.Writer == nil {
		//写入文件
		if that.config.Path != "" {
			that.printToFile(now, buffer)
		}
		//写入标准输出
		if that.config.StdoutPrint {
			if _, err := std.Write(buffer.Bytes()); err != nil {
				intlog.Error(err)
			}
		}
	} else {
		//写入自定义设备
		if _, err := that.config.Writer.Write(buffer.Bytes()); err != nil {
			// panic(err)
			intlog.Error(err)
		}
	}
}

// 打印日志到文件
func (that *Logger) printToFile(now time.Time, buffer *bytes.Buffer) {
	var (
		logFilePath   = that.getFilePath(now)
		memoryLockKey = "dlog.printToFile:" + logFilePath
	)
	//内存锁
	dmlock.Lock(memoryLockKey)
	defer dmlock.Unlock(memoryLockKey)

	//如果日志文件的容量超出了设置的大小，则备份该日志文件，并重新生成新文件
	if that.config.RotateSize > 0 {
		if dfile.Size(logFilePath) > that.config.RotateSize {
			that.rotateFileBySize(now)
		}
	}
	//获取文件的指针
	if file := that.getFilePointer(logFilePath); file == nil {
		intlog.Errorf(`got nil file pointer for: %s`, logFilePath)
	} else {
		if _, err := file.Write(buffer.Bytes()); err != nil {
			intlog.Error(err)
		}
		if err := file.Close(); err != nil {
			intlog.Error(err)
		}
	}
}

// getFilePointer retrieves and returns a file pointer from file pool.
func (that *Logger) getFilePointer(path string) *dfpool.File {
	file, err := dfpool.Open(
		path,
		defaultFileFlags,
		defaultFilePerm,
		defaultFileExpire,
	)
	if err != nil {
		// panic(err)
		intlog.Error(err)
	}
	return file
}

//打印日志到标准输出
func (that *Logger) printStd(lead string, value ...interface{}) {
	that.print(os.Stdout, lead, value...)
}

// 打印错误日志
func (that *Logger) printErr(lead string, value ...interface{}) {
	if that.config.StStatus == 1 {
		if s := that.GetStack(); s != "" {
			value = append(value, "\nStack:\n"+s)
		}
	}
	// In matter of sequence, do not use stderr here, but use the same stdout.
	that.print(os.Stdout, lead, value...)
}

// 格式化
func (that *Logger) format(format string, value ...interface{}) string {
	return fmt.Sprintf(format, value...)
}

// PrintStack 打印栈堆
func (that *Logger) PrintStack(skip ...int) {
	if s := that.GetStack(skip...); s != "" {
		that.Println("Stack:\n" + s)
	} else {
		that.Println()
	}
}

// GetStack 获取栈堆
func (that *Logger) GetStack(skip ...int) string {
	stackSkip := that.config.StSkip
	if len(skip) > 0 {
		stackSkip += skip[0]
	}
	filters := []string{pathFilterKey}
	if that.config.StFilter != "" {
		filters = append(filters, that.config.StFilter)
	}
	return ddebug.StackWithFilters(filters, stackSkip)
}
