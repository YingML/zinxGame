package core

import (
	"fmt"
	"testing"
)

func TestNewAOIMap(t *testing.T) {
	am := NewAOIMap(0,0,250,250,5,5)
	fmt.Println(am)
}

func TestAOIMap_GetSurroundsGridsByGid(t *testing.T) {
	am := NewAOIMap(0,0,250,250,5,5)
	//fmt.Printf("%+v", am.GetSurroundsGridsByGid(24))
	fmt.Printf("%+v", am.GetSurroundsGridsByGid(5))
}
