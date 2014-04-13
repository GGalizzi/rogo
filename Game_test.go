package main

import (
	"testing"

	sf "bitbucket.org/krepa098/gosfml2"
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
			ss := readSettings().SpriteSize
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
	if actual := readSettings().resW; actual != 840 {
		t.Errorf("resW = %v, expected: %v", actual, 840)
	}
	if actual := readSettings().resH; actual != 780 {
		t.Errorf("resH = %v, expected: %v", actual, 780)
	}
}

func TestStateCommands(t *testing.T) {
	g := MockNewGame()

	g.handleInput('x')
	if state := g.state; state != LOOK {
		t.Errorf("Expected: %v, Got: %v", state, LOOK)
	}
	g.handleInput('2')
	ePPos := sf.Vector2i{12, 12}
	eCPos := sf.Vector2i{12, 13}
	if pPos, cPos := g.player.Position(), g.cursor.Position(); pPos != ePPos || cPos != eCPos {
		t.Errorf("PlayerPos: %v,%v. CursorPos: %v,%v. Player should be 12,12 Cursor should be 12,13", pPos.X, pPos.Y, cPos.X, cPos.Y)
	}

	g.handleInput(27)
	if state := g.state; state != PLAY {
		t.Errorf("Expected: %v, Got: %v", state, PLAY)
	}
}

func MockNewGame() *Game {
	g := new(Game)
	g.area = NewArea()
	g.player = NewEntity(0, 0, 12, 12, g.area)
	g.cursor = NewEntity(0, 0, 12, 12, g.area)
	g.gameView = sf.NewView()
	return g
}
