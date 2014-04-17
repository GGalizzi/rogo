package main

import "testing"

func TestInventory(t *testing.T) {
	g := MockNewGame()

	g.area.items = append(g.area.items, NewEntityFromFile("potion", 3, 3))
	g.player.Place(3, 3)
	g.area.mobs = append(g.area.mobs, g.player)

	g.handleInput('g')

	if g.player.inventory["potion"] == nil {
		t.Errorf("Expected player to have potion in inventory: %v", g.player.inventory)
	}

	g.area.items = append(g.area.items, NewEntityFromFile("potion", 3, 3))

	g.handleInput('g')

	if g.player.inventory["potion"].stack != 2 {
		t.Errorf("Player should have two potions, it has: %d", g.player.inventory["potion"].stack)
	}

	t.Logf("Player should have potions in inventory: %+v", g.player.inventory)
}

func TestKeys(t *testing.T) {
	g := MockNewGame()

	g.area.placeTile("lockedDoor", 12, 12)
	tile := g.area.tiles[12+12*g.area.width]

	var validKey *Entity
	validKey = nil

	for _, v := range g.area.items {
		if v.itype == KEY && v.linkedDoor == tile {
			validKey = v
		}
	}

	if validKey == nil {
		t.Errorf("A valid key should be placed somewhere in the map after a lockedDoor is created.")
	}

	t.Logf("A key was created after placing a lockedDoor on the area. %+v links to %p.", validKey.Item, tile)

	g.player.pickUp(validKey)

	g.player.Place(11, 12)
	g.handleInput('6')

	if tile.locked || tile.blocks {
		t.Errorf("The key should have unlocked the tile. %+v points to %p: %+v", validKey.Item, tile, tile)
	}

}

func TestPotion(t *testing.T) {

	g := MockNewGame()

	potion := NewEntityFromFile("potion", 3, 3)
	potion.potency = 10
	g.player.inventory["potion"] = potion.Item

	g.player.maxhp = 10
	g.player.curhp = 1
	prevHP := g.player.curhp
	g.player.use(g.player.inventory["potion"])

	if afterHP := g.player.curhp; afterHP != 10 {
		t.Errorf("Expected player to heal by using potion, without going over 10(max). Was at: %v, Is now at: %v", prevHP, afterHP)
	}

	t.Logf("Player healed %d -> %d", prevHP, g.player.curhp)

}
