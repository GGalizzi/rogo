package main

import sf "bitbucket.org/krepa098/gosfml2"

var (
	EntitiesTexture, _ = sf.NewTextureFromFile("ascii.png", nil)
)

type Entity struct {
	Sprite *sf.Sprite
	X      int
	Y      int

	Map *Map
}

func NewEntity(spriteX, spriteY, posX, posY int, ma *Map) *Entity {
	e := new(Entity)

	e.X = posX
	e.Y = posY
	e.Map = ma

	e.Sprite, _ = sf.NewSprite(EntitiesTexture)

	e.Sprite.SetTextureRect(sf.IntRect{SpriteSize * spriteX, SpriteSize * spriteY, SpriteSize, SpriteSize})
	e.Sprite.SetPosition(sf.Vector2f{float32(e.X * SpriteSize), float32(e.Y * SpriteSize)})

	return e
}

func (e *Entity) Move(x, y int) {
	if !e.Map.IsBlocked(e.X+x, e.Y+y) {
		e.X += x
		e.Y += y
		e.Sprite.Move(sf.Vector2f{float32(SpriteSize * x), float32(SpriteSize * y)})
	}
}

func (e *Entity) Draw(w *sf.RenderWindow) {
	e.Sprite.Draw(w, sf.DefaultRenderStates())
}
