package main

import sf "bitbucket.org/krepa098/gosfml2"

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
	g.player = NewEntity(0, 0, 2, 2, g.area)
	g.cursor = NewEntity(0, 0, 2, 2, g.area)

	for i := 0; i < 3; i++ {
		g.entities = append(g.entities, NewEntityFromFile("orc", 3, 1, g.area))
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
		for _, d := range g.entities {
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

func (g *Game) describe(e *Entity) {
	appendString(g.lookText, e.name)
	g.lookText.SetPosition(e.PosVector())
}

func (g *Game) handleInput(key rune) {
	var inControl *Entity
	if g.state == PLAY {
		inControl = g.player
	} else if g.state == LOOK {
		inControl = g.cursor
	}

	move := func(x, y int) {
		cp := inControl.Position()
		cp.X += x
		cp.Y += y
		describing := false
		for _, e := range g.entities {
			if ep := e.Position(); ep == cp {
				subject := e
				switch {
				case e.Mob != nil && inControl == g.player:
					inControl.attack(subject)
					return
				case inControl == g.cursor:
					if !describing {
						g.lookText.SetString("")
					}
					g.describe(subject)
					describing = true
					break
				}
			}
		}
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
	case 'x':
		g.state = LOOK
		g.cursor.SetPosition(g.player.Position())
	case 27: //ESC key
		g.state = PLAY
		g.gameView.SetCenter(g.player.PosVector())
	case 'Q':
		g.window.Close()
	}
}
