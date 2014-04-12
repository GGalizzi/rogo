package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"encoding/json"
	"os"
)

type Area struct {
	tiles  []*Tile
	width  int
	height int
}

func NewArea() *Area {
	a := new(Area)
	a.width = 50
	a.height = 20

	a.tiles = make([]*Tile, a.height*a.width)

	for i, _ := range a.tiles {
		a.tiles[i] = NewTile()
	}

	for x := 0; x < a.width; x++ {
		for y := 0; y < a.height; y++ {
			if y == 0 || y == a.height-1 || x == 0 || x == a.width-1 {
				a.placeTile("wall", x, y)
			} else {
				a.tiles[x+y*a.width].Blocks = false
			}
		}
	}

	return a
}

func (a *Area) Draw(window *sf.RenderWindow) {
	for x := 0; x < a.width; x++ {
		for y := 0; y < a.height; y++ {
			a.tiles[x+y*a.width].Draw(window, x, y)
		}
	}
}

func (a *Area) placeTile(name string, x, y int) {
	file, err := os.Open("tiles/" + name + ".tile")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	jParser := json.NewDecoder(file)

	var t interface{}

	if err = jParser.Decode(&t); err != nil {
		panic(err)
	}

	data := t.(map[string]interface{})
	a.tiles[x+y*a.width].Blocks = data["blocks"].(bool)
	SetSprite(a.tiles[x+y*a.width], int(data["spriteX"].(float64)), int(data["spriteY"].(float64)))
}

func (a *Area) IsBlocked(x, y int) bool {
	return a.tiles[x+y*a.width].Blocks
}
