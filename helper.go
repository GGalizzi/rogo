package main

import (
	sf "bitbucket.org/krepa098/gosfml2"
	"encoding/json"
	"os"
	"strconv"
	"strings"
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

type Settings struct {
	Resolution string `json:"resolution"`
	SpriteSize int
	resW       uint
	resH       uint
}

func ReadSettings() Settings {
	file, err := os.Open("conf.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jParser := json.NewDecoder(file)
	var s Settings
	if err = jParser.Decode(&s); err != nil {
		panic(err)
	}
	rs := strings.Split(s.Resolution, "x")
	resW, err := strconv.Atoi(rs[0])
	s.resW = uint(resW)
	if err != nil {
		panic(s)
	}

	resH, err := strconv.Atoi(rs[1])
	s.resH = uint(resH)
	if err != nil {
		panic(err)
	}

	return s
}
