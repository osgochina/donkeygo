package dlog

// Print 打印日志
func Print(v ...interface{}) {
	logger.Print(v...)
}

// Printf 打印日志
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// Println 打印一行日志
func Println(v ...interface{}) {
	logger.Println(v...)
}

// Fatal 打印致命日志，并且结束当前进程
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

// Fatalf 打印致命日志，并且结束当前进程
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

// Panic 打印重要日志，并且结束当前协程
func Panic(v ...interface{}) {
	logger.Panic(v...)
}

// Panicf 打印重要日志，并且结束当前协程
func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

// Info 打印一行日志
func Info(v ...interface{}) {
	logger.Info(v...)
}

// Infof 打印一行日志
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

// Debug 打印debug日志
func Debug(v ...interface{}) {
	logger.Debug(v...)
}

// Debugf 打印debug日志
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

// Notice 打印通知日志
func Notice(v ...interface{}) {
	logger.Notice(v...)
}

// Noticef 打印通知日志
func Noticef(format string, v ...interface{}) {
	logger.Noticef(format, v...)
}

// Warning 打印警告日志
func Warning(v ...interface{}) {
	logger.Warning(v...)
}

// Warningf 打印警告日志
func Warningf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}

// Error 打印错误日志
func Error(v ...interface{}) {
	logger.Error(v...)
}

// Errorf 打印错误日志
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

// Critical 打印重要日志
func Critical(v ...interface{}) {
	logger.Critical(v...)
}

// Criticalf 打印重要日志
func Criticalf(format string, v ...interface{}) {
	logger.Criticalf(format, v...)
}
