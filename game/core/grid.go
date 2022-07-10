package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID int
	MinX int
	MaxX int
	MinY int
	MaxY int
	playerIDs map[int32]bool
	pIDLock sync.RWMutex
}

func NewGrid(gID int, minX ,maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int32]bool),
		pIDLock:   sync.RWMutex{},
	}
}

func (g *Grid) Add(playerID int32)  {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

func (g *Grid) Remove(playerID int32)  {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int32)  {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid: {id: %d, minX: %d, minY: %d, maxX: %d, maxY: %d, playerIDs: %v}\n", g.GID, g.MinX, g.MinY, g.MaxX, g.MaxY, g.playerIDs)
}