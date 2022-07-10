package znet

import "zinxGame/zinx/ziface"

type Request struct {
	id uint32
	// 已经和客户端建立好的连接
	conn ziface.IConnection
	data ziface.IMessage
}

func NewRequest(id uint32,c ziface.IConnection, msg ziface.IMessage) *Request {
	return &Request{
		id: id,
		conn: c,
		data: msg,
	}
}

func (r *Request) GetReqID() uint32 {
	return r.id
}

func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data.GetData()
}

func (r *Request) GetMsgID() uint32  {
	return r.data.GetMsgID()
}

var _ ziface.IRequest = &Request{}
