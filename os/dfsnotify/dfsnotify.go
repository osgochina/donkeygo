package dfsnotify

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/container/dqueue"
	"github.com/osgochina/donkeygo/container/dset"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/internal/intlog"
	"github.com/osgochina/donkeygo/os/dcache"
	"sync"
	"time"
)

// Watcher 监听文件变化
type Watcher struct {
	watcher   *fsnotify.Watcher // 监视器
	events    *dqueue.Queue     // 事件
	cache     *dcache.Cache     // 缓存
	nameSet   *dset.StrSet      //
	callbacks *dmap.StrAnyMap   // 回调函数
	closeChan chan struct{}     //
}

// Callback 回调函数结构体
type Callback struct {
	Id        int
	Func      func(event *Event)
	Path      string
	name      string
	elem      *dlist.Element
	recursive bool
}

// Event 事件接头体
type Event struct {
	event   fsnotify.Event
	Path    string
	Op      Op
	Watcher *Watcher
}

// Op 操作标识
type Op uint32

const (
	CREATE Op = 1 << iota
	WRITE
	REMOVE
	RENAME
	CHMOD
)

const (
	repeatEventFilterDuration = time.Millisecond //防止并发触发事件，在该时间内，多个事件只会触发一次
	callbackExitEventPanicStr = "exit"           // 回调方法主动退出的标识
)

var (
	mu                  sync.Mutex
	defaultWatcher      *Watcher
	callbackIdMap       = dmap.NewIntAnyMap(true)
	callbackIdGenerator = dtype.NewInt()
)

func New() (*Watcher, error) {
	w := &Watcher{
		cache:     dcache.New(),
		events:    dqueue.New(),
		nameSet:   dset.NewStrSet(true),
		closeChan: make(chan struct{}),
		callbacks: dmap.NewStrAnyMap(true),
	}
	if watcher, err := fsnotify.NewWatcher(); err == nil {
		w.watcher = watcher
	} else {
		intlog.Printf("New watcher failed: %v", err)
		return nil, err
	}
	w.watchLoop()
	w.eventLoop()
	return w, nil
}

// 获取默认的监视器
func getDefaultWatcher() (*Watcher, error) {
	mu.Lock()
	defer mu.Unlock()
	if defaultWatcher != nil {
		return defaultWatcher, nil
	}
	var err error
	defaultWatcher, err = New()
	return defaultWatcher, err
}

// Exit 退出文件监听
func Exit() {
	panic(callbackExitEventPanicStr)
}

// Add 添加路径<path>的监听
func Add(path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	w, err := getDefaultWatcher()
	if err != nil {
		return nil, err
	}
	return w.Add(path, callbackFunc, recursive...)
}

// AddOnce 添加路径<path>的监听，使用name区分监听
func AddOnce(name, path string, callbackFunc func(event *Event), recursive ...bool) (callback *Callback, err error) {
	w, err := getDefaultWatcher()
	if err != nil {
		return nil, err
	}
	return w.AddOnce(name, path, callbackFunc, recursive...)
}

// Remove 移除监听
func Remove(path string) error {
	w, err := getDefaultWatcher()
	if err != nil {
		return err
	}
	return w.Remove(path)
}

// RemoveCallback 移除回调函数
func RemoveCallback(callbackId int) error {
	w, err := getDefaultWatcher()
	if err != nil {
		return err
	}
	callback := (*Callback)(nil)
	if r := callbackIdMap.Get(callbackId); r != nil {
		callback = r.(*Callback)
	}
	if callback == nil {
		return errors.New(fmt.Sprintf(`callback for id %d not found`, callbackId))
	}
	w.RemoveCallback(callbackId)
	return nil
}
