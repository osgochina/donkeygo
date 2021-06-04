package dfile

import (
	"github.com/osgochina/donkeygo/os/dcmd"
	"time"
)

//默认缓存超时时间
const dDefaultCacheExpire = time.Minute

var (
	//通过环境变量获取文件超时时间，如果没有获取到，则使用默认的超时时间
	cacheExpire = dcmd.GetOptWithEnv("dk.dfile.cache", dDefaultCacheExpire).Duration()
)
