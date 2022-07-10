package api

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinxGame/game/core"
	"zinxGame/game/pb"
	"zinxGame/zinx/ziface"
	"zinxGame/zinx/znet"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (w *WorldChatApi) Handle(request ziface.IRequest) {
	talk := &pb.Talk{}
	if err := proto.Unmarshal(request.GetData(), talk); err != nil {
		fmt.Println("proto unmarshal failed:", err)
		return
	}

	pid, err := request.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println("Get Pid failed: ", err)
		return
	}

	player := core.WorldMgr.GetPlayerByPid(pid.(int32))

	player.Talk(talk.Content)
}

var _ ziface.IRouter = &WorldChatApi{}
