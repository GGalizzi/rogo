package main

import (
	"fmt"

	sf "bitbucket.org/krepa098/gosfml2"
)

//Mob contains the data that an entity of type mob can use, meaning, any NPC.
type Mob struct {
	*stats
	*Entity

	pouch SoulPouch

	inventory Inventory
	faction   []Faction
}

//NewEntity initializes an Entity with the given data.
func NewMob(name string, spriteX, spriteY, posX, posY int) *Mob {

	m := new(Mob)
	m.stats = new(stats)
	m.maxhp, m.curhp = 130, 130
	m.atk = 10
	m.def = 4
	m.faction = append(m.faction, PLAYER)
	m.inventory = make(Inventory)

	m.Entity = NewEntity(name, spriteX, spriteY, posX, posY)

	return m
}

func NewMobFromFile(name string, x, y int) *Mob {
	m := new(Mob)
	ent, data := NewEntityFromFile(name, x, y)
	m.Entity = ent

	if data["type"].(string) == "mob" {
		m.stats = new(stats)
		m.maxhp, m.curhp = int(data["hp"].(float64)), int(data["hp"].(float64))
		m.atk = int(data["atk"].(float64))
		m.def = int(data["def"].(float64))

		m.faction = make([]Faction, 1)
		//e.faction = append(e.faction, data["faction"].([]interface{})...)
		for _, v := range data["faction"].([]interface{}) {
			m.faction = append(m.faction, Faction(v.(string)))
		}
	}

	return m
}

//Move should take ints between -1 and 1. That is, the direction where to move.
//To specify any tile in the map Place or SetPosition should be used.
func (e *Mob) Move(x, y int, g *Game) {

	dx := e.x + x
	dy := e.y + y

	var om *Mob
	var oi *Item

	var longest int
	if len(g.area.mobs) > len(g.area.items) {
		longest = len(g.area.mobs)
	} else {
		longest = len(g.area.items)
	}

	for i := 0; i < longest; i++ {
		if i < len(g.area.mobs) {
			om = g.area.mobs[i]
		} else {
			om = nil
		}

		//Only check items if we are moving the cursor.
		if e.name == "cursor" && i < len(g.area.items) {
			oi = g.area.items[i]
		} else {
			oi = nil
		}

		if om != nil && (dx == om.Position().X && dy == om.Position().Y) {
			if e.name == "cursor" {
				g.describe(om.Entity)
			} else {
				e.attack(om)
				return
			}
		}

		//Only check for items if we are moving the cursor.
		if e.name == "cursor" && oi != nil && (dx == oi.Position().X && dy == oi.Position().Y) {
			g.describe(oi.Entity)
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

func (e *Mob) moveTowards(oe *Entity, g *Game) {

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

func (attacker *Mob) attack(defender *Mob) {
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

func (m *Mob) die() {
	m.sprite.SetColor(sf.ColorRed())
}

func (m *Mob) pickUp(i *Item) {
	log(fmt.Sprintf("You pickup: %v", i.name))
	if m.inventory[i.name] != nil {
		m.inventory[i.name].stack++ //Add a stack to it if we already had the item.
		return
	}
	m.inventory[i.name] = i
}

func (m *Mob) use(i *Item) {
	i.effect(i, m)
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

func (m *Mob) getSoul() *Soul {
	soul := new(Soul)
	soul.stats = new(stats)
	soul.name = soul.genName(m.name)

	percOfPower := 0.1

	soul.maxhp = int(float64(m.maxhp) * percOfPower)
	soul.atk = int(float64(m.atk) * percOfPower)
	soul.def = int(float64(m.def) * percOfPower)

	return soul
}

func (m *Mob) tryOpen(t *Tile) {
	if t.locked {
		for _, v := range m.inventory {
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

func (e *Mob) isAlliedWith(oe *Mob) bool {

	if e == nil || oe == nil {
		fmt.Printf("Tried to see if someone is allied with someone, one of those was an invalid pointer. %+v isAlliedWith? %+v", e, oe)
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
