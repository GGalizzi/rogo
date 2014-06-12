package main

import (
	"math/rand"

	sf "bitbucket.org/krepa098/gosfml2"
)

//Area contains data that relates to an area, a map, a dungeon. Basically, a set of tiles.
type Area struct {
	width  int
	height int

	mobs  []*Entity
	items []*Entity

	tiles []*Tile
}

//NewArea initializes an Area struct with a basic map.
func NewArea() *Area {
	a := &Area{width: 220, height: 120}

	a.tiles = make([]*Tile, a.height*a.width)

	for i := range a.tiles {
		a.tiles[i] = NewTile()
	}

	seed := rand.Uint32()
	a.genFromPerlin(uint(seed))

	return a
}

func (a *Area) genTestRoom() {
	for x := 0; x < a.width; x++ {
		for y := 0; y < a.height; y++ {
			if y == 0 || y == a.height-1 || x == 0 || x == a.width-1 || (x == a.width/2) {
				a.placeTile("wall", x, y)
			}
			if y == a.height/2 && x == a.width/2 {
				a.placeTile("lockedDoor", x, y)
			}
		}
	}
}

func (a *Area) genFromPerlin(seed uint) {
	pn := NewPerlin(seed)

	for Y := 0.0; Y < float64(a.height); Y++ {
		for X := 0.0; X < float64(a.width); X++ {
			x := X / float64(a.width)
			y := Y / float64(a.height)

			n := pn.noise(float64(10*x), float64(10*y), float64(0.8))

			xx := int(X)
			yy := int(Y)
			if n < 0.33 || xx == 0 || xx == a.width-1 || yy == 0 || yy == a.height-1 {
				a.placeTile("wall", xx, yy)
				continue
			}
			if n > 0.8 {
				a.placeTile("water", xx, yy)
			}
		}
	}
}

//Draw draws all the tiles that make the area.
func (a *Area) Draw(window *sf.RenderWindow) {
	var fromX int
	var toX int
	var fromY int
	var toY int

	sight := 8

	player := a.mobs[0]
	if player.x-sight < 0 {
		fromX = 0
	} else {
		fromX = player.x - sight
	}

	if player.x+sight > a.width {
		toX = a.width
	} else {
		toX = player.x + sight
	}

	if player.y-sight < 0 {
		fromY = 0
	} else {
		fromY = player.y - sight
	}

	if player.y+sight > a.height {
		toY = a.height
	} else {
		toY = player.y + sight
	}

	for x := fromX; x < toX; x++ {
		for y := fromY; y < toY; y++ {
			a.tiles[x+y*a.width].Draw(window, x, y)
		}
	}
}

//placeTile places the given tile by reading from the JSON that contains its data.
func (a *Area) placeTile(name string, x, y int) {
	data := ReadJSON("tiles", name)

	t := a.tiles[x+y*a.width]
	t.blocks = data["blocks"].(bool)
	if data["color"] != nil {
		r := data["color"].([]interface{})[0].(float64)
		g := data["color"].([]interface{})[1].(float64)
		b := data["color"].([]interface{})[2].(float64)
		t.SetColor(sf.Color{byte(r), byte(g), byte(b), 255})
	}

	if locked := data["locked"]; locked != nil {
		t.locked = locked.(bool)
		t.door = true

		//Call some function to determine the location of the key for the door.
		key := NewEntityFromFile("key", x-5, y)
		key.linkedDoor = t

		a.items = append(a.items, key)
	}

	t.setSprite(int(data["spriteX"].(float64)), int(data["spriteY"].(float64)))
}

//IsBlocked checks if the tile in the given coords blocks movement.
func (a *Area) IsBlocked(x, y int) bool {
	return a.tiles[x+y*a.width].blocks
}

func (a *Area) isDoor(x, y int) bool {
	return a.tiles[x+y*a.width].door
}
