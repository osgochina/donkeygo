package dfpool

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Open 打开一个文件，并创建该文件的指针池，返回该文件的句柄
func Open(path string, flag int, perm os.FileMode, ttl ...time.Duration) (file *File, err error) {
	var fpTTL time.Duration
	if len(ttl) > 0 {
		fpTTL = ttl[0]
	}
	pool := pools.GetOrSetFuncLock(
		fmt.Sprintf("%s&%d&%d&%d", path, flag, fpTTL, perm),
		func() interface{} {
			return New(path, flag, perm, fpTTL)
		},
	).(*Pool)
	return pool.File()
}

func (that *File) Stat() (os.FileInfo, error) {
	if that.stat == nil {
		return nil, errors.New("file stat is empty")
	}
	return that.stat, nil
}

// Close 关闭文件指针对象池
func (that *File) Close() error {
	if that.pid == that.pool.id.Val() {
		return that.pool.pool.Put(that)
	}
	return nil
}
