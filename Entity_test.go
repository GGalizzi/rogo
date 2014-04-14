package main

import (
	"testing"

	sf "bitbucket.org/krepa098/gosfml2"
)

func TestNewMobFromFile(t *testing.T) {
	a := PrepareArea()
	e := NewEntityFromFile("orc", 3, 3, a)

	ss := readSettings().SpriteSize
	if pos, vec := e.Position(), e.PosVector(); pos.X != 3 || pos.Y != 3 || vec.X != float32(3*ss) || vec.Y != float32(3*ss) {
		t.Errorf("Expected position to be 3,3. pos:%v,%v; vec:%v,%v)", pos.X, pos.Y, vec.X, vec.Y)
	}

	// An orc should have 30hp, 5atk, 2def
	if maxhp, curhp, atk, def := e.maxhp, e.curhp, e.atk, e.def; maxhp != 30 || curhp != 30 || atk != 5 || def != 2 {
		t.Errorf("Expected: maxhp: 30, curhp: 30, atk: 5, def: 2; Got: maxhp: %v, curhp: %v, atk: %v, def: %v", maxhp, curhp, atk, def)
	}
}

func TestBasicAi(t *testing.T) {
	a := PrepareArea()
	e := NewEntityFromFile("orc", 3, 3, a)
	g := MockNewGame()

	op := e.Position()
	opv := e.PosVector()

	g.entities = append(g.entities, e, g.player)

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
	g.processAI(e)
	aph := g.player.curhp

	if ap := e.Position(); ap.X == 5 || ap.Y == 5 {
		t.Errorf("Expected orc not to move, should be attacking. Actual Pos: %v", ap)
	}

	if oph == aph {
		t.Errorf("Expected player to lose health. Before Process: %v, After Process: %v", oph, aph)
	}
}

func TestFactionAttack(t *testing.T) {
	a := PrepareArea()

	orc1 := NewEntityFromFile("orc", 3, 3, a)
	orc2 := NewEntityFromFile("orc", 3, 4, a)

	g := MockNewGame()
	g.entities = append(g.entities, orc1, orc2)

	orc2hp := orc2.curhp
	orc1.moveTowards(orc2, g)

	if actual := orc2.curhp; actual != orc2hp {
		t.Errorf("Mobs within the same faction shouldn't hit each other when moving. Pre-HP: %v, Post-HP: %v", orc2hp, actual)
	}
}
