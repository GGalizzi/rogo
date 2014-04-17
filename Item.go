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
	effect func(*Item, *Mob)
	itype  ItemType

	name    string
	potency int
	stack   int
}

type Inventory map[string]*Item

func potionEffect(i *Item, m *Mob) {
	m.heal(i.potency)
	log("Drank " + i.name)
}
