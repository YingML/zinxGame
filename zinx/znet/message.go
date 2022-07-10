package znet

import "zinxGame/zinx/ziface"

type Message struct {
	id uint32 //消息的ID
	dataLen uint32 //消息的长度
	data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.id
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetMsgID(id uint32) {
	m.id = id
}

func (m *Message) SetDataLen(l uint32) {
	m.dataLen = l
}

func (m *Message) SetData(data []byte) {
	m.data = data
}

var _ ziface.IMessage = &Message{}