package znet

import (
	"fmt"
	"zinxGame/zinx/ziface"
)

const (
	DefaultWorkPoolSize uint32 = 10
	DefaultTaskQueueSize uint32 = 1024
)

type WorkPool struct {
	taskQueues       []chan ziface.IRequest
	maxTaskQueueSize uint32
	maxWorkPoolSize  uint32
	scheduleHandle   func(request ziface.IRequest) uint32
}

func (w *WorkPool) GetMaxWorkPoolSize() uint32 {
	return w.maxWorkPoolSize
}
func NewWorkPool(options ...workPoolOption) *WorkPool {
	wp := &WorkPool{
		taskQueues:       make([]chan ziface.IRequest, DefaultWorkPoolSize),
		maxTaskQueueSize: DefaultTaskQueueSize,
		maxWorkPoolSize:  DefaultWorkPoolSize,
	}
	for _, fn := range options {
		fn(wp)
	}
	return wp
}

func (w *WorkPool) Start(workerHandle func(workID uint32, queue chan ziface.IRequest)) {
	for i := uint32(0); i < w.maxWorkPoolSize; i++ {
		w.taskQueues[i] = make(chan ziface.IRequest, w.maxTaskQueueSize)
		go workerHandle(i, w.taskQueues[i])
	}
}

func (w *WorkPool) SendToTask(request ziface.IRequest) {
	var workID uint32
	if w.scheduleHandle != nil{
		workID = w.scheduleHandle(request)
	}else {
		workID = request.GetReqID() % w.maxWorkPoolSize
	}
	w.taskQueues[workID] <- request
	fmt.Printf("Send Request[%d] To worker[%d]\n", request.GetReqID(), workID)
}

var _ ziface.IWorkPool = &WorkPool{}
