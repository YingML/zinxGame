package znet

import (
	"fmt"
	"net"
	"zinxGame/zinx/ziface"
)

type Server struct {
	name string
	// 服务器绑定的IP版本
	ipVersion string
	ip        string
	port      int
	// 该server的消息路由
	handler ziface.IMsgHandler
	// 该server的连接管理
	cmr           ziface.IConnManager
	onConnStart   func(conn ziface.IConnection)
	onConnStop    func(conn ziface.IConnection)
	dp            *DataPack
	maxConnection uint32
}

func (s *Server) GetDataPack() ziface.IDataPack {
	return s.dp
}

func (s *Server) SetOnConnStart(f func(conn ziface.IConnection)) {
	s.onConnStart = f
}

func (s *Server) SetOnConnStop(f func(conn ziface.IConnection)) {
	s.onConnStop = f
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.onConnStart != nil {
		fmt.Println("--> Call onConnStart() ...")
		s.onConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.onConnStop != nil {
		fmt.Println("--> Call onConnStop() ...")
		s.onConnStop(conn)
	}
}

func NewServer(h ziface.IMsgHandler, options ...serverOption) *Server {
	s := &Server{
		cmr:           NewConnManager(),
		name:          "Zinx v1.0",
		ipVersion:     "tcp4",
		ip:            "0.0.0.0",
		port:          8999,
		handler:       h,
		dp:            NewDataPack(),
		maxConnection: 1024,
	}

	for _, fn := range options {
		fn(s)
	}

	return s
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.cmr
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener as ip %s, port %d, is starting\n", s.ip, s.port)
	wp := s.handler.GetWorkPool()
	if wp != nil {
		wp.Start(func(workID uint32, reqChan chan ziface.IRequest) {
			fmt.Printf("[Start Worker, ID: %d]\n", workID)
			for {
				select {
					case r := <- reqChan:
						s.handler.Serve(r)
				}
			}
		})
	}
	go func() {
		// 1、获取1个TCP的addr
		addr, err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			fmt.Println("resolve tcp addr error : ", err)
			return
		}
		// 2、监听服务器的地址
		listener, err := net.ListenTCP(s.ipVersion, addr)
		if err != nil {
			fmt.Println("listen", addr, "err", err)
			return
		}
		fmt.Println("Start Zinx server success,", s.name, "Listening...")

		var connID uint32 = 0
		// 3、 阻塞的等待客户端链接，处理客户端连接业务(读写)
		for {
			// 如果有客户端链接进来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}
			// 如果超过最大连接数，则拒绝此新的连接
			if uint32(s.cmr.Len()+1) > s.maxConnection {
				// TODO 给客户端一个超出最大连接的错误包
				fmt.Println("Too Many Client Connection To Server.")
				conn.Close()
				continue
			}

			c := NewConn(s, conn, connID, s.handler)

			connID++
			//已经与客户端建立连接，处理与客户端业务
			go c.Start()
		}
	}()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {

	}
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态会一些一开辟的链接信息，进行关闭或回收
	fmt.Println("[STOP] Zinx server name ", s.name)
	s.cmr.RemoveAll()
}

var _ ziface.IServer = &Server{}