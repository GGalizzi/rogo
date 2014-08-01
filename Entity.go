package main

import sf "bitbucket.org/krepa098/gosfml2"

//Faction represents the different groups of factions an NPC or player can belong to.
type Faction string

const (
	//ORCS faction belongs to any orc.
	ORCS Faction = "orcs"
	//PLAYER faction represents allies to the player, and the player itself.
	PLAYER Faction = "player"
)

//Entity contains the data that represents any entity that can appear on an Area that is not a tile.
type Entity struct {
	x int
	y int

	name string

	sprite *Graph
}

type stats struct {
	maxhp int
	curhp int
	atk   int
	def   int
}

func NewEntity(name string, spriteX, spriteY, posX, posY int) *Entity {

	sprite := NewGraph(spriteX, spriteY)

	sprite.setSprite(spriteX, spriteY)
	sprite.SetPosition(sf.Vector2f{float32(posX * sprite.size), float32(posY * sprite.size)})

	return &Entity{x: posX, y: posY, sprite: sprite, name: name}
}

//NewEntityFromFile initializes an Entity with the data stored in the given JSON file.
func NewEntityFromFile(name string, x, y int) (e *Entity, data map[string]interface{}) {
	data = ReadJSON("entities", name)
	e = &Entity{x: x, y: y, name: name}

	sx, sy := int(data["spriteX"].(float64)), int(data["spriteY"].(float64))

	//e.sprite, _ = sf.NewGraph(EntitiesTexture)
	e.sprite = NewGraph(sx, sy)
	e.sprite.SetPosition(sf.Vector2f{float32(e.x * e.sprite.size), float32(e.y * e.sprite.size)})

	return
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
