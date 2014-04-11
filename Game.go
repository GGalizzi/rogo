package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

type Game struct {
	Window *sf.RenderWindow
	Map    *Map
	Player *Entity

	Texture *sf.Texture
}

func NewGame() *Game {
	g := new(Game)
	g.Window = sf.NewRenderWindow(sf.VideoMode{840, 780, 32}, "GoSFMLike", sf.StyleDefault, sf.DefaultContextSettings())

	g.Map = NewMap()
	g.Player = NewEntity(0, 0, 2, 2, g.Map)

	return g
}

func (g *Game) run() {
	for g.Window.IsOpen() {
		for event := g.Window.PollEvent(); event != nil; event = g.Window.PollEvent() {
			switch et := event.(type) {
			case sf.EventClosed:
				g.Window.Close()
			case sf.EventKeyPressed:
				g.handleInput(et.Code)
			}
		}
		g.Window.Clear(sf.ColorBlack())
		g.Map.Draw(g.Window)
		g.Player.Sprite.Draw(g.Window, sf.DefaultRenderStates())
		g.Window.Display()
	}
}

func (g *Game) handleInput(key sf.KeyCode) {
	switch key {
	case sf.KeyNumpad2:
		g.Player.Move(0, 1)
	case sf.KeyNumpad3:
		g.Player.Move(1, 1)
	case sf.KeyNumpad6:
		g.Player.Move(1, 0)
	case sf.KeyNumpad9:
		g.Player.Move(1, -1)
	case sf.KeyNumpad8:
		g.Player.Move(0, -1)
	case sf.KeyNumpad7:
		g.Player.Move(-1, -1)
	case sf.KeyNumpad4:
		g.Player.Move(-1, 0)
	case sf.KeyNumpad1:
		g.Player.Move(-1, 1)
	}
}
