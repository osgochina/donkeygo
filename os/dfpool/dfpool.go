package dfpool

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/container/dpool"
	"github.com/osgochina/donkeygo/container/dtype"
	"os"
	"time"
)

// Pool 文件指针池子
type Pool struct {
	id   *dtype.Int    // 文件池id
	pool *dpool.Pool   // 对象复用池
	init *dtype.Bool   // 是否初始化
	ttl  time.Duration // 过期时间
}

// File 文件对象
type File struct {
	*os.File             //基础的文件句柄
	stat     os.FileInfo // 文件状态
	pid      int         //所属的pool id
	pool     *Pool       // 所属的pool
	flag     int         // 打开文件的flag
	perm     os.FileMode // 打开文件的Permission
	path     string      //文件的绝对路径
}

var (
	// 全局的文件指针池
	pools = dmap.NewStrAnyMap(true)
)
