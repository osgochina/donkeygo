package dfile

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/container/darray"
)

// Search 在以下路径中按名字有限搜索文件
// prioritySearchPaths, Pwd()、SelfDir()、MainPkgPath().
// 如果找到，则返回绝对路径，如果未找到则返回空字符串，并报错
func Search(name string, prioritySearchPaths ...string) (realPath string, err error) {

	//获取指定文件的真是路径
	realPath = RealPath(name)
	if realPath != "" {
		return
	}
	//需要扫描的文件夹
	array := darray.NewStrArray()
	array.Append(prioritySearchPaths...)
	array.Append(Pwd(), SelfDir())
	if path := MainPkgPath(); path != "" {
		array.Append(path)
	}
	// 删除重复的目录
	array.Unique()

	// 搜索目录下是否存在该文件
	array.RLockFunc(func(array []string) {
		path := ""
		for _, v := range array {
			path = RealPath(v + Separator + name)
			if path != "" {
				realPath = path
				break
			}
		}
	})

	//全部搜索完毕，还没搜索到要找的文件
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		buffer.WriteString(fmt.Sprintf("cannot find file/folder \"%s\" in following paths:", name))
		array.RLockFunc(func(array []string) {
			for k, v := range array {
				buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
			}
		})
		err = errors.New(buffer.String())
	}

	return
}
