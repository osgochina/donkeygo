package dfile

import (
	"github.com/osgochina/donkeygo/text/dstr"
)

// ReplaceFile 替换文件中的指定内容
//从文件中查找内容，并替换，然后写入文件
func ReplaceFile(search, replace, path string) error {
	return PutContents(path, dstr.Replace(GetContents(path), search, replace))
}

// ReplaceFileFunc 通过func处理得到要替换的文件内容，替换path中的内容
func ReplaceFileFunc(f func(path, content string) string, path string) error {
	data := GetContents(path)
	result := f(path, data)
	if len(data) != len(result) && data != result {
		return PutContents(path, result)
	}
	return nil
}

// ReplaceDir 替换文件夹下的文件中的指定内容
func ReplaceDir(search, replace, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = ReplaceFile(search, replace, file); err != nil {
			return err
		}
	}
	return err
}

// ReplaceDirFunc 使用回调函数替换文件夹下文件中的内容
func ReplaceDirFunc(f func(path, content string) string, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = ReplaceFileFunc(f, file); err != nil {
			return err
		}
	}
	return err
}
