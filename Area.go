package main

import sf "bitbucket.org/krepa098/gosfml2"

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
				a.tiles[x+y*a.width].SetSprite(0, 9)
				a.tiles[x+y*a.width].Blocks = true
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

func (a *Area) IsBlocked(x, y int) bool {
	return a.tiles[x+y*a.width].Blocks
}
