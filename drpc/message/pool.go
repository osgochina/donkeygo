package message

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/drpc/tfilter"
	"sync"
)

var messagePool = sync.Pool{
	New: func() interface{} {
		return NewMessage()
	},
}

//GetMessage 从对象池中获取message
func GetMessage(settings ...MsgSetting) Message {
	m := messagePool.Get().(*message)
	m.doSetting(settings...)
	return m
}

//PutMessage 把message对象返回对象池
func PutMessage(m Message) {
	m.Reset()
	messagePool.Put(m)
}

//NewMessage 创建message对象
func NewMessage(settings ...MsgSetting) Message {
	var m = &message{
		meta:        dmap.New(true),
		pipeTFilter: tfilter.NewPipeTFilter(),
	}
	m.doSetting(settings...)
	return m
}
