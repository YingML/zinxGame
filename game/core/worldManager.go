package core

import "sync"

type WorldManager struct {
	// AOIMap 当前地图AOI模块
	AoiMap *AOIMap
	// 当前全部在线的players集合
	Players map[int32]*Player
	// 保护Players集合的锁
	plock sync.RWMutex
}

var WorldMgr *WorldManager

func init(){
	WorldMgr = &WorldManager{
		AoiMap:  NewAOIMap(AoiMinX, AoiMinY, AoiMaxX, AoiMaxY, AoiCntX, AoiCntY),
		Players: make(map[int32]*Player),
	}
}

func (w *WorldManager) AddPlayer(player *Player)  {
	w.plock.Lock()
	w.Players[player.Pid] = player
	w.plock.Unlock()

	w.AoiMap.AddToGridByPos(player.Pid, player.X, player.Z)
}

func (w *WorldManager) RemovePlayerByPid(pid int32)  {
	player := w.Players[pid]
	w.AoiMap.RemoveFromGridByPos(pid, player.X, player.Z)

	w.plock.Lock()
	delete(w.Players, pid)
	w.plock.Unlock()
}

func (w *WorldManager) GetPlayerByPid(pid int32) *Player {
	w.plock.RLock()
	defer w.plock.RUnlock()

	return w.Players[pid]
}

func (w *WorldManager) GetAllPlayers() []*Player {
	w.plock.RLock()
	defer w.plock.RUnlock()

	players := make([]*Player, 0)
	for _, player := range w.Players {
		players = append(players, player)
	}
	return players
}