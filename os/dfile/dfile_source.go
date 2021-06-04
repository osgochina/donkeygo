package dfile

import (
	"github.com/osgochina/donkeygo/text/dregex"
	"github.com/osgochina/donkeygo/text/dstr"
	"os"
	"runtime"
	"strings"
)

var (
	// goRootForFilter is used for stack filtering purpose.
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

//MainPkgPath 返回main包的绝对文件路径，
//其中主要包含入口功能。
//它只在开发环境中可用。
//注1:仅对源代码开发环境有效，
//只对生成此可执行文件的系统有效。
//注意2:当该方法第一次被调用时，如果它是在一个异步的goroutine中，该方法可能不会获得主包路径。
func MainPkgPath() string {

	if goRootForFilter == "" {
		return ""
	}
	path := mainPkgPath.Val()
	if path != "" {
		return path
	}
	var lastFile string
	for i := 1; i < 10000; i++ {
		pc, file, _, ok := runtime.Caller(i)
		if ok {
			if goRootForFilter != "" &&
				len(file) >= len(goRootForFilter) &&
				file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			if Ext(file) != ".go" {
				continue
			}
			lastFile = file

			//检查它是否在包初始化函数中被调用，在这里它不能检索主包路径，所以它只是返回，可以进行下一次检查。
			if fn := runtime.FuncForPC(pc); fn != nil {
				array := dstr.Split(fn.Name(), ".")
				if array[0] != "main" {
					continue
				}
			}
			if dregex.IsMatchString(`package\s+main\s+`, GetContents(file)) {
				mainPkgPath.Set(Dir(file))
				return Dir(file)
			}
		} else {
			break
		}
	}

	//如果它仍然找不到包main的路径，
	//它递归地搜索最后一个go文件的目录及其父目录。
	//业务项目的单元测试用例通常是必要的。
	if lastFile != "" {
		for path = Dir(lastFile); len(path) > 1 && Exists(path) && path[len(path)-1] != os.PathSeparator; {
			files, _ := ScanDir(path, "*.go")
			for _, v := range files {
				if dregex.IsMatchString(`package\s+main\s+`, GetContents(v)) {
					mainPkgPath.Set(path)
					return path
				}
			}
			path = Dir(path)
		}
	}

	return ""
}
