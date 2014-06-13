package main

import "testing"

func TestPlaceTile(t *testing.T) {
	a := PrepareArea()

	a.placeTile("wall", 3, 3)

	if blocks := a.tiles[3+3*a.width].blocks; !blocks {
		t.Errorf("The tile at 3,3 should be blocking. Blocks: %v", blocks)
	}
}

func TestMoveToBlockedTile(t *testing.T) {
	g := MockNewGame()

	g.area.placeTile("wall", 3, 3)

	p := NewEntity("player", 0, 0, 2, 3)

	p.Move(1, 0, g)

	if pos := p.Position(); pos.X == 3 && pos.Y == 3 {
		t.Errorf("Position should have remained 2,3. Pos:%v,%v", pos.X, pos.Y)
	}
}

func TestDownStairs(t *testing.T) {
	g := MockNewGame()
	prevAreaName := g.area.name
	t.Logf("PrevArea: %s", prevAreaName)

	g.area.placeTile("downStair", 3, 3)
	p := g.player
	p.Place(2, 3)
	p.Move(1, 0, g)

	g.handleInput('>')

	if g.area.name == prevAreaName {
		t.Errorf("Area should have changed: PrevArea: %s, CurArea: %v", prevAreaName, g.area.name)
	}

	t.Logf("CurArea: %s", g.area.name)
}

func PrepareArea() *Area {
	a := new(Area)
	a.name = "Testing Area"
	a.width = 20
	a.height = 20
	a.tiles = make([]*Tile, a.height*a.width)

	for i := range a.tiles {
		a.tiles[i] = NewTile()
	}

	return a
}
