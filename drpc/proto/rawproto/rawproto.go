package rawproto

import (
	"donkeygo/drpc/proto"
	"io"
	"sync"
)

var RawProtoFunc = func(rw proto.IOWithReadBuffer) proto.Proto {
	return &rawProto{
		id:   6,
		name: "raw",
		r:    rw,
		w:    rw,
	}
}

var _ proto.Proto = new(rawProto)

type rawProto struct {
	r    io.Reader
	w    io.Writer
	rMu  sync.Mutex
	name string
	id   byte
}

func (that *rawProto) Version() (byte, string) {
	return that.id, that.name
}

func (that *rawProto) Pack(m proto.Message) error {
	return nil
}

func (that *rawProto) Unpack(m proto.Message) error {
	return nil
}
