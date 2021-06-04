package dfile

import (
	"bufio"
	"github.com/osgochina/donkeygo/util/dconv"
	"io"
	"io/ioutil"
	"os"
)

var (
	// DefaultReadBuffer 读取文件的缓冲区大小
	DefaultReadBuffer = 1024
)

// GetContents 获取文件内容，返回字符串信息
func GetContents(path string) string {
	return dconv.UnsafeBytesToStr(GetBytes(path))
}

// GetBytes 获取文件内容，返回byte数组
func GetBytes(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

//写入内容到文件
func putContents(path string, data []byte, flag int, perm os.FileMode) error {
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return err
		}
	}
	//使用特殊的标识打开文件
	f, err := OpenWithFlagPerm(path, flag, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	if n, e := f.Write(data); e != nil {
		return e
	} else if n < len(data) {
		return io.ErrShortWrite
	}
	return nil
}

// Truncate 截断文件到指定的大小
func Truncate(path string, size int) error {
	return os.Truncate(path, int64(size))
}

// PutContents 写入字符串内容到文件，文件不存在则创建，如果存在则从头开始写
func PutContents(path string, content string) error {
	return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultPermOpen)
}

// PutContentsAppend 追加写入字符串内容到文件末尾，文件不存在则创建，已存在则写到末尾
func PutContentsAppend(path string, content string) error {
	return putContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, DefaultPermOpen)
}

// PutBytes 写入字节数组内容到文件，文件不存在则创建，如果存在则从头开始写
func PutBytes(path string, content []byte) error {
	return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultPermOpen)
}

// PutBytesAppend 追加写入字节数组内容到文件末尾，文件不存在则创建，已存在则写到末尾
func PutBytesAppend(path string, content []byte) error {
	return putContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, DefaultPermOpen)
}

// GetNextCharOffset 从指定位置开始读取数据，碰到与char相等的字节返回它所在的位置
func GetNextCharOffset(reader io.ReaderAt, char byte, start int64) int64 {
	buffer := make([]byte, DefaultReadBuffer)
	offset := start
	for {
		if n, err := reader.ReadAt(buffer, offset); n > 0 {
			for i := 0; i < n; i++ {
				if buffer[i] == char {
					return int64(i) + offset
				}
			}
			offset += int64(n)
		} else if err != nil {
			break
		}
	}
	return -1
}

// GetNextCharOffsetByPath 获取指定文件中char所在的位置，从start开始查找
func GetNextCharOffsetByPath(path string, char byte, start int64) int64 {
	if f, err := OpenWithFlagPerm(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetNextCharOffset(f, char, start)
	}
	return -1
}

// GetBytesTilChar 获取从指定字符到结尾所有的字符数组
func GetBytesTilChar(reader io.ReaderAt, char byte, start int64) ([]byte, int64) {
	if offset := GetNextCharOffset(reader, char, start); offset != -1 {
		return GetBytesByTwoOffsets(reader, start, offset+1), offset
	}
	return nil, -1
}

// GetBytesTilCharByPath 获取文件中从start位置开始，找到的第一个char字符，一直到文件结尾的内容
func GetBytesTilCharByPath(path string, char byte, start int64) ([]byte, int64) {
	if f, err := OpenWithFlagPerm(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetBytesTilChar(f, char, start)
	}
	return nil, -1
}

// GetBytesByTwoOffsets 获取指定区间的字符内容
func GetBytesByTwoOffsets(reader io.ReaderAt, start int64, end int64) []byte {
	buffer := make([]byte, end-start)
	if _, err := reader.ReadAt(buffer, start); err != nil {
		return nil
	}
	return buffer
}

// GetBytesByTwoOffsetsByPath 获取文件中从start开始到end结束区间的内容
func GetBytesByTwoOffsetsByPath(path string, start int64, end int64) []byte {
	if f, err := OpenWithFlagPerm(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetBytesByTwoOffsets(f, start, end)
	}
	return nil
}

// ReadLines 逐行读取文件内容，读取一行则回调一次callback方法
func ReadLines(file string, callback func(text string) error) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = callback(scanner.Text()); err != nil {
			return err
		}
	}
	return nil
}

// ReadByteLines 逐行读取文件内容，读取一行回调一次callback方法
func ReadByteLines(file string, callback func(bytes []byte) error) error {
	return ReadLinesBytes(file, callback)
}

// ReadLinesBytes 逐行读取文件内容，读取一行回调一次callback方法
func ReadLinesBytes(file string, callback func(bytes []byte) error) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = callback(scanner.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
