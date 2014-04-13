package main

import sf "bitbucket.org/krepa098/gosfml2"

//Tile struct represents data about a tile on a map, its sprite and any special configuration it should have.
type Tile struct {
	sprite *sf.Sprite

	Blocks bool
}

//NewTile initializes a Tile with basic (floor) data.
func NewTile() *Tile {
	t := new(Tile)

	t.sprite, _ = sf.NewSprite(EntitiesTexture)
	t.sprite.SetTextureRect(sf.IntRect{ReadSettings().SpriteSize * 2, ReadSettings().SpriteSize * 2, ReadSettings().SpriteSize, ReadSettings().SpriteSize})

	return t
}

//Sprite to implement Spriter interface, returns a Sprite.
func (t *Tile) Sprite() *sf.Sprite {
	return t.sprite
}

//Draw draws the tiles sprite in the given position in the given window.
func (t *Tile) Draw(window *sf.RenderWindow, x, y int) {
	t.sprite.SetPosition(sf.Vector2f{float32(x * ReadSettings().SpriteSize), float32(y * ReadSettings().SpriteSize)})
	t.sprite.Draw(window, sf.DefaultRenderStates())
}
