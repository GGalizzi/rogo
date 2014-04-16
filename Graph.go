package main

import sf "bitbucket.org/krepa098/gosfml2"

//Graph represents an SFML sprite together with the size and texture image.
type Graph struct {
	*sf.Sprite

	size int
}

//NewGraph initializes a Graph struct with the size read on the configuration fille, and to the correct
//TextureRect of that file.
func NewGraph(x, y int) *Graph {
	graph := new(Graph)
	var err error
	graph.Sprite, err = sf.NewSprite(SpriteSheet)
	if err != nil {
		panic(err)
	}
	graph.size = readSettings().SpriteSize
	graph.setSprite(x, y)
	return graph
}

func (gr *Graph) setSprite(x, y int) {
	gr.SetTextureRect(sf.IntRect{gr.size * x, gr.size * y, gr.size, gr.size})
}

func (gr *Graph) setPosition(x, y int) {

	gr.SetPosition(sf.Vector2f{float32(gr.size * x), float32(gr.size * y)})
}
