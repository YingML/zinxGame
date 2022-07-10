package znet

import "zinxGame/zinx/ziface"

type BaseRouter struct {
	
}
// 之所以方法都为空
// 是因为有的Router不希望有PreHandle、PostHandle这两个业务
// 所以Router全部集成Base的好处就是， 不需要实现PreHandle、PostHandle

func (r *BaseRouter) PreHandle(request ziface.IRequest) {}

func (r *BaseRouter) Handle(request ziface.IRequest) {}

func (r *BaseRouter) PostHandle(request ziface.IRequest) {}

var _ ziface.IRouter = &BaseRouter{}
