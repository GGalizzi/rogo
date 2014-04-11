package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

type Drawer interface {
	Draw(*sf.RenderWindow)
}

func (g *Game) Draw(d Drawer) {
	d.Draw(g.window)
}

type Game struct {
	window *sf.RenderWindow
	area   *Area
	player *Entity

	drawers []Drawer // Drawable entities on map.

	gameView *sf.View
}

func NewGame() *Game {
	g := new(Game)
	g.window = sf.NewRenderWindow(sf.VideoMode{840, 780, 32}, "GoSFMLike", sf.StyleDefault, sf.DefaultContextSettings())

	g.area = NewArea()
	g.player = NewEntity(0, 0, 2, 2, g.area)

	g.drawers = append(g.drawers, g.area, g.player)

	g.gameView = sf.NewView()
	g.gameView.SetCenter(g.player.PosVector())
	g.gameView.SetSize(sf.Vector2f{150, 150})
	g.gameView.SetViewport(sf.FloatRect{0, 0, .75, .75})

	return g
}

func (g *Game) run() {
	for g.window.IsOpen() {
		for event := g.window.PollEvent(); event != nil; event = g.window.PollEvent() {
			switch et := event.(type) {
			case sf.EventClosed:
				g.window.Close()
			case sf.EventTextEntered:
				g.handleInput(et.Char)
			}
		}
		g.window.Clear(sf.ColorBlack())

		g.window.SetView(g.gameView)
		for _, d := range g.drawers {
			g.Draw(d)
		}
		g.window.Display()
	}
}

func (g *Game) handleInput(key rune) {
	inControl := g.player

	move := func(x, y int) {
		inControl.Move(x, y)
		g.gameView.SetCenter(inControl.PosVector())
	}

	switch key {
	case '2':
		move(0, 1)
	case '3':
		move(1, 1)
	case '6':
		move(1, 0)
	case '9':
		move(1, -1)
	case '8':
		move(0, -1)
	case '7':
		move(-1, -1)
	case '4':
		move(-1, 0)
	case '1':
		move(-1, 1)
	case 'Q':
		g.window.Close()
	}
}
