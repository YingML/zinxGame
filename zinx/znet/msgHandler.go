package znet

import (
	"fmt"
	"zinxGame/zinx/ziface"
)

type MsgHandler struct {
	routers map[uint32]ziface.IRouter
	wp      ziface.IWorkPool
}

type msgOptions func (h *MsgHandler)
func WithRouters(routers map[uint32]ziface.IRouter) msgOptions {
	return func(h *MsgHandler) {
		h.routers = routers
	}
}
func WithWorkPool(wp ziface.IWorkPool) msgOptions {
	return func(h *MsgHandler) {
		h.wp = wp
	}
}

func NewMsgHandler(options ...msgOptions) *MsgHandler {
	h := &MsgHandler{
		routers: make(map[uint32]ziface.IRouter),
	}
	for _, o := range options {
		o(h)
	}
	return h
}

func (m *MsgHandler) GetWorkPool() ziface.IWorkPool {
	return m.wp
}

func (m *MsgHandler) Serve(request ziface.IRequest) {
	msgID := request.GetMsgID()
	if h, ok := m.routers[msgID]; ok {
		h.PreHandle(request)
		h.Handle(request)
		h.PostHandle(request)
	}else {
		fmt.Println("Invalid MessageID: ", msgID)
		return
	}
}

func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	m.routers[msgID]=router
}

var _ ziface.IMsgHandler = &MsgHandler{}