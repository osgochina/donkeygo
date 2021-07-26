package dfsnotify

import (
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/internal/intlog"
)

// Add 添加文件路径的监视，并且绑定回调函数
// 如果传入的路径是目录，并且recursive为true，则执行递归监听
// recursive默认是true
func (that *Watcher) Add(path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	return that.AddOnce("", path, callbackFunc, recursive...)
}

// AddOnce 使用唯一的<name>添加path的监控,相同name多次监控，只会成功一次
func (that *Watcher) AddOnce(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	that.nameSet.AddIfNotExistFuncLock(name, func() bool {
		// 添加监听路径
		callback, err = that.addWithCallbackFunc(name, path, callbackFunc, recursive...)
		if err != nil {
			return false
		}
		// 如果传入的路径是文件夹，并且需要递归监听
		if fileIsDir(path) && (len(recursive) == 0 || recursive[0]) {
			// 获取当前文件夹下的所有文件夹地址
			for _, subPath := range fileAllDirs(path) {
				if fileIsDir(subPath) {
					if err := that.watcher.Add(subPath); err != nil {
						intlog.Error(err)
					} else {
						intlog.Printf("watcher adds monitor for: %s", subPath)
					}
				}
			}
		}

		if name == "" {
			return false
		}
		return true
	})
	return callback, err
}

// 将路径添加到底层监视器，创建并返回回调对象
// 需要注意的是，如果它多次调用相同的' path '，最新的路径将覆盖之前的路径。
func (that *Watcher) addWithCallbackFunc(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {

	//判断要监听的路径是否存在
	if t := fileRealPath(path); t == "" {
		return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
	} else {
		path = t
	}
	callback = &Callback{
		Id:        callbackIdGenerator.Add(1),
		Func:      callbackFunc,
		Path:      path,
		name:      name,
		recursive: true,
	}
	if len(recursive) > 0 {
		callback.recursive = recursive[0]
	}
	//把callback方法加入到list中
	that.callbacks.LockFunc(func(m map[string]interface{}) {
		list := (*dlist.List)(nil)
		if v, ok := m[path]; !ok {
			list = dlist.New(true)
			m[path] = list
		} else {
			list = v.(*dlist.List)
		}
		callback.elem = list.PushBack(callback)
	})
	// 添加path的监听
	if err = that.watcher.Add(path); err != nil {
		intlog.Error(err)
	} else {
		intlog.Printf("watcher adds monitor for: %s", path)
	}
	// 把callback加入到map中
	callbackIdMap.Set(callback.Id, callback)

	return
}

// Close 关闭监控器
func (that *Watcher) Close() {
	that.events.Close()
	if err := that.watcher.Close(); err != nil {
		intlog.Error(err)
	}
	close(that.closeChan)
}

// RemoveCallback 删除回调函数
func (that *Watcher) RemoveCallback(callbackId int) {
	callback := (*Callback)(nil)
	if r := callbackIdMap.Get(callbackId); r != nil {
		callback = r.(*Callback)
	}
	if callback != nil {
		if r := that.callbacks.Get(callback.Path); r != nil {
			r.(*dlist.List).Remove(callback.elem)
		}
		callbackIdMap.Remove(callbackId)
		if callback.name != "" {
			that.nameSet.Remove(callback.name)
		}
	}
}

// Remove 移除path的监听
func (that *Watcher) Remove(path string) error {
	// Firstly remove the callbacks of the path.
	if r := that.callbacks.Remove(path); r != nil {
		list := r.(*dlist.List)
		for {
			if r := list.PopFront(); r != nil {
				callbackIdMap.Remove(r.(*Callback).Id)
			} else {
				break
			}
		}
	}
	// Secondly remove monitor of all sub-files which have no callbacks.
	if subPaths, err := fileScanDir(path, "*", true); err == nil && len(subPaths) > 0 {
		for _, subPath := range subPaths {
			if that.checkPathCanBeRemoved(subPath) {
				if err := that.watcher.Remove(subPath); err != nil {
					intlog.Error(err)
				}
			}
		}
	}
	// Lastly remove the monitor of the path from underlying monitor.
	return that.watcher.Remove(path)
}

// 判断path是否可以移除监听
func (that *Watcher) checkPathCanBeRemoved(path string) bool {
	// Firstly check the callbacks in the watcher directly.
	if v := that.callbacks.Get(path); v != nil {
		return false
	}
	// Secondly check its parent whether has callbacks.
	dirPath := fileDir(path)
	if v := that.callbacks.Get(dirPath); v != nil {
		for _, c := range v.(*dlist.List).FrontAll() {
			if c.(*Callback).recursive {
				return false
			}
		}
		return false
	}
	// Recursively check its parent.
	parentDirPath := ""
	for {
		parentDirPath = fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if v := that.callbacks.Get(parentDirPath); v != nil {
			for _, c := range v.(*dlist.List).FrontAll() {
				if c.(*Callback).recursive {
					return false
				}
			}
			return false
		}
		dirPath = parentDirPath
	}
	return true
}
