package dcompress

import (
	"bytes"
	"compress/gzip"
	"io"
)

//使用gzip压缩
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

//解压缩gzip
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
