package main

import sf "bitbucket.org/krepa098/gosfml2"

var (
	EntitiesTexture, _ = sf.NewTextureFromFile("ascii.png", nil)
)

type Entity struct {
	sprite *sf.Sprite
	x      int
	y      int

	area *Area
}

func NewEntity(spriteX, spriteY, posX, posY int, a *Area) *Entity {
	e := new(Entity)

	e.x = posX
	e.y = posY
	e.area = a

	e.sprite, _ = sf.NewSprite(EntitiesTexture)

	e.sprite.SetTextureRect(sf.IntRect{SpriteSize * spriteX, SpriteSize * spriteY, SpriteSize, SpriteSize})
	e.sprite.SetPosition(sf.Vector2f{float32(e.x * SpriteSize), float32(e.y * SpriteSize)})

	return e
}

func (e *Entity) Move(x, y int) {
	if !e.area.IsBlocked(e.x+x, e.y+y) {
		e.x += x
		e.y += y
		e.sprite.Move(sf.Vector2f{float32(SpriteSize * x), float32(SpriteSize * y)})
	}
}

func (e *Entity) Draw(w *sf.RenderWindow) {
	e.sprite.Draw(w, sf.DefaultRenderStates())
}

func (e *Entity) Place(x, y int) {
	e.x = x
	e.y = y
}

func (e *Entity) PosVector() sf.Vector2f {
	return e.sprite.GetPosition()
}

func (e *Entity) Position() sf.Vector2i {
	return sf.Vector2i{e.x, e.y}
}
