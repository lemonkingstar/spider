package handler

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type Handler interface {
	Name() string
	Execute() error
}

type messageInfo struct {
	message *protogen.Message
	fields  map[string]*protogen.Field
	oneofs  map[string]*protogen.Oneof
}

func buildFieldIndex(messages []*protogen.Message) map[string]*messageInfo {
	index := make(map[string]*messageInfo)
	walkMessages(messages, func(msg *protogen.Message) {
		info := &messageInfo{
			message: msg,
			fields:  make(map[string]*protogen.Field, len(msg.Fields)),
			oneofs:  make(map[string]*protogen.Oneof, len(msg.Oneofs)),
		}
		for _, f := range msg.Fields {
			info.fields[f.GoName] = f
		}

		for _, o := range msg.Oneofs {
			info.oneofs[o.GoName] = o
			// 兼容 oneof字段类型
			for _, f := range o.Fields {
				index[f.GoIdent.GoName] = &messageInfo{
					message: msg,
					fields: map[string]*protogen.Field{
						f.GoName: f,
					},
				}
			}
		}
		index[msg.GoIdent.GoName] = info
	})
	return index
}

func walkMessages(messages []*protogen.Message, fn func(*protogen.Message)) {
	for _, msg := range messages {
		if msg.Desc.IsMapEntry() {
			continue
		}
		fn(msg)
		walkMessages(msg.Messages, fn)
	}
}
