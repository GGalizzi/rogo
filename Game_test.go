package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"os"
	"testing"
)

func TestHandleInput(t *testing.T) {
	var g *Game
	g = MockNewGame()

	pos := g.player.Position
	// ey, ex = Expected; cy,cx = Current/Actual coords, ox,oy = Original Coords before moving.
	for r := '1'; r <= '9'; r++ {
		x, y := pos().X, pos().Y
		g.handleInput(r)

		assertMove := func(ex, ey int) {
			if cx, cy := pos().X, pos().Y; cx != ex || cy != ey {
				t.Errorf("Position.Y = %v, want: %v", cy, ey)
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

func MockNewGame() *Game {
	g := new(Game)
	g.area = NewArea()
	g.player = NewEntity(0, 0, 12, 12, g.area)
	g.gameView = sf.NewView()
	return g
}

type MockWindow struct{}

func (w *MockWindow) Close() {
	os.Exit(0)
}
