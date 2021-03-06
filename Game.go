package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"

	sf "bitbucket.org/krepa098/gosfml2"
)

//State represents the state of the game (i.e: playing, in inventory, dead, in the menu, etc)
type State int

const (
	//PLAY state means the player is in control of its character.
	PLAY State = iota
	//QUIT state means the game loop should end.
	QUIT
	//LOOK state means the player is using the look command.
	LOOK
	//LOG state means the player is looking at the log
	LOG
	//INVENTORY state means the player is looking at an inventory list.
	INVENTORY
)

//Drawer is implemented on types that can be drawn on the window.
type Drawer interface {
	Draw(*sf.RenderWindow)
}

//Draw draws any drawer on the window.
func (g *Game) Draw(d Drawer) {
	d.Draw(g.window)
}

type Announcement struct {
	rect *sf.RectangleShape
	text *sf.Text
}

func NewAnnouncement(s string) *Announcement {
	if len(s) == 0 {
		return nil
	}
	a := new(Announcement)
	a.rect, _ = sf.NewRectangleShape()
	a.rect.SetFillColor(sf.ColorBlack())
	a.rect.SetOutlineThickness(1.0)
	a.rect.SetOutlineColor(sf.ColorWhite())

	a.text, _ = sf.NewText(Font)
	a.text.SetCharacterSize(12)
	a.text.SetColor(sf.ColorWhite())
	a.text.SetPosition(sf.Vector2f{1, 1})

	a.text.SetString(s)
	a.text.SetColor(sf.ColorWhite())

	return a
}

func (a *Announcement) Draw(w *sf.RenderWindow) {
	rect := a.text.GetLocalBounds()
	padding := float32(10)
	resW := readSettings().resW
	resH := readSettings().resH
	a.rect.SetSize(sf.Vector2f{rect.Width + padding, rect.Height + padding})
	a.rect.SetPosition(sf.Vector2f{(resW/2 - rect.Width) + padding, (resH/2 - rect.Height) + padding})
	rectBounds := a.rect.GetGlobalBounds()
	a.text.SetPosition(sf.Vector2f{rectBounds.Left + padding/2, rectBounds.Top + padding/2})
	a.rect.Draw(w, sf.DefaultRenderStates())

	a.text.Draw(w, sf.DefaultRenderStates())
}

//Game contains the base data of the game, from the window, to its current entities and area currently in memory.
type Game struct {
	window *sf.RenderWindow
	area   *Area
	player *Mob
	cursor *Mob

	state State
	Settings

	gameView   *sf.View
	statusView *sf.View
	logView    *sf.View

	hpText   *sf.Text
	lookText *sf.Text
	logText  *sf.Text

	announcement *Announcement
}

