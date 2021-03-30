package gzip

import (
	"bytes"
	"compress/gzip"
	"donkeygo/drpc/tfilter"
	"fmt"
	"io/ioutil"
	"sync"
)

// Reg registers a gzip filter for transfer.
func Reg(id byte, name string, level int) {
	tfilter.Reg(newGzip(id, name, level))
}

type Gzip struct {
	id    byte
	name  string
	level int
	wPool sync.Pool
	rPool sync.Pool
}

func newGzip(id byte, name string, level int) *Gzip {
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		panic(fmt.Sprintf("gzip: invalid compression level: %d", level))
	}
	g := new(Gzip)
	g.level = level
	g.id = id
	g.name = name
	g.wPool = sync.Pool{
		New: func() interface{} {
			gw, _ := gzip.NewWriterLevel(nil, g.level)
			return gw
		},
	}
	g.rPool = sync.Pool{
		New: func() interface{} {
			return new(gzip.Reader)
		},
	}
	return g
}

// ID returns transfer filter id.
func (g *Gzip) ID() byte {
	return g.id
}

// ID returns transfer filter name.
func (g *Gzip) Name() string {
	return g.name
}

// OnPack performs filtering on packing.
func (g *Gzip) OnPack(src []byte) ([]byte, error) {
	var bb = new(bytes.Buffer)
	gw := g.wPool.Get().(*gzip.Writer)
	gw.Reset(bb)
	_, _ = gw.Write(src)
	err := gw.Close()
	gw.Reset(nil)
	g.wPool.Put(gw)
	if err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

// OnUnpack performs filtering on unpacking.
func (g *Gzip) OnUnpack(src []byte) (dest []byte, err error) {
	if len(src) == 0 {
		return src, nil
	}
	gr := g.rPool.Get().(*gzip.Reader)
	err = gr.Reset(bytes.NewReader(src))
	if err == nil {
		dest, err = ioutil.ReadAll(gr)
	}
	_ = gr.Close()
	g.rPool.Put(gr)
	return dest, err
}
