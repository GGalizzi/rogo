package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
)

//Entity contains the data that represents any entity that can appear on an Area that is not a tile.
type Entity struct {
	x int
	y int

	name string

	area *Area
	*Mob

	sprite *Graph
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

	sprite := NewGraph(spriteX, spriteY)

	sprite.setSprite(spriteX, spriteY)
	sprite.SetPosition(sf.Vector2f{float32(posX * sprite.size), float32(posY * sprite.size)})

	m := new(Mob)
	m.maxhp, m.curhp = 120, 120
	m.atk = 10
	m.def = 5

	return &Entity{x: posX, y: posY, area: a, sprite: sprite, Mob: m, name: "You"}
}

//NewEntityFromFile initializes an Entity with the data stored in the given JSON file.
func NewEntityFromFile(name string, x, y int, a *Area) *Entity {

	data := ReadJSON("entities", name)
	e := &Entity{x: x, y: y, area: a, name: name}

	sx, sy := int(data["spriteX"].(float64)), int(data["spriteY"].(float64))

	//e.sprite, _ = sf.NewGraph(EntitiesTexture)
	e.sprite = NewGraph(sx, sy)
	e.sprite.SetPosition(sf.Vector2f{float32(e.x * e.sprite.size), float32(e.y * e.sprite.size)})

	e.Mob = nil
	if data["type"].(string) == "mob" {
		e.Mob = new(Mob)
		e.maxhp, e.curhp = int(data["hp"].(float64)), int(data["hp"].(float64))
		e.atk = int(data["atk"].(float64))
		e.def = int(data["def"].(float64))
	}

	return e
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

func (attacker *Entity) attack(defender *Entity) {
	defender.curhp -= attacker.atk - defender.def
	if defender.curhp <= 0 {
		defender.die()
	}
}

func (e *Entity) die() {
	e.Mob = nil
	e.sprite.SetColor(sf.ColorRed())
	//TODO: make him an item.
}

//Draw draws the sprite on the window.
func (e *Entity) Draw(w *sf.RenderWindow) {
	e.sprite.Draw(w, sf.DefaultRenderStates())
}

//Place places the entity in the given coordinates, as well as set the sprite position to its correct place.
func (e *Entity) Place(x, y int) {
	e.x = x
	e.y = y
	e.sprite.setPosition(x, y)
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
