package dcfg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/gogf/gf/os/gres"
	"github.com/gogf/gf/os/gspath"
	"github.com/gogf/gf/util/gmode"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dcmd"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/text/dstr"
)

// New 创建一个配置文件对象，可以传入指定的文件名
func New(file ...string) *Config {

	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	} else {
		//从命令行或环境自定义缺省配置文件名称。
		customFile := dcmd.GetOptWithEnv(fmt.Sprintf("%s.file", cmdEnvKey)).String()
		if len(customFile) > 0 {
			name = customFile
		}
	}
	c := &Config{
		defaultName: name,
		searchPaths: darray.NewStrArray(true),
		jsonMap:     dmap.NewStrAnyMap(true),
	}
	customPath := dcmd.GetOptWithEnv(fmt.Sprintf("%s.path", cmdEnvKey)).String()
	if len(customPath) > 0 {
		if dfile.Exists(customPath) {
			_ = c.SetPath(customPath)
		} else {
			if errorPrint() {
				dlog.Errorf("[dcfg] Configuration directory path does not exist: %s", customPath)
			}
		}
	} else {
		//把当前working目录加入搜索范围
		if err := c.AddPath(dfile.Pwd()); err != nil {
			intlog.Error(err)
		}

		// 把当前代码的main包所在的目录加入到搜索范围
		if mainPath := dfile.MainPkgPath(); mainPath != "" && dfile.Exists(mainPath) {
			if err := c.AddPath(mainPath); err != nil {
				intlog.Error(err)
			}
		}

		// 把当前程序所在的目录放入搜索范围
		if selfPath := dfile.SelfDir(); selfPath != "" && dfile.Exists(selfPath) {
			if err := c.AddPath(selfPath); err != nil {
				intlog.Error(err)
			}
		}

	}
	return c
}

// Instance 通过名字获取全局配置实例.
// 默认是使用"toml"后缀的配置文件
func Instance(name ...string) *Config {
	key := DefaultName
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	//加锁获取，如果存在则直接返回，如果不存在则调用方法生成，然后返回
	return instances.GetOrSetFuncLock(key, func() interface{} {
		c := New()
		if key != DefaultName || !c.Available() {
			for _, fileType := range supportedFileTypes {
				file := fmt.Sprintf(`%s.%s`, key, fileType)
				if c.Available(file) {
					c.SetFileName(file)
					break
				}
			}
		}
		return c
	}).(*Config)
}

// SetPath 设置配置文件的搜索路径，可以是绝对路径也可以是相对路径
func (that *Config) SetPath(path string) error {

	realPath, err := that.realPath(path)
	if err != nil {
		return err
	}
	//重复检查以下该路径是否已经生成过
	if that.searchPaths.Search(realPath) != -1 {
		return nil
	}
	// 清空配置信息，等待生成
	that.jsonMap.Clear()
	that.searchPaths.Clear()
	that.searchPaths.Append(realPath)
	intlog.Print("SetPath:", realPath)
	return nil
}

// SetViolenceCheck 设置是否进行分级冲突检查。
//当键名中有级别符号时，需要启用此功能。
//默认是关闭的。
//注意，打开这个特性的代价是非常昂贵的，不建议在键名中允许分隔符。最好在应用程序端避免这种情况。
func (that *Config) SetViolenceCheck(check bool) {
	that.violenceCheck = check
	that.Clear()
}

// AddPath 添加路径到搜索路径列表
func (that *Config) AddPath(path string) error {
	realPath, err := that.realPath(path)
	if err != nil {
		return err
	}
	// 检查路径是否已经存在
	if that.searchPaths.Search(realPath) != -1 {
		return nil
	}
	that.searchPaths.Append(realPath)
	intlog.Print("AddPath:", realPath)
	return nil
}

// 生成真实可用的path
func (that *Config) realPath(path string) (string, error) {
	isDir := true
	realPath := ""

	// 判断需要设置的文件夹路径是否被资源化缓存
	if file := gres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		//真实路径
		realPath = dfile.RealPath(path)
		if realPath == "" {
			//去搜索路径找到合适的并且存在的文件夹
			that.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := gspath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = dfile.IsDir(realPath)
		}
	}
	//如果路径还不存在，生成错误信息并返回这些错误信息
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		if that.searchPaths.Len() > 0 {

			buffer.WriteString(fmt.Sprintf("[dcfg] SetPath failed: cannot find directory \"%s\" in following paths:", path))
			that.searchPaths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`[dcfg] SetPath failed: path "%s" does not exist`, path))
		}
		err := errors.New(buffer.String())
		if errorPrint() {
			dlog.Error(err)
		}
		return "", err
	}

	//如果传入的路径不是文件夹，则生成错误信息返回
	if !isDir {
		err := fmt.Errorf(`[dcfg] SetPath failed: path "%s" should be directory type`, path)
		if errorPrint() {
			dlog.Error(err)
		}
		return "", err
	}
	return realPath, nil
}

// SetFileName 设置配置文件名
func (that *Config) SetFileName(name string) *Config {
	that.defaultName = name
	return that
}

