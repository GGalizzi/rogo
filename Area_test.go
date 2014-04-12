package main

import "testing"

func TestPlaceTile(t *testing.T) {
	a := PrepareArea()

	a.placeTile("wall", 3, 3)

	if blocks := a.tiles[3+3*a.width].Blocks; !blocks {
		t.Errorf("The tile at 3,3 should be blocking. Blocks: %v", blocks)
	}
}

func TestMoveToBlockedTile(t *testing.T) {
	a := PrepareArea()

	a.placeTile("wall", 3, 3)

	p := NewEntity(0, 0, 2, 3, a)

	p.Move(1, 0)

	if pos := p.Position(); pos.X == 3 && pos.Y == 3 {
		t.Errorf("Position should have remained 2,3. Pos:%v,%v", pos.X, pos.Y)
	}
}

func PrepareArea() *Area {
	a := new(Area)
	a.width = 20
	a.height = 20
	a.tiles = make([]*Tile, a.height*a.width)

	for i, _ := range a.tiles {
		a.tiles[i] = NewTile()
	}

	return a
}
