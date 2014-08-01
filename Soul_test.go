package main

import "testing"

func TestSoulAbsorb(t *testing.T) {

	p := NewMob("player", -1, 0, 2, 3)
	e := NewMob("enemy", -1, 0, 3, 3)

	e.curhp = 1
	e.faction = nil

	prevLen := len(p.pouch)
	//When player kills en enemy, he gets his soul into the pouch.
	p.attack(e)
	afterLen := len(p.pouch)
	if afterLen != prevLen+1 {
		t.Errorf("Expected the amount of souls in player pouch to be 1 more than before. Was: %d Now: %d", prevLen, afterLen)
	}

	t.Logf("Pouch len. Was: %d, Now: %d", prevLen, afterLen)
}

func TestGetSoul(t *testing.T) {
	mob := NewMob("mob", -1, 0, 2, 2)

	mob.maxhp = 93
	mob.atk = 13
	mob.def = 193

	soul := mob.getSoul()

	perc := 0.1
	if soul.maxhp != int(float64(mob.maxhp)*perc) || soul.atk != int(float64(mob.atk)*perc) || soul.def != int(float64(mob.def)*perc) {
		t.Errorf("Expected the souls stats to be 10percent of the mob it came from.\nsoul: %v\n mob: %v\n", soul.stats, mob.stats)
	}
	t.Logf("Soul power is: %v, from: %v", soul.stats, mob.stats)
}
