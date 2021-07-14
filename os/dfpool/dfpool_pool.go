package dfpool

import (
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/osgochina/donkeygo/container/dpool"
	"github.com/osgochina/donkeygo/container/dtype"
	"os"
	"time"
)

// New 新建并返回一个文件指针池
// Note the expiration logic:
// ttl = 0 : not expired;
// ttl < 0 : immediate expired after use;
// ttl > 0 : timeout expired;
// It is not expired in default.
func New(path string, flag int, perm os.FileMode, ttl ...time.Duration) *Pool {
	var fpTTL time.Duration
	if len(ttl) > 0 {
		fpTTL = ttl[0]
	}
	p := &Pool{
		id:   dtype.NewInt(),
		ttl:  fpTTL,
		init: dtype.NewBool(),
	}
	p.pool = newFilePool(p, path, flag, perm, fpTTL)
	return p
}

func newFilePool(p *Pool, path string, flag int, perm os.FileMode, ttl time.Duration) *dpool.Pool {
	pool := dpool.New(ttl, func() (interface{}, error) {
		file, err := os.OpenFile(path, flag, perm)
		if err != nil {
			return nil, err
		}
		return &File{
			File: file,
			pid:  p.id.Val(),
			pool: p,
			flag: flag,
			perm: perm,
			path: path,
		}, nil
	}, func(i interface{}) {
		_ = i.(*File).Close()
	})
	return pool
}

// File 从文件池中获取文件的指针
func (that *Pool) File() (*File, error) {

	if v, err := that.pool.Get(); err != nil {
		return nil, err
	} else {
		f := v.(*File)
		f.stat, err = os.Stat(f.path)
		//文件支持自动创建
		if f.flag&os.O_CREATE > 0 {
			if os.IsNotExist(err) {
				if f.File, err = os.OpenFile(f.path, f.flag, f.perm); err != nil {
					return nil, err
				} else {
					if f.stat, err = f.File.Stat(); err != nil {
						return nil, err
					}
				}
			}
		}
		//如果需要截断文件
		if f.flag&os.O_TRUNC > 0 {
			if f.stat.Size() > 0 {
				if err = f.Truncate(0); err != nil {
					return nil, err
				}
			}
		}
		//是否要追加写入
		if f.flag&os.O_APPEND > 0 {
			// offset 开始偏移量，也就是代表需要移动偏移的字节数
			// whence：给offset参数一个定义，表示要从哪个位置开始偏移；0代表从文件开头开始算起，1代表从当前位置开始算起，2代表从文件末尾算起。
			// 这里表示从文件末尾开始写入
			if _, err = f.Seek(0, 2); err != nil {
				return nil, err
			}
		} else {
			// 表示从头覆盖写
			if _, err = f.Seek(0, 0); err != nil {
				return nil, err
			}
		}

		//如果文件池没有初始化，则对文件加入事件监听
		if !that.init.Val() && that.init.Cas(false, true) {
			_, _ = gfsnotify.Add(f.path, func(event *gfsnotify.Event) {
				//如果要监听的文件被删除了或者改名了，则把该文件指针池清除了
				if event.IsRemove() || event.IsRename() {
					//删除老的池子
					that.id.Add(1)
					that.pool.Clear()
					that.id.Add(1)
				}
			}, false)
		}

		return f, nil
	}
}

// Close 关闭当前文件指针池
func (that *Pool) Close() {
	that.pool.Close()
}
