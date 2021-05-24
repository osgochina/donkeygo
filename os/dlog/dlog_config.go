package dlog

import "io"

// SetConfig 设置配置信息
func SetConfig(config Config) error {
	return logger.SetConfig(config)
}

// SetConfigWithMap 通过map格式设置配置信息
func SetConfigWithMap(m map[string]interface{}) error {
	return logger.SetConfigWithMap(m)
}

// SetPath 设置存储目录
func SetPath(path string) error {
	return logger.SetPath(path)
}

// GetPath 获取存储目录
func GetPath() string {
	return logger.GetPath()
}

// SetFile 设置存储文件
func SetFile(pattern string) {
	logger.SetFile(pattern)
}

// SetLevel 设置日志级别
func SetLevel(level int) {
	logger.SetLevel(level)
}

// GetLevel 获取日志级别
func GetLevel() int {
	return logger.GetLevel()
}

// SetWriter 设置自定义writer
func SetWriter(writer io.Writer) {
	logger.SetWriter(writer)
}

// GetWriter 获取当前的writer对象
func GetWriter() io.Writer {
	return logger.GetWriter()
}

// SetDebug 开启debug
func SetDebug(debug bool) {
	logger.SetDebug(debug)
}

// SetAsync 设置异步写文件
func SetAsync(enabled bool) {
	logger.SetAsync(enabled)
}

// SetStdoutPrint 设置标准输出
func SetStdoutPrint(enabled bool) {
	logger.SetStdoutPrint(enabled)
}

// SetHeaderPrint 设置日志头输出
func SetHeaderPrint(enabled bool) {
	logger.SetHeaderPrint(enabled)
}

// SetPrefix 设置日志前缀
func SetPrefix(prefix string) {
	logger.SetPrefix(prefix)
}

// SetFlags 设置日志扩展标识
func SetFlags(flags int) {
	logger.SetFlags(flags)
}

// GetFlags 获取日志扩展标识
func GetFlags() int {
	return logger.GetFlags()
}

// SetCtxKeys 设置上下文key
func SetCtxKeys(keys ...interface{}) {
	logger.SetCtxKeys(keys...)
}

// GetCtxKeys 获取上下文key
func GetCtxKeys() []interface{} {
	return logger.GetCtxKeys()
}

// PrintStack 打印栈堆层数
func PrintStack(skip ...int) {
	logger.PrintStack(skip...)
}

// GetStack 获取栈信息
func GetStack(skip ...int) string {
	return logger.GetStack(skip...)
}

// SetStack 开启栈堆打印
func SetStack(enabled bool) {
	logger.SetStack(enabled)
}

// SetLevelStr 通过字符串设置日志级别
func SetLevelStr(levelStr string) error {
	return logger.SetLevelStr(levelStr)
}

// SetLevelPrefix 设置日志级别前缀
func SetLevelPrefix(level int, prefix string) {
	logger.SetLevelPrefix(level, prefix)
}

// SetLevelPrefixes 设置日志级别的前缀
func SetLevelPrefixes(prefixes map[int]string) {
	logger.SetLevelPrefixes(prefixes)
}

// GetLevelPrefix 获取日志级别的前缀
func GetLevelPrefix(level int) string {
	return logger.GetLevelPrefix(level)
}
