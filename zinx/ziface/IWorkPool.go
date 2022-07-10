package ziface

type IWorkPool interface {
	Start(workerHandle func(workID uint32, queue chan IRequest))
	SendToTask(request IRequest)
	GetMaxWorkPoolSize() uint32
}
