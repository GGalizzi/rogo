package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"encoding/json"
	"os"
)

type Spriter interface {
	GetSprite() *sf.Sprite
}

func SetSprite(obj Spriter, x, y int) {
	obj.GetSprite().SetTextureRect(sf.IntRect{ReadSettings().SpriteSize * x, ReadSettings().SpriteSize * y, ReadSettings().SpriteSize, ReadSettings().SpriteSize})
}

func ReadJson(folder, name string) map[string]interface{} {
	file, err := os.Open(folder + "/" + name + ".json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	jParser := json.NewDecoder(file)

	var t interface{}

	if err = jParser.Decode(&t); err != nil {
		panic(err)
	}

	data := t.(map[string]interface{})

	return data
}
