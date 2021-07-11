package dmd5

import (
	"crypto/md5"
	"fmt"
	"github.com/osgochina/donkeygo/util/dconv"
	"io"
	"os"
)

// Encrypt 把任何对象加密成md5
func Encrypt(data interface{}) (encrypt string, err error) {
	return EncryptBytes(dconv.Bytes(data))
}

// MustEncrypt 把任何对象加密成md5,如果加密失败则直接panic
func MustEncrypt(data interface{}) string {
	result, err := Encrypt(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptBytes 把[]byte数组加密成md5字符串
func EncryptBytes(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptBytes 把[]byte数组加密成md5字符串,如果加密失败则panic
func MustEncryptBytes(data []byte) string {
	result, err := EncryptBytes(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptString 把字符串加密成md5
func EncryptString(data string) (encrypt string, err error) {
	return EncryptBytes([]byte(data))
}

// MustEncryptString 把字符串加密成md5，如果失败则panic
func MustEncryptString(data string) string {
	result, err := EncryptString(data)
	if err != nil {
		panic(err)
	}
	return result
}

// EncryptFile 获取指定path文件的md5值
func EncryptFile(path string) (encrypt string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// MustEncryptFile 计算指定文件的md5值，如果计算失败，则panic
func MustEncryptFile(path string) string {
	result, err := EncryptFile(path)
	if err != nil {
		panic(err)
	}
	return result
}
