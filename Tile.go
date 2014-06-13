package main

import sf "bitbucket.org/krepa098/gosfml2"

//Tile struct represents data about a tile on a map, its sprite and any special configuration it should have.
type Tile struct {
	*Graph

	//Generic
	blocks bool

	//Door stuff
	locked bool
	door   bool

	//Stairs
	downStair   bool
	upStair     bool
	linkedArea  *Area
	linkedStair *Tile
}

//NewTile initializes a Tile with basic (floor) data.
func NewTile() *Tile {
	t := new(Tile)

	t.Graph = NewGraph(1, 9)

	return t
}

//Draw draws the tiles sprite in the given position in the given window.
func (t *Tile) Draw(window *sf.RenderWindow, x, y int) {
	t.setPosition(x, y)
	t.Sprite.Draw(window, sf.DefaultRenderStates())
}
