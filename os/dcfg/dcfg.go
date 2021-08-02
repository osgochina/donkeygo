package dcfg

import (
	"context"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dcmd"
)

// Config 配置信息结构体
type Config struct {
	defaultName   string           //默认的配置文件名字
	searchPaths   *darray.StrArray //配置文件的搜索路径
	jsonMap       *dmap.StrAnyMap  //把配置文件转换成json对象后存储在这里
	violenceCheck bool             //
}

const (
	DefaultName       = "config"             //
	DefaultConfigFile = "config.toml"        //默认的配置文件名字
	cmdEnvKey         = "dk.dcfg"            //通过环境变量传入配置文件路径使用的key
	errorPrintKey     = "dk.dcfg.errorprint" // 通过环境变量传入是否打开错误信息配置
)

var (
	// 支持以下格式的配置文件解析
	supportedFileTypes = []string{"toml", "yaml", "yml", "json", "ini", "xml"}
	//资源文件夹的路径名，会去以下这些文件夹下查找配置文件
	resourceTryFiles = []string{"", "/", "config/", "config", "/config", "/config/"}

	//配置文件实例
	instances = dmap.NewStrAnyMap(true)
	//自定制的配置文件内容
	customConfigContentMap = dmap.NewStrStrMap(true)
)

// SetContent 写入配置文件内容到指定的文件映射对象，注意不会改变源文件
// 如果不传入文件名，则写入到默认配置文件中
func SetContent(content string, file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	instances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				v.(*Config).jsonMap.Remove(name)
			}
		}
		customConfigContentMap.Set(name, content)
	})
}

// GetContent 获取配置文件内容
func GetContent(file ...string) string {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}

	return customConfigContentMap.Get(name)
}

// RemoveContent 设置指定的配置文件内容
func RemoveContent(file ...string) {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	}
	instances.LockFunc(func(m map[string]interface{}) {
		if customConfigContentMap.Contains(name) {
			for _, v := range m {
				v.(*Config).jsonMap.Remove(name)
			}
		}
		customConfigContentMap.Remove(name)
	})

	intlog.Printf(context.TODO(), `RemoveContent: %s`, name)
}

// ClearContent 清空所有的配置文件内容
func ClearContent() {
	customConfigContentMap.Clear()
	// Clear cache for all instances.
	instances.LockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			v.(*Config).jsonMap.Clear()
		}
	})

	intlog.Print(context.TODO(), `RemoveConfig`)
}

// 是否需要打印错误信息
func errorPrint() bool {
	return dcmd.GetOptWithEnv(errorPrintKey, true).Bool()
}
