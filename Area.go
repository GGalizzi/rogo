package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

//Area contains data that relates to an area, a map, a dungeon. Basically, a set of tiles.
type Area struct {
	width  int
	height int

	tiles []*Tile
}

//NewArea initializes an Area struct with a basic map.
func NewArea() *Area {
	a := &Area{width: 20, height: 20}

	a.tiles = make([]*Tile, a.height*a.width)

	for i := range a.tiles {
		a.tiles[i] = NewTile()
	}

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

	return a
}

//Draw draws all the tiles that make the area.
func (a *Area) Draw(window *sf.RenderWindow) {
	for x := 0; x < a.width; x++ {
		for y := 0; y < a.height; y++ {
			a.tiles[x+y*a.width].Draw(window, x, y)
		}
	}
}

//placeTile places the given tile by reading from the JSON that contains its data.
func (a *Area) placeTile(name string, x, y int) {
	data := ReadJSON("tiles", name)

	t := a.tiles[x+y*a.width]
	t.blocks = data["blocks"].(bool)

	if locked := data["locked"]; locked != nil {
		t.locked = locked.(bool)
		t.door = true
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
