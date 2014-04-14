package main

import (
	"fmt"

	sf "bitbucket.org/krepa098/gosfml2"
)

//State represents the state of the game (i.e: playing, in inventory, dead, in the menu, etc)
type State int

const (
	//PLAY state means the player is in control of its character.
	PLAY State = iota
	//LOOK state means the player is using the look command.
	LOOK
)

//Drawer is implemented on types that can be drawn on the window.
type Drawer interface {
	Draw(*sf.RenderWindow)
}

//Draw draws any drawer on the window.
func (g *Game) Draw(d Drawer) {
	d.Draw(g.window)
}

//Game contains the base data of the game, from the window, to its current entities and area currently in memory.
type Game struct {
	window *sf.RenderWindow
	area   *Area
	player *Entity
	cursor *Entity

	entities []*Entity

	state State
	Settings

	gameView *sf.View
	lookText *sf.Text
}

//NewGame initializes a Game struct.
func NewGame() *Game {
	g := new(Game)
	g.Settings = readSettings()
	g.window = sf.NewRenderWindow(sf.VideoMode{g.resW, g.resH, 32}, "GoSFMLike", sf.StyleDefault, sf.DefaultContextSettings())
	g.state = PLAY

	g.area = NewArea()
	g.player = NewEntity("player", 0, 0, 12, 2, g.area)
	g.cursor = NewEntity("cursor", 0, 0, 2, 2, g.area)

	for i := 0; i < 3; i++ {
		g.entities = append(g.entities, NewEntityFromFile("orc", 3+i, 1, g.area))
	}
	g.entities = append(g.entities, g.player)

	g.gameView = sf.NewView()
	g.gameView.SetCenter(g.player.PosVector())
	g.gameView.SetSize(sf.Vector2f{150, 150})
	g.gameView.SetViewport(sf.FloatRect{0, 0, .75, .75})

	var err error
	g.lookText, err = sf.NewText(Font)
	if err != nil {
		panic(err)
	}
	g.lookText.SetCharacterSize(12)

	return g
}

func (g *Game) run() {

	for g.window.IsOpen() {

		wait := true
		for event := g.window.PollEvent(); event != nil; event = g.window.PollEvent() {
			switch et := event.(type) {
			case sf.EventClosed:
				g.window.Close()
			case sf.EventTextEntered:
				wait = g.handleInput(et.Char)
			}
		}
		g.window.Clear(sf.ColorBlack())

		g.window.SetView(g.gameView)
		for _, d := range g.entities {
			if !wait && d != g.player && d.Mob != nil {
				g.processAI(d)
				if g.player.Mob == nil {
					fmt.Print("Game Over, you died.\n")
					g.window.Close()
					return
				}
			}
			g.Draw(d)
		}
		if g.state == LOOK {
			g.Draw(g.cursor)
			g.lookText.Draw(g.window, sf.DefaultRenderStates())
		}

		g.Draw(g.area)
		g.window.Display()
	}

}

func (g *Game) processAI(e *Entity) {
	e.moveTowards(g.player, g)
}

func (g *Game) describe(e *Entity) {
	appendString(g.lookText, e.name)
	g.lookText.SetPosition(e.PosVector())
}

func (g *Game) handleInput(key rune) (wait bool) {
	wait = true
	var inControl *Entity
	if g.state == PLAY {
		inControl = g.player
	} else if g.state == LOOK {
		inControl = g.cursor
	}

	move := func(x, y int) {
		if g.state == LOOK {
			g.lookText.SetString("")
		}
		inControl.Move(x, y, g)
		if g.state == PLAY {
			wait = false
		}
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
	case '5':
		wait = false
	case 'x':
		wait = false
		g.state = LOOK
		g.cursor.SetPosition(g.player.Position())
	case 27: //ESC key
		wait = false
		g.state = PLAY
		g.gameView.SetCenter(g.player.PosVector())
	case 'Q':
		g.window.Close()
	}

	return
}
