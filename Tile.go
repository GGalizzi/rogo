package main

import sf "bitbucket.org/krepa098/gosfml2"

type Tile struct {
	x      int
	y      int
	Sprite *sf.Sprite

	Blocks bool
}

func NewTile() *Tile {
	t := new(Tile)

	t.Sprite, _ = sf.NewSprite(EntitiesTexture)
	t.Sprite.SetTextureRect(sf.IntRect{ReadSettings().SpriteSize * 2, ReadSettings().SpriteSize * 2, ReadSettings().SpriteSize, ReadSettings().SpriteSize})

	return t
}

func (t *Tile) SetSprite(x, y int) {
	t.Sprite.SetTextureRect(sf.IntRect{ReadSettings().SpriteSize * x, ReadSettings().SpriteSize * y, ReadSettings().SpriteSize, ReadSettings().SpriteSize})
}

func (t *Tile) Draw(window *sf.RenderWindow, x, y int) {
	t.Sprite.SetPosition(sf.Vector2f{float32(x * ReadSettings().SpriteSize), float32(y * ReadSettings().SpriteSize)})
	t.Sprite.Draw(window, sf.DefaultRenderStates())
}
