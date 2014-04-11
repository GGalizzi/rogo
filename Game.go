package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

type Drawer interface {
	Draw(*sf.RenderWindow)
}

func (g *Game) Draw(d Drawer) {
	d.Draw(g.Window)
}

type Game struct {
	Window *sf.RenderWindow
	Map    *Map
	Player *Entity

	Drawers []Drawer // Drawable entities on map.
}

func NewGame() *Game {
	g := new(Game)
	g.Window = sf.NewRenderWindow(sf.VideoMode{840, 780, 32}, "GoSFMLike", sf.StyleDefault, sf.DefaultContextSettings())

	g.Map = NewMap()
	g.Player = NewEntity(0, 0, 2, 2, g.Map)

	g.Drawers = append(g.Drawers, g.Map, g.Player)

	return g
}

func (g *Game) run() {
	for g.Window.IsOpen() {
		for event := g.Window.PollEvent(); event != nil; event = g.Window.PollEvent() {
			switch et := event.(type) {
			case sf.EventClosed:
				g.Window.Close()
			case sf.EventTextEntered:
				g.handleInput(et.Char)
			}
		}
		g.Window.Clear(sf.ColorBlack())

		for _, d := range g.Drawers {
			g.Draw(d)
		}
		g.Window.Display()
	}
}

func (g *Game) handleInput(key rune) {
	switch key {
	case '2':
		g.Player.Move(0, 1)
	case '3':
		g.Player.Move(1, 1)
	case '6':
		g.Player.Move(1, 0)
	case '9':
		g.Player.Move(1, -1)
	case '8':
		g.Player.Move(0, -1)
	case '7':
		g.Player.Move(-1, -1)
	case '4':
		g.Player.Move(-1, 0)
	case '1':
		g.Player.Move(-1, 1)
	case 'Q':
		g.Window.Close()
	}
}