//NewGame initializes a Game struct.
func NewGame() *Game {
	g := new(Game)
	g.Settings = readSettings()
	g.window = sf.NewRenderWindow(sf.VideoMode{uint(g.resW), uint(g.resH), 32}, "GoSFMLike", sf.StyleDefault, sf.DefaultContextSettings())
	g.state = PLAY

	g.area = NewArea()
	g.player = NewMob("player", 0, 0, 3, 4)
	g.cursor = NewMob("cursor", 0, 0, 2, 2)

	for i := 0; i < 3; i++ {
		g.area.mobs = append(g.area.mobs, NewMobFromFile("orc", 3+i, 1))
		g.area.mobs = append(g.area.mobs, NewMobFromFile("ogre", 4+i, 2))
		g.area.items = append(g.area.items, NewItemFromFile("potion", 4, 4))
	}
	g.area.mobs = append(g.area.mobs, g.player)

	g.gameView = sf.NewView()
	g.gameView.SetCenter(g.player.PosVector())
	g.gameView.SetSize(sf.Vector2f{g.resW * 0.75, g.resH * 0.75})
	g.gameView.SetViewport(sf.FloatRect{0, 0, .75, .75})

	g.statusView = sf.NewView()
	g.statusView.SetSize(sf.Vector2f{g.resW * 0.25, g.resH})
	g.statusView.SetCenter(sf.Vector2f{(g.resW * 0.25) / 2, g.resH / 2})
	g.statusView.SetViewport(sf.FloatRect{.77, 0, .25, 1})

	g.hpText, _ = sf.NewText(Font)
	g.hpText.SetCharacterSize(12)

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

	for g.window.IsOpen() && g.state != QUIT {
		wait := true
	pollLoop:
		//Event / input queue polling
		for event := g.window.PollEvent(); event != nil; event = g.window.PollEvent() {
			switch et := event.(type) {
			case sf.EventClosed:
				g.window.Close()
			case sf.EventTextEntered:
				wait = g.handleInput(et.Char)
				break pollLoop
			}
		}
		g.window.Clear(sf.ColorBlack())

		// Draw status stuff.
		g.window.SetView(g.statusView)
		g.hpText.SetString("HP: " + strconv.Itoa(g.player.curhp) + "/" + strconv.Itoa(g.player.maxhp))
		g.hpText.Draw(g.window, sf.DefaultRenderStates())

		g.drawLog()

		if g.state != INVENTORY {
			g.window.SetView(g.gameView)

			//Draw area (tiles)
			g.Draw(g.area)

			//Draw items
			for _, i := range g.area.items {
				g.Draw(i)
			}

			//Process mobs Ai, check for deaths and draw them.
			for i, m := range g.area.mobs {
				if m.curhp <= 0 {
					m.die(g.area, i)
				}
				//Check if we should process AI
				if !wait && m != g.player {
					g.processAI(m)
				}
				g.Draw(m)
			}

			//Check if player died.
			if g.player.curhp <= 0 {
				fmt.Print("Game Over, you died.\n")
				g.announcement = NewAnnouncement("Your soul has perished.\nPress Escape to exit.")
				g.state = QUIT
			}

			//Display announcement if there is one.
			if g.announcement != nil {
				g.window.SetView(g.window.GetDefaultView())
				g.announcement.Draw(g.window)
				g.window.Display()

				/*
					for !sf.KeyboardIsKeyPressed(sf.KeyEscape) {
						//Pause the game.
					}
				*/

				for !sf.KeyboardIsKeyPressed(sf.KeyEscape) {
					//While the escape key has not been pressed yet
					for event := g.window.PollEvent(); event != nil; event = g.window.PollEvent() {
						//Ignore and pop all events of the queue
					}
				}

				g.announcement = nil
				continue
			}

			if g.state == LOOK {
				g.Draw(g.cursor)
				g.lookText.Draw(g.window, sf.DefaultRenderStates())
			}

			g.window.Display()
		}
	}

}

func (g *Game) drawLog() {
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
}

func (g *Game) openLog() {
	g.logView.SetSize(sf.Vector2f{g.resW, g.resH * 0.85})
	g.logView.SetViewport(sf.FloatRect{.1, .05, 1, .85})
	g.logView.SetCenter(sf.Vector2f{g.resW / 2, (g.resH * .85) / 2})
}

func (g *Game) tryPickUp() {
	for l, i := range g.area.items {
		if g.player.Position() == i.Position() {
			g.player.pickUp(i)
			//g.area.items = removeFromList(g.area.items, l).([]*Item)
			s := g.area.items
			s[l], s = s[len(s)-1], s[:len(s)-1]
			g.area.items = s
			return
		}
	}
}

