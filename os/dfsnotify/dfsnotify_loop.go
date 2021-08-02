package dfsnotify

import (
	"context"
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/internal/intlog"
)

//监听事件
func (that *Watcher) watchLoop() {
	go func() {
		for {
			select {
			//收到关闭信号，结束循环
			case <-that.closeChan:
				return
			case ev := <-that.watcher.Events:
				//防止事件并发发送过快，在自定义时间内同样的事件只发送一次
				_, _ = that.cache.SetIfNotExist(ev.String(), func() (interface{}, error) {
					that.events.Push(&Event{
						event:   ev,
						Path:    ev.Name,
						Op:      Op(ev.Op),
						Watcher: that,
					})
					return struct{}{}, nil
				}, repeatEventFilterDuration)
			case err := <-that.watcher.Errors:
				intlog.Error(context.TODO(), err)
			}
		}
	}()
}

// 事件循环
func (that *Watcher) eventLoop() {
	go func() {
		for {
			if v := that.events.Pop(); v != nil {
				event := v.(*Event)
				//获取该路径上注册的所有回调方法
				callbacks := that.getCallbacks(event.Path)
				if len(callbacks) == 0 {
					_ = that.watcher.Remove(event.Path)
					continue
				}
				switch {
				case event.IsRemove():
					if fileExists(event.Path) {
						// 如果该名字的文件还存在，重新把改名字的文件添加到监听中
						if err := that.watcher.Add(event.Path); err != nil {
							intlog.Error(context.TODO(), err)
						} else {
							intlog.Printf(context.TODO(), "fake remove event, watcher re-adds monitor for: %s", event.Path)
						}
						// 如果该名字的文件还存在，就不能认为它是删除，把事件变成改名
						event.Op = RENAME
					}
				case event.IsRename():
					if fileExists(event.Path) {
						//如果该名字的文件还存在,则再次加入监听
						if err := that.watcher.Add(event.Path); err != nil {
							intlog.Error(context.TODO(), err)
						} else {
							intlog.Printf(context.TODO(), "fake rename event, watcher re-adds monitor for: %s", event.Path)
						}
						// 如果该名字的文件还存在，就不能认为它是改名了，把事件更改为修改权限
						event.Op = CHMOD
					}
				case event.IsCreate():
					if fileIsDir(event.Path) {
						// 如果被创建的是一个文件夹，则把该文件夹递归的添加到监听该文件夹的所有事件
						for _, subPath := range fileAllDirs(event.Path) {
							if fileIsDir(subPath) {
								if err := that.watcher.Add(subPath); err != nil {
									intlog.Error(context.TODO(), err)
								} else {
									intlog.Printf(context.TODO(), "folder creation event, watcher adds monitor for: %s", subPath)
								}
							}
						}
					} else {
						// 如果它是一个文件，则把该文件路径添加到监听列表中
						if err := that.watcher.Add(event.Path); err != nil {
							intlog.Error(context.TODO(), err)
						} else {
							intlog.Printf(context.TODO(), "file creation event, watcher adds monitor for: %s", event.Path)
						}
					}
				}
				//开启协程执行要回调的方法
				for _, v := range callbacks {
					go func(callback *Callback) {
						defer func() {
							if err := recover(); err != nil {
								switch err {
								// 如果是回调方法主动退出，则删除该回调
								case callbackExitEventPanicStr:
									that.RemoveCallback(callback.Id)
								default:
									panic(err)
								}
							}
						}()
						callback.Func(event)
					}(v)
				}
			} else {
				break
			}
		}
	}()
}

// 获取改文件的回调方法
func (that *Watcher) getCallbacks(path string) (callbacks []*Callback) {

	//添加自身的回调方法
	if v := that.callbacks.Get(path); v != nil {
		for _, v := range v.(*dlist.List).FrontAll() {
			callback := v.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}
	//获取监听该文件所在文件夹的事件回调方法
	dirPath := fileDir(path)
	if v := that.callbacks.Get(dirPath); v != nil {
		for _, v := range v.(*dlist.List).FrontAll() {
			callback := v.(*Callback)
			callbacks = append(callbacks, callback)
		}
	}

	//递归向上级目录获取监听回调方法
	for {
		parentDirPath := fileDir(dirPath)
		if parentDirPath == dirPath {
			break
		}
		if v := that.callbacks.Get(parentDirPath); v != nil {
			for _, v := range v.(*dlist.List).FrontAll() {
				callback := v.(*Callback)
				if callback.recursive {
					callbacks = append(callbacks, callback)
				}
			}
		}
		dirPath = parentDirPath
	}
	return
}
