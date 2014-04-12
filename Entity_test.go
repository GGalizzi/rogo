package main

import "testing"

func TestNewMobFromFile(t *testing.T) {
	a := PrepareArea()
	e := NewEntityFromFile("orc", 3, 3, a)

	ss := ReadSettings().SpriteSize
	if pos, vec := e.Position(), e.PosVector(); pos.X != 3 || pos.Y != 3 || vec.X != float32(3*ss) || vec.Y != float32(3*ss) {
		t.Errorf("Expected position to be 3,3. pos:%v,%v; vec:%v,%v)", pos.X, pos.Y, vec.X, vec.Y)
	}

	// An orc should have 30hp, 5atk, 2def
	if maxhp, curhp, atk, def := e.maxhp, e.curhp, e.atk, e.def; maxhp != 30 || curhp != 30 || atk != 5 || def != 2 {
		t.Errorf("Expected: maxhp: 30, curhp: 30, atk: 5, def: 2; Got: maxhp: %v, curhp: %v, atk: %v, def: %v", maxhp, curhp, atk, def)
	}
}
