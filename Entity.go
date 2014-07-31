package main

import (
	"fmt"

	sf "bitbucket.org/krepa098/gosfml2"
)

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

	*Mob
	*Item

	sprite *Graph
}

type stats struct {
	maxhp int
	curhp int
	atk   int
	def   int
}

//Mob contains the data that an entity of type mob can use, meaning, any NPC.
type Mob struct {
	*stats

	pouch SoulPouch

	inventory Inventory
	faction   []Faction
}

//NewEntity initializes an Entity with the given data.
func NewEntity(name string, spriteX, spriteY, posX, posY int) *Entity {

	sprite := NewGraph(spriteX, spriteY)

	sprite.setSprite(spriteX, spriteY)
	sprite.SetPosition(sf.Vector2f{float32(posX * sprite.size), float32(posY * sprite.size)})

	m := new(Mob)
	m.stats = new(stats)
	m.maxhp, m.curhp = 130, 130
	m.atk = 10
	m.def = 4
	m.faction = append(m.faction, PLAYER)
	m.inventory = make(Inventory)

	return &Entity{x: posX, y: posY, sprite: sprite, Mob: m, name: name}
}

//NewEntityFromFile initializes an Entity with the data stored in the given JSON file.
func NewEntityFromFile(name string, x, y int) *Entity {

	data := ReadJSON("entities", name)
	e := &Entity{x: x, y: y, name: name}

	sx, sy := int(data["spriteX"].(float64)), int(data["spriteY"].(float64))

	//e.sprite, _ = sf.NewGraph(EntitiesTexture)
	e.sprite = NewGraph(sx, sy)
	e.sprite.SetPosition(sf.Vector2f{float32(e.x * e.sprite.size), float32(e.y * e.sprite.size)})

	e.Mob = nil
	if data["type"].(string) == "mob" {
		e.Mob = new(Mob)
		e.stats = new(stats)
		e.maxhp, e.curhp = int(data["hp"].(float64)), int(data["hp"].(float64))
		e.atk = int(data["atk"].(float64))
		e.def = int(data["def"].(float64))

		e.faction = make([]Faction, 1)
		//e.faction = append(e.faction, data["faction"].([]interface{})...)
		for _, v := range data["faction"].([]interface{}) {
			e.faction = append(e.faction, Faction(v.(string)))
		}
	}
	e.Item = nil
	if data["type"].(string) == "item" {
		e.Item = new(Item)
		e.Item.name = e.name
		e.stack = 1
		e.itype = ItemType(data["itemType"].(string))
		switch e.itype {
		case "potion":
			e.effect = potionEffect
			e.potency = int(data["potency"].(float64))
		case "key":
		}
	}

	return e
}

//Move should take ints between -1 and 1. That is, the direction where to move.
//To specify any tile in the map Place or SetPosition should be used.
func (e *Entity) Move(x, y int, g *Game) {

	dx := e.x + x
	dy := e.y + y

	ents := append(g.area.mobs, g.area.items...)

	for _, oe := range ents {
		if dx == oe.Position().X && dy == oe.Position().Y {
			if e.name == "cursor" {
				g.describe(oe)
			} else if oe.Mob != nil {
				e.attack(oe)
				return
			}
		}
	}
	if !g.area.IsBlocked(dx, dy) {
		e.Place(dx, dy)
		return
	}

	if g.area.isDoor(dx, dy) {
		e.tryOpen(g.area.tiles[dx+dy*g.area.width])
	}
}

func (e *Entity) moveTowards(oe *Entity, g *Game) {
	if e.isAlliedWith(oe) {
		return
	}
	ep := e.Position()   //Entity Position
	oep := oe.Position() //Other Entity Position.

	dx, dy := 0, 0

	switch {
	case ep.X < oep.X:
		dx = 1
	case ep.X > oep.X:
		dx = -1
	case ep.X == oep.X:
		dx = 0
	}

	switch {
	case ep.Y < oep.Y:
		dy = 1
	case ep.Y > oep.Y:
		dy = -1
	case ep.Y == oep.Y:
		dy = 0
	}

	e.Move(dx, dy, g)
}

func (attacker *Entity) attack(defender *Entity) {
	if !attacker.isAlliedWith(defender) {
		curhp := defender.curhp
		afterhp := curhp - (attacker.atk - defender.def)
		if afterhp <= 0 {
			attacker.absorb(defender.getSoul())
			defender.die()
			return
		}
		if afterhp > curhp {
			defender.curhp = curhp
			return
		}
		defender.curhp = afterhp
		damaged := curhp - defender.curhp
		log(fmt.Sprintf("%v attacks %v for %v damage.", attacker.name, defender.name, damaged))

	}
}

func (e *Entity) die() {
	e.Mob = nil
	e.sprite.SetColor(sf.ColorRed())

	//Becomes a corpse
	e.name = e.name + "'s corpse"
	e.Item = new(Item)
}

func (e *Entity) pickUp(i *Entity) {
	log(fmt.Sprintf("You pickup: %v", i.name))
	if e.inventory[i.name] != nil {
		e.inventory[i.name].stack++ //Add a stack to it if we already had the item.
		return
	}
	e.inventory[i.name] = i.Item
}

func (e *Entity) use(i *Item) {
	i.effect(i, e.Mob)
}

func (m *Mob) heal(amount int) {
	m.curhp += amount
	if m.curhp > m.maxhp {
		m.curhp = m.maxhp
	}
}

func (m *Mob) absorb(soul *Soul) {
	m.pouch = append(m.pouch, soul)
}

func (m *Entity) getSoul() *Soul {
	soul := new(Soul)
	soul.stats = new(stats)
	soul.name = soul.genName(m.name)

	percOfPower := 0.1

	soul.maxhp = int(float64(m.maxhp) * percOfPower)
	soul.atk = int(float64(m.atk) * percOfPower)
	soul.def = int(float64(m.def) * percOfPower)

	return soul
}

func (e *Entity) tryOpen(t *Tile) {
	if t.locked {
		for _, v := range e.inventory {
			if v.itype == KEY && v.linkedDoor == t {
				t.blocks = false
				t.locked = false
				t.setSprite(1, 9)
				log("You unlock the door")
				return
			}
		}
		log("The door is locked.")
		return
	}

	t.blocks = false
	t.setSprite(1, 9)
	return
}

func (e *Entity) isAlliedWith(oe *Entity) bool {

	if e.Mob == nil || oe.Mob == nil {
		fmt.Printf("Tried to see if someone is allied with someone, one of those was an invalid pointer. %+v isAlliedWith? %+v", e.Mob, oe.Mob)
		return true
	}
	tef := e.faction
	toef := oe.faction
	for _, ef := range tef {
		for _, oef := range toef {
			if ef == oef {
				return true
			}
		}
	}
	return false
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
