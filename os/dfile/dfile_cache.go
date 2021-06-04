package dfile

import (
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/osgochina/donkeygo/os/dcache"
	"github.com/osgochina/donkeygo/os/dcmd"
	"time"
)

//默认缓存超时时间
const dDefaultCacheExpire = time.Minute

var (
	//通过环境变量获取文件超时时间，如果没有获取到，则使用默认的超时时间
	cacheExpire   = dcmd.GetOptWithEnv("dk.dfile.cache", dDefaultCacheExpire).Duration()
	internalCache = dcache.New()
)

// GetContentsWithCache 从缓存中获取文件信息，返回字符串
func GetContentsWithCache(path string, duration ...time.Duration) string {
	return string(GetBytesWithCache(path, duration...))
}

// GetBytesWithCache 获取文件信息，并且缓存
func GetBytesWithCache(path string, duration ...time.Duration) []byte {
	key := cacheKey(path)
	expire := cacheExpire
	if len(duration) > 0 {
		expire = duration[0]
	}
	//获取文件信息，如果不存在该文件，则执行方法创建
	r, _ := internalCache.GetOrSetFuncLock(key, func() (interface{}, error) {
		b := GetBytes(path)
		//如果文件存在，则监听文件变化，有任何变化则删除缓存
		if b != nil {
			_, _ = gfsnotify.Add(path, func(event *gfsnotify.Event) {
				_, _ = internalCache.Remove(key)
				gfsnotify.Exit()
			})
		}
		return b, nil
	}, expire)
	if r != nil {
		return r.([]byte)
	}
	return nil
}

//获取文件缓存的key
func cacheKey(path string) string {
	return "dk.dfile.cache:" + path
}
