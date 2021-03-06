package main

import (
	"testing"

	sf "bitbucket.org/krepa098/gosfml2"
)

func TestNewMobFromFile(t *testing.T) {
	e := NewMobFromFile("orc", 3, 3)

	ss := readSettings().SpriteSize
	if pos, vec := e.Position(), e.PosVector(); pos.X != 3 || pos.Y != 3 || vec.X != float32(3*ss) || vec.Y != float32(3*ss) {
		t.Errorf("Expected position to be 3,3. pos:%v,%v; vec:%v,%v)", pos.X, pos.Y, vec.X, vec.Y)
	}

	// An orc should have 30hp, 5atk, 2def
	if maxhp, curhp, atk, def := e.maxhp, e.curhp, e.atk, e.def; maxhp != 30 || curhp != 30 || atk != 5 || def != 2 {
		t.Errorf("Expected: maxhp: 30, curhp: 30, atk: 5, def: 2; Got: maxhp: %v, curhp: %v, atk: %v, def: %v", maxhp, curhp, atk, def)
	}

	t.Logf("Mob was created as per orc.json data: %+v", *e)
}

func TestBasicAi(t *testing.T) {
	e := NewMobFromFile("orc", 3, 3)
	g := MockNewGame()

	op := e.Position()
	opv := e.PosVector()

	g.area.mobs = append(g.area.mobs, e, g.player)

	g.processAI(e)

	if ap, apv := e.Position(), e.PosVector(); ap == op || opv == apv {
		t.Errorf("Expected orc to move from 3,3 -> %v | Original Vector:%v should != Actual Vector: %v", ap, opv, apv)
	}

	if ap := e.Position(); ap.X != 4 || ap.Y != 4 {
		t.Errorf("Expected orc to move one tile only. Actual: %v", ap)
	}

	g.player.SetPosition(sf.Vector2i{5, 5})

	op = e.Position()
	oph := g.player.curhp
	g.player.def = 0
	g.processAI(e)
	aph := g.player.curhp

	if ap := e.Position(); ap.X == 5 || ap.Y == 5 {
		t.Errorf("Expected orc not to move, should be attacking. Actual Pos: %v", ap)
	}

	if aph >= oph {
		t.Errorf("Expected player to lose health. Before Process: %v, After Process: %v . atk(%v) - def(%v)", oph, aph, e.atk, g.player.def)
	}
}

func TestFactionAttack(t *testing.T) {

	orc1 := NewMobFromFile("orc", 3, 3)
	orc2 := NewMobFromFile("orc", 3, 4)

	g := MockNewGame()
	g.area.mobs = append(g.area.mobs, orc1, orc2)

	orc2hp := orc2.curhp
	orc1.moveTowards(orc2.Entity, g)

	if actual := orc2.curhp; actual != orc2hp {
		t.Errorf("Mobs within the same faction shouldn't hit each other when moving. Pre-HP: %v, Post-HP: %v", orc2hp, actual)
	}
}

func TestMobDeath(t *testing.T) {
	g := MockNewGame()
	g.player.Place(3, 3)
	orc := NewMobFromFile("orc", 3, 4)
	orc.curhp = 1
	g.area.mobs = append(g.area.mobs, orc, g.player)

	for i, m := range g.area.mobs {
		if m == orc {
			orc.die(g.area, i)
		}
	}

	found := false
	for _, m := range g.area.mobs {
		if m == orc {
			found = true
		}
	}

	if found {
		t.Errorf("Expected orc to die. i.e: Mob be nil, Item not be nil ->: %v", orc)
	}
}

func TestPlayerAttack(t *testing.T) {

	g := MockNewGame()

	g.player.SetPosition(sf.Vector2i{3, 3})
	orc := NewMobFromFile("orc", 3, 4)

	g.area.mobs = append(g.area.mobs, g.player, orc)

	prevHP := orc.curhp
	g.handleInput('2')
	postHP := orc.curhp

	if postHP >= prevHP {
		t.Errorf("Expected orc to lose health. Before player move: %v, After: %v | %v vs. %v", prevHP, postHP, g.player, orc)
	}
	t.Logf("Orcs health is reduced when player bumps into him. Orc health before: %v, after: %v", prevHP, postHP)
}
