package znet

import "zinxGame/zinx/ziface"

type serverOption func(server *Server)

func WithMaxPackageSize(maxPackageSize uint32) serverOption {
	return func(server *Server) {
		server.dp.maxPackageSize = maxPackageSize
	}
}

func WithMaxConnection(maxConnection uint32) serverOption {
	return func(server *Server) {
		server.maxConnection = maxConnection
	}
}

func WithServerName(name string) serverOption {
	return func(server *Server) {
		server.name = name
	}
}

func WithIPVersion(ipVersion string) serverOption {
	return func(server *Server) {
		server.ipVersion = ipVersion
	}
}

func WithIP(ip string) serverOption {
	return func(server *Server) {
		server.ip = ip
	}
}

func WithPort(port int) serverOption {
	return func(server *Server) {
		server.port = port
	}
}

type workPoolOption func(wp *WorkPool)
func WithTaskQueueSize(TaskQueueSize uint32) workPoolOption {
	return func(wp *WorkPool) {
		wp.maxTaskQueueSize = TaskQueueSize
	}
}
func WithWorkPoolSize(WorkPoolSize uint32) workPoolOption {
	return func(wp *WorkPool) {
		wp.taskQueues = make([]chan ziface.IRequest, WorkPoolSize)
		wp.maxWorkPoolSize = WorkPoolSize
	}
}
func WithScheduler(scheduler func(request ziface.IRequest) uint32) workPoolOption {
	return func(wp *WorkPool) {
		wp.scheduleHandle = scheduler
	}
}
