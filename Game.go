package main

import (
	"fmt"
	"io/ioutil"

	sf "bitbucket.org/krepa098/gosfml2"
)

//State represents the state of the game (i.e: playing, in inventory, dead, in the menu, etc)
type State int

const (
	//PLAY state means the player is in control of its character.
	PLAY State = iota
	//LOOK state means the player is using the look command.
	LOOK
	//LOG state means the player is looking at the log
	LOG
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

	mobs  []*Entity
	items []*Entity

	state State
	Settings

	gameView *sf.View
	logView  *sf.View

	lookText *sf.Text
	logText  *sf.Text
}

//NewGame initializes a Game struct.
func NewGame() *Game {
	g := new(Game)
	g.Settings = readSettings()
	g.window = sf.NewRenderWindow(sf.VideoMode{uint(g.resW), uint(g.resH), 32}, "GoSFMLike", sf.StyleDefault, sf.DefaultContextSettings())
	g.state = PLAY

	g.area = NewArea()
	g.player = NewEntity("player", 0, 0, 3, 4, g.area)
	g.cursor = NewEntity("cursor", 0, 0, 2, 2, g.area)

	for i := 0; i < 1; i++ {
		g.mobs = append(g.mobs, NewEntityFromFile("orc", 3+i, 1, g.area))
	}
	g.items = append(g.items, NewEntityFromFile("potion", 4, 4, g.area))
	g.mobs = append(g.mobs, g.player)

	g.gameView = sf.NewView()
	g.gameView.SetCenter(g.player.PosVector())
	g.gameView.SetSize(sf.Vector2f{g.resW * 0.75, g.resH * 0.75})
	g.gameView.SetViewport(sf.FloatRect{0, 0, .75, .75})

	g.logView = sf.NewView()

	var err error
	g.lookText, err = sf.NewText(Font)
	if err != nil {
		panic(err)
	}
	g.lookText.SetCharacterSize(12)

	g.logText, _ = sf.NewText(Font)
	g.logText.SetCharacterSize(12)

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

		//Draw items
		for _, i := range g.items {
			g.Draw(i)
		}

		//Process mobs Ai and draw them.
		for _, d := range g.mobs {
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

		logFile, err := ioutil.ReadFile("log.txt")
		if err != nil {
			fmt.Println("Can't open the log file log.txt: ERR: ", err)
		}

		if g.state != LOG {
			g.logView.SetSize(sf.Vector2f{g.resW * 0.8, g.resH * 0.25})
			g.logView.SetViewport(sf.FloatRect{0.01, .70, .8, .25})
			glb := g.logText.GetGlobalBounds()
			lvCenterX := (g.resW * 0.8) / 2
			lvCenterY := (g.resH * 0.25) / 2
			g.logView.SetCenter(sf.Vector2f{lvCenterX, glb.Height - lvCenterY})
		}

		g.window.SetView(g.logView)
		g.logText.SetString(string(logFile))
		g.logText.Draw(g.window, sf.DefaultRenderStates())

		g.window.Display()
	}

}

func (g *Game) openLog() {
	g.logView.SetSize(sf.Vector2f{g.resW, g.resH * 0.85})
	g.logView.SetViewport(sf.FloatRect{.1, .05, 1, .85})
	g.logView.SetCenter(sf.Vector2f{g.resW / 2, (g.resH * .85) / 2})
}

func (g *Game) tryPickUp() {
	for l, i := range g.items {
		if g.player.Position() == i.Position() {
			g.player.pickUp(i)
			g.items[len(g.items)-1], g.items[l], g.items = nil, g.items[len(g.items)-1], g.items[:len(g.items)-1]
		}
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
		if g.state != LOG {
			inControl.Move(x, y, g)

			g.gameView.SetCenter(inControl.PosVector())
		}
		if g.state == PLAY {
			wait = false
		}
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
	case 'g':
		g.tryPickUp()
	case 'x':
		wait = true
		g.state = LOOK
		g.cursor.SetPosition(g.player.Position())
	case 'L':
		wait = true
		g.state = LOG
		g.openLog()
	case 27: //ESC key
		wait = true
		g.state = PLAY
		g.gameView.SetCenter(g.player.PosVector())
	case 'Q':
		g.window.Close()
	default:
		fmt.Println("Can't recognize command: ", key)
	}

	return
}
