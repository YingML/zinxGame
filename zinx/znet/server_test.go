package znet

import (
	"fmt"
	"testing"
	"zinxGame/game/utils"
	"zinxGame/zinx/ziface"
)

type PingRouter struct {
	BaseRouter
}

func (r *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle ...")
	request.GetConn().SendMsg(uint32(200), []byte("before ping...\n"))
}

func (r *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ...")
	request.GetConn().SendMsg(uint32(200), []byte("ping...ping...ping...\n"))
}

func (r *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle ...")
	request.GetConn().SendMsg(uint32(200), []byte("after ping...\n"))
}

type HelloRouter struct {
	BaseRouter
}

func (r *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Hello Router Handle ...")
	request.GetConn().SendMsg(uint32(200), []byte("Hello...Hello...Hello...\n"))
}

func TestServer(t *testing.T) {
	// 1、创建1个Server句柄
	rs := make(map[uint32]ziface.IRouter)
	rs[0] = &PingRouter{}
	scheduler := func(request ziface.IRequest) uint32 {
		return request.GetReqID() % utils.Config.MaxWorkPoolSize
	}
	wp := NewWorkPool(WithScheduler(scheduler))
	h := NewMsgHandler(WithRouters(rs), WithWorkPool(wp))
	h.AddRouter(1, &PingRouter{})
	h.AddRouter(2, &HelloRouter{})
	server := NewServer(h)

	server.SetOnConnStart(func(conn ziface.IConnection) {
		fmt.Println("==> onConnStart is Called.")
		conn.SendMsg(202, []byte("DoConnection BEGIN"))
	})
	server.SetOnConnStop(func(conn ziface.IConnection) {
		fmt.Println("==> OnConnStop is Called.")
		fmt.Println("connID:", conn.ConnID(),"is lost.")
	})
	// 2、启动Server
	server.Serve()
}
