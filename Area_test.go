package main

import "testing"

func TestPlaceTile(t *testing.T) {
	a := new(Area)
	a.width = 20
	a.height = 20
	a.tiles = make([]*Tile, a.height*a.width)

	for i, _ := range a.tiles {
		a.tiles[i] = NewTile()
	}

	a.placeTile("wall", 3, 3)

	if blocks := a.tiles[3+3*a.width].Blocks; !blocks {
		t.Errorf("The tile at 3,3 should be blocking. Blocks: %v", blocks)
	}
}