// GetFileName 获取配置文件名
func (that *Config) GetFileName() string {
	return that.defaultName
}

// Available 检查指定的配置文件中的配置是否可用
func (that *Config) Available(file ...string) bool {
	var name string
	if len(file) > 0 && file[0] != "" {
		name = file[0]
	} else {
		name = that.defaultName
	}
	if path, _ := that.GetFilePath(name); path != "" {
		return true
	}
	if GetContent(name) != "" {
		return true
	}
	return false
}

func (that *Config) GetFilePath(file ...string) (path string, err error) {
	name := that.defaultName
	if len(file) > 0 {
		name = file[0]
	}
	//如果资源不为空
	if !gres.IsEmpty() {
		// 去查一下指定的配置文件是否在资源文件夹下面
		for _, v := range resourceTryFiles {
			if file := gres.Get(v + name); file != nil {
				path = file.Name()
				return
			}
		}
		// 去搜索路径下查找配置文件是否在资源文件夹中
		that.searchPaths.RLockFunc(func(array []string) {
			for _, prefix := range array {
				for _, v := range resourceTryFiles {
					if file := gres.Get(prefix + v + name); file != nil {
						path = file.Name()
						return
					}
				}
			}
		})
	}
	that.autoCheckAndAddMainPkgPathToSearchPaths()

	// 在搜索目录列表中查找是否存在指定的配置文件
	that.searchPaths.RLockFunc(func(array []string) {
		for _, prefix := range array {
			prefix = dstr.TrimRight(prefix, `\/`)
			//查到搜索目录
			if path, _ = gspath.Search(prefix, name); path != "" {
				return
			}
			//查找搜索路径下的目录中的config目录
			if path, _ = gspath.Search(prefix+dfile.Separator+"config", name); path != "" {
				return
			}
		}
	})

	if path == "" {
		var (
			buffer = bytes.NewBuffer(nil)
		)
		if that.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf(`[dcfg] cannot find config file "%s" in resource manager or the following paths:`, name))
			that.searchPaths.RLockFunc(func(array []string) {
				index := 1
				for _, v := range array {
					v = dstr.TrimRight(v, `\/`)
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, v))
					index++
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, v+dfile.Separator+"config"))
					index++
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf("[dcfg] cannot find config file \"%s\" with no path configured", name))
		}
		err = derror.New(buffer.String())
	}
	return
}

//判断当前是否处于开发环境，如果处于开发环境，则自动把main包所在的目录添加到配置搜索目录
func (that *Config) autoCheckAndAddMainPkgPathToSearchPaths() {
	if gmode.IsDevelop() {
		mainPkgPath := dfile.MainPkgPath()
		if mainPkgPath != "" {
			if !that.searchPaths.Contains(mainPkgPath) {
				that.searchPaths.Append(mainPkgPath)
			}
		}
	}
}

// 获取配置文件的json格式
func (that *Config) getJson(file ...string) *gjson.Json {
	var name string
	if len(file) > 0 && file[0] != "" {
		name = file[0]
	} else {
		name = that.defaultName
	}

	r := that.jsonMap.GetOrSetFuncLock(name, func() interface{} {
		var (
			err      error
			content  string
			filePath string
		)
		isFromConfigContent := true
		if content = GetContent(name); content == "" {
			//内存中不存在
			isFromConfigContent = false
			//获取配置文件真实地址
			filePath, err = that.GetFilePath(name)
			if err != nil && errorPrint() {
				dlog.Error(err)
			}
			if filePath == "" {
				return nil
			}
			// 先从资源池中获取文件内容
			if file := gres.Get(filePath); file != nil {
				content = string(file.Content())
			} else {
				//直接获取原始文件内容
				content = dfile.GetContents(filePath)
			}
		}

		var (
			j *gjson.Json
		)
		//获取配置文件格式
		dataType := dfile.ExtName(name)
		// 判断gjson是否支持该配置文件类型,并且载入文件到内存
		if gjson.IsValidDataType(dataType) && !isFromConfigContent {
			j, err = gjson.LoadContentType(dataType, content, true)
		} else {
			j, err = gjson.LoadContent(content, true)
		}
		if err != nil {
			//如果
			if errorPrint() {
				if filePath != "" {
					dlog.Criticalf(`[dcfg] load config file "%s" failed: %s`, filePath, err.Error())
				} else {
					dlog.Criticalf(`[dcfg] load configuration failed: %s`, err.Error())
				}
			}
			return nil
		}
		//设置数据分层的暴力检查
		j.SetViolenceCheck(that.violenceCheck)

		//文件不存在资源池，添加文件变更事件，实时更新配置文件
		if filePath != "" && !gres.Contains(filePath) {
			//文件有任何变更，则从内存中删除该文件
			_, err = gfsnotify.Add(filePath, func(event *gfsnotify.Event) {
				that.jsonMap.Remove(name)
			})
			if err != nil && errorPrint() {
				dlog.Error(err)
			}
		}
		return j
	})
	if r != nil {
		return r.(*gjson.Json)
	}
	return nil
}
