package main

import "testing"

func TestInventory(t *testing.T) {
	g := MockNewGame()
	a := PrepareArea()

	g.items = append(g.items, NewEntityFromFile("potion", 3, 3, a))
	g.player.Place(3, 3)
	g.mobs = append(g.mobs, g.player)

	g.handleInput('g')

	if g.player.inventory["potion"] == nil {
		t.Errorf("Expected player to have potion in inventory: %v", g.player.inventory)
	}

	g.items = append(g.items, NewEntityFromFile("potion", 3, 3, a))

	g.handleInput('g')

	if g.player.inventory["potion"].stack != 2 {
		t.Errorf("Player should have two potions, it has: %d", g.player.inventory["potion"].stack)
	}

	t.Logf("Player should have potions in inventory: %+v", g.player.inventory)
}

func TestPotion(t *testing.T) {

	g := MockNewGame()
	a := PrepareArea()

	potion := NewEntityFromFile("potion", 3, 3, a)
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
