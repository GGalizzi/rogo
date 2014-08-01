package main

//ItemType represents the type of item. Which can be one of the consts described below.
type ItemType string

const (
	//POTION is any item that can be consumed and has an effect on the consumer.
	POTION ItemType = "potion"
	//KEY is any item that can be used to open any kind of locked door
	KEY ItemType = "key"
)

type Item struct {
	*Entity

	name  string
	itype ItemType

	//Key stuff
	linkedDoor *Tile

	effect  func(*Item, *Mob)
	potency int
	stack   int
}

type Inventory map[string]*Item

func potionEffect(i *Item, m *Mob) {
	m.heal(i.potency)
	log("Drank " + i.name)
}

func NewItemFromFile(name string, x, y int) *Item {
	i := new(Item)
	ent, data := NewEntityFromFile(name, x, y)
	i.Entity = ent

	if data["type"].(string) == "item" {
		i.name = name
		i.stack = 1
		i.itype = ItemType(data["itemType"].(string))
		switch i.itype {
		case "potion":
			i.effect = potionEffect
			i.potency = int(data["potency"].(float64))
		case "key":
		}
	}

	return i
}
