package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"encoding/json"
	"os"
)

var (
	EntitiesTexture, _ = sf.NewTextureFromFile("ascii.png", nil)
)

type Entity struct {
	sprite *sf.Sprite
	x      int
	y      int

	area *Area

	*Mob
}

type Mob struct {
	maxhp int
	curhp int
	atk   int
	def   int
}

func NewEntity(spriteX, spriteY, posX, posY int, a *Area) *Entity {
	e := new(Entity)

	e.x = posX
	e.y = posY
	e.area = a

	e.sprite, _ = sf.NewSprite(EntitiesTexture)

	e.sprite.SetTextureRect(sf.IntRect{ReadSettings().SpriteSize * spriteX, ReadSettings().SpriteSize * spriteY, ReadSettings().SpriteSize, ReadSettings().SpriteSize})
	e.sprite.SetPosition(sf.Vector2f{float32(e.x * ReadSettings().SpriteSize), float32(e.y * ReadSettings().SpriteSize)})

	return e
}

func NewEntityFromFile(name string, x, y int, a *Area) *Entity {

	file, err := os.Open("entities/" + name + ".ent")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jParser := json.NewDecoder(file)

	var d interface{}

	if err = jParser.Decode(&d); err != nil {
		panic(err)
	}

	data := d.(map[string]interface{})

	e := new(Entity)

	e.x = x
	e.y = y
	e.area = a

	sx, sy := int(data["spriteX"].(float64)), int(data["spriteY"].(float64))

	e.sprite, _ = sf.NewSprite(EntitiesTexture)
	SetSprite(e, sx, sy)
	e.sprite.SetPosition(sf.Vector2f{float32(e.x * ReadSettings().SpriteSize), float32(e.y * ReadSettings().SpriteSize)})

	if data["type"].(string) == "mob" {
		e.Mob = new(Mob)
		e.maxhp, e.curhp = int(data["hp"].(float64)), int(data["hp"].(float64))
		e.atk = int(data["atk"].(float64))
		e.def = int(data["def"].(float64))
	}

	return e
}

func (e *Entity) GetSprite() *sf.Sprite {
	return e.sprite
}

//Move should take ints between -1 and 1. That is, the direction where to move.
//To specify any tile in the map Place should be used.
func (e *Entity) Move(x, y int) {
	if !e.area.IsBlocked(e.x+x, e.y+y) {
		e.x += x
		e.y += y
		e.sprite.Move(sf.Vector2f{float32(ReadSettings().SpriteSize * x), float32(ReadSettings().SpriteSize * y)})
	}
}

func (e *Entity) Draw(w *sf.RenderWindow) {
	e.sprite.Draw(w, sf.DefaultRenderStates())
}

func (e *Entity) Place(x, y int) {
	e.x = x
	e.y = y
}

func (e *Entity) PosVector() sf.Vector2f {
	return e.sprite.GetPosition()
}

func (e *Entity) Position() sf.Vector2i {
	return sf.Vector2i{e.x, e.y}
}
