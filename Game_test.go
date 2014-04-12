package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"testing"
)

func TestHandleInput(t *testing.T) {
	var g *Game
	g = MockNewGame()

	pos := g.player.Position
	vec := g.player.PosVector
	// ey, ex = Expected; cy,cx = Current/Actual coords, ox,oy = Original Coords before moving.
	for r := '1'; r <= '9'; r++ {
		x, y := pos().X, pos().Y
		g.handleInput(r)

		assertMove := func(ex, ey int) {
			ss := ReadSettings().SpriteSize
			if cx, cy, vx, vy := pos().X, pos().Y, vec().X, vec().Y; cx != ex || cy != ey || vx != float32(ex*ss) || vy != float32(ey*ss) {
				t.Errorf("Expected %v,%v;%v,%v; Got: %v,%v;%v,%v", ex, ey, ex*ss, ey*ss, cx, cy, vx, vy)
			}
		}
		switch r {
		case '1':
			assertMove(x-1, y+1)
		case '2':
			assertMove(x, y+1)
		case '3':
			assertMove(x+1, y+1)
		case '4':
			assertMove(x-1, y)
		case '5':
			assertMove(x, y)
		case '6':
			assertMove(x+1, y)
		case '7':
			assertMove(x-1, y-1)
		case '8':
			assertMove(x, y-1)
		case '9':
			assertMove(x+1, y-1)
		}
	}
}

func TestSettings(t *testing.T) {
	if actual := ReadSettings().resW; actual != 840 {
		t.Errorf("resW = %v, expected: %v", actual, 840)
	}
	if actual := ReadSettings().resH; actual != 780 {
		t.Errorf("resH = %v, expected: %v", actual, 780)
	}
}

func MockNewGame() *Game {
	g := new(Game)
	g.area = NewArea()
	g.player = NewEntity(0, 0, 12, 12, g.area)
	g.gameView = sf.NewView()
	return g
}