//listUsables lists all the items that have an effect, and prompts the user to use one.
func (g *Game) listUsables() {
	letter := 'a'
	listText, _ := sf.NewText(Font)
	listText.SetCharacterSize(12)
	listText.SetPosition(sf.Vector2f{12, 12})
	usables := make(map[rune]*Item)
	for k, i := range g.player.inventory {
		if i.effect != nil {
			appendString(listText, strconv.QuoteRune(letter)+" - "+k+" x"+strconv.Itoa(i.stack))
			usables[letter] = i
			letter++
		}
	}

listLoop:
	for g.window.IsOpen() {
		for event := g.window.PollEvent(); event != nil; event = g.window.PollEvent() {
			switch et := event.(type) {
			case sf.EventTextEntered:
				done, used := g.inventoryInput(et.Char, usables)
				if used != nil {
					if used.stack > 1 {
						used.stack--
						break listLoop
					}
					delete(g.player.inventory, used.name)
					break listLoop
				}
				if done {
					break listLoop
				}
			}
		}
		g.window.Clear(sf.ColorBlack())

		g.window.SetView(g.logView)
		g.drawLog()
		g.window.SetView(g.gameView)
		listText.Draw(g.window, sf.DefaultRenderStates())
		g.window.Display()
	}

	g.state = PLAY
}

func (g *Game) processAI(m *Mob) {
	m.moveTowards(g.player.Entity, g)
}

func (g *Game) describe(e *Entity) {
	appendString(g.lookText, e.name)
	g.lookText.SetPosition(e.PosVector())
}

func (g *Game) handleInput(key rune) (wait bool) {
	wait = true
	var inControl *Mob
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

	// The ESC key should be usable in all states to exit back toPLAY
	if key == 27 {
		wait = true
		g.state = PLAY
		g.gameView.SetCenter(g.player.PosVector())
		return
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
	//Get/Grab/Pickup Item
	case 'g':
		wait = false
		g.tryPickUp()
	//Use items menu
	case 'u':
		wait = true
		g.state = INVENTORY
		g.listUsables()
	//Look command
	case 'x':
		wait = true
		g.state = LOOK
		g.cursor.SetPosition(g.player.Position())
	//Show message log
	case 'L':
		wait = true
		g.state = LOG
		g.openLog()

		//Open Souls menu
	case 'S':
		wait = true
		g.announcement = NewAnnouncement(g.player.pouch.getList())

	case 'D':
		//debug
		wait = true
		g.announcement = NewAnnouncement("Debugging")
	//Go down stairs
	case '>':
		if stair := g.area.tiles[g.player.x+g.player.y*g.area.width]; stair.downStair {
			g.switchArea(stair)
		} else {
			log("No stairs there.")
		}

	//Go up stairs
	case '<':
		if stair := g.area.tiles[g.player.x+g.player.y*g.area.width]; stair.upStair {
			g.switchArea(stair)
		} else {
			log("No stairs there.")
		}

	case 'Q':
		g.window.Close()
	default:
		fmt.Println("Can't recognize command: ", key, "\n")
	}

	return
}

//Logic for changing areas when walking up/down stairs.
//Needs a stair arguments to check if the stair has already been used.
func (g *Game) switchArea(stair *Tile) {
	if stair.linkedArea == nil {
		prevArea := g.area
		g.area = NewArea()
		g.area.mobs = append(g.area.mobs, g.player)
		stair.linkedArea = g.area
		if stair.downStair {
			g.area.genTestRoom()
			rx := rand.Intn(14) + 1
			ry := rand.Intn(13) + 1
			returnStair := g.area.placeTile("upStair", rx, ry)
			g.player.SetPosition(sf.Vector2i{rx, ry})
			returnStair.linkedArea = prevArea
			returnStair.linkedStair = stair
			stair.linkedStair = returnStair
		}
	} else {
		g.area = stair.linkedArea
		posf := stair.linkedStair.GetPosition()
		px := int(posf.X) / readSettings().SpriteSize
		py := int(posf.Y) / readSettings().SpriteSize
		g.player.SetPosition(sf.Vector2i{px, py})
	}

}

func (g *Game) inventoryInput(key rune, items map[rune]*Item) (done bool, used *Item) {
	fmt.Printf("Pressed: %q. Corresponds to: %+v", key, items[key])
	if key == 27 {
		return true, nil
	}

	if items[key] != nil && items[key].effect != nil {
		g.player.use(items[key])
		done, used = true, items[key]
		return
	}

	log("Can't use that.")
	done, used = false, nil
	return
}
