package dfile

import (
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/text/dstr"
	"os"
	"path/filepath"
	"sort"
)

const (
	// 目录扫描的最大递归深度
	maxScanDepth = 100000
)

// ScanDir 扫描文件夹，找到匹配的文件夹
func ScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, nil)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFunc 扫描文件夹，找到匹配的文件夹，并且执行回调函数
func ScanDirFunc(path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	list, err := doScanDir(0, path, pattern, recursive, handler)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFile 遍历返回文件夹中的文件
func ScanDirFile(path string, pattern string, recursive ...bool) ([]string, error) {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, func(path string) string {
		if IsDir(path) {
			return ""
		}
		return path
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFileFunc 遍历文件夹中的文件，并执行回调函数
func ScanDirFileFunc(path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	list, err := doScanDir(0, path, pattern, recursive, func(path string) string {
		if IsDir(path) {
			return ""
		}
		return handler(path)
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

//递归扫描遍历文件夹的文件
func doScanDir(depth int, path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	if depth >= maxScanDepth {
		return nil, derror.Newf("directory scanning exceeds max recursive depth: %d", maxScanDepth)
	}
	list := ([]string)(nil)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	var (
		filePath = ""
		patterns = dstr.SplitAndTrim(pattern, ",")
	)
	for _, name := range names {
		filePath = path + Separator + name
		if IsDir(filePath) && recursive {
			array, _ := doScanDir(depth+1, filePath, pattern, true, handler)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		// Handler filtering.
		if handler != nil {
			filePath = handler(filePath)
			if filePath == "" {
				continue
			}
		}
		// If it meets pattern, then add it to the result list.
		for _, p := range patterns {
			if match, err := filepath.Match(p, name); err == nil && match {
				filePath = Abs(filePath)
				if filePath != "" {
					list = append(list, filePath)
				}
			}
		}
	}
	return list, nil
}
