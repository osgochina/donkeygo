package dcompress

import (
	"bytes"
	"compress/gzip"
	"github.com/osgochina/donkeygo/os/dfile"
	"io"
)

// Gzip 使用gzip压缩
func Gzip(data []byte, level ...int) ([]byte, error) {
	var (
		writer *gzip.Writer
		buf    bytes.Buffer
		err    error
	)
	if len(level) > 0 {
		writer, err = gzip.NewWriterLevel(&buf, level[0])
		if err != nil {
			return nil, err
		}
	} else {
		writer = gzip.NewWriter(&buf)
	}
	if _, err = writer.Write(data); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnGzip 解压缩gzip
func UnGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(&buf, reader); err != nil {
		return nil, err
	}
	if err = reader.Close(); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// GzipFile 压缩文件
func GzipFile(src, dst string, level ...int) error {
	var (
		writer *gzip.Writer
		err    error
	)
	srcFile, err := dfile.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := dfile.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if len(level) > 0 {
		writer, err = gzip.NewWriterLevel(dstFile, level[0])
		if err != nil {
			return err
		}
	} else {
		writer = gzip.NewWriter(dstFile)
	}
	defer writer.Close()

	_, err = io.Copy(writer, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// UnGzipFile 解压缩文件
func UnGzipFile(src, dst string) error {
	srcFile, err := dfile.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := dfile.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	reader, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	if _, err = io.Copy(dstFile, reader); err != nil {
		return err
	}
	return nil
}
