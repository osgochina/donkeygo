package dfile

import (
	"github.com/osgochina/donkeygo/util/dconv"
	"io/ioutil"
)

var (
	// Buffer size for reading file content.
	DefaultReadBuffer = 1024
)

// GetContents 获取文件内容
func GetContents(path string) string {
	return dconv.UnsafeBytesToStr(GetBytes(path))
}

func GetBytes(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}
