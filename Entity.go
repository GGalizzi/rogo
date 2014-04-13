package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

var (
	//EntitiesTexture is the file data which contains all the sprites that should be used.
	EntitiesTexture, _ = sf.NewTextureFromFile("ascii.png", nil)
)

//Entity contains the data that represents any entity that can appear on an Area that is not a tile.
type Entity struct {
	x int
	y int

	area *Area
	*Mob

	sprite *sf.Sprite
}

//Mob contains the data that an entity of type mob can use, meaning, any NPC.
type Mob struct {
	maxhp int
	curhp int
	atk   int
	def   int
}

//NewEntity initializes an Entity with the given data.
func NewEntity(spriteX, spriteY, posX, posY int, a *Area) *Entity {

	sprite, err := sf.NewSprite(EntitiesTexture)
	if err != nil {
		panic(err)
	}

	sprite.SetTextureRect(sf.IntRect{ReadSettings().SpriteSize * spriteX, ReadSettings().SpriteSize * spriteY, ReadSettings().SpriteSize, ReadSettings().SpriteSize})
	sprite.SetPosition(sf.Vector2f{float32(posX * ReadSettings().SpriteSize), float32(posY * ReadSettings().SpriteSize)})

	return &Entity{x: posX, y: posY, area: a, sprite: sprite}
}

//NewEntityFromFile initializes an Entity with the data stored in the JSON file given.
func NewEntityFromFile(name string, x, y int, a *Area) *Entity {

	data := ReadJSON("entities", name)
	e := &Entity{x: x, y: y, area: a}

	sx, sy := int(data["spriteX"].(float64)), int(data["spriteY"].(float64))

	e.sprite, _ = sf.NewSprite(EntitiesTexture)
	SetSprite(e, sx, sy)
	e.sprite.SetPosition(sf.Vector2f{float32(e.x * ReadSettings().SpriteSize), float32(e.y * ReadSettings().SpriteSize)})

	if data["type"].(string) == "mob" {
		e.Mob = new(Mob)
		e.maxhp, e.curhp = int(data["hp"].(float64)), int(data["hp"].(float64))
		e.atk = int(data["atk"].(float64))
		e.def = int(data["def"].(float64))
	}

	return e
}

//Sprite to implements Spriter interface, returns the entities sprite.
func (e *Entity) Sprite() *sf.Sprite {
	return e.sprite
}

//Move should take ints between -1 and 1. That is, the direction where to move.
//To specify any tile in the map Place or SetPosition should be used.
func (e *Entity) Move(x, y int) {
	if !e.area.IsBlocked(e.x+x, e.y+y) {
		dx := e.x + x
		dy := e.y + y
		e.Place(dx, dy)
	}
}

//Draw draws the sprite on the window.
func (e *Entity) Draw(w *sf.RenderWindow) {
	e.sprite.Draw(w, sf.DefaultRenderStates())
}

//Place places the entity in the given coordinates, as well as set the sprite position to its correct place.
func (e *Entity) Place(x, y int) {
	e.x = x
	e.y = y
	e.sprite.SetPosition(sf.Vector2f{float32(ReadSettings().SpriteSize * e.x), float32(ReadSettings().SpriteSize * e.y)})
}

//PosVector returns the position of the sprite, without using the tiled coordinate system, but the position based on the pixels of the window.
func (e *Entity) PosVector() sf.Vector2f {
	return e.sprite.GetPosition()
}

//Position returns the position of the entity in the tile coordinate system.
func (e *Entity) Position() sf.Vector2i {
	return sf.Vector2i{e.x, e.y}
}

//SetPosition places the entity in the given coordinates, as well as set the sprite position to its correct place.
func (e *Entity) SetPosition(pos sf.Vector2i) {
	e.Place(pos.X, pos.Y)
}
