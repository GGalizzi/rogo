package main

import sf "bitbucket.org/krepa098/gosfml2"

type Map struct {
	Tiles  []*Tile
	Width  int
	Height int
}

func NewMap() *Map {
	m := new(Map)
	m.Width = 50
	m.Height = 20

	m.Tiles = make([]*Tile, m.Height*m.Width)

	for i, _ := range m.Tiles {
		m.Tiles[i] = NewTile()
	}

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if y == 0 || y == m.Height-1 || x == 0 || x == m.Width-1 {
				m.Tiles[x+y*m.Width].SetSprite(0, 9)
				m.Tiles[x+y*m.Width].Blocks = true
			} else {
				m.Tiles[x+y*m.Width].Blocks = false
			}
		}
	}

	return m
}

func (m *Map) Draw(window *sf.RenderWindow) {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.Tiles[x+y*m.Width].Draw(window, x, y)
		}
	}
}

func (m *Map) IsBlocked(x, y int) bool {
	return m.Tiles[x+y*m.Width].Blocks
}
