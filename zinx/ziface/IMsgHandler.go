package ziface

type IMsgHandler interface {
	// 指定对应的Router消息处理方法
	Serve(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// 获取WorkPool
	GetWorkPool() IWorkPool
}
