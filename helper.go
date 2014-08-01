package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	sf "bitbucket.org/krepa098/gosfml2"
)

var (
	//SpriteSheet is the file data which contains all the sprites that should be used.
	SpriteSheet, _ = sf.NewTextureFromFile("ascii.png", nil)
	Font, _        = sf.NewFontFromFile("font.ttf")
)

//ReadJSON reads the given json file in the given folder, and returns a map of any type, representing the JSON.
func ReadJSON(folder, name string) map[string]interface{} {
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

func appendString(text *sf.Text, s string) {
	switch text.GetString() {
	case "":
		text.SetString(s)
	default:
		text.SetString(text.GetString() + "\n" + s)
	}
}

func log(s string) {
	logFile, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logFile.WriteString(s + "\n")
	if err != nil {
		panic(err)
	}
}

func removeFromList(a interface{}, i int) interface{} {
	//a[len(a)-1], a[i], a = nil, a[len(a)-1], a[:len(a)-1] // Deletes completely
	switch t := a.(type) {
	case []*Mob:
		fmt.Printf("type: %v", t)
		if len(t) <= i {
			fmt.Printf("Trying to remove something that doesn't exist? %v [%d]\n", a, i)
			return t
		}
		t[i], t = t[len(t)-1], t[:len(t)-1]

		return a
	}

	return a
}

//Settings struct defines all the variable settings of the game, which are to be stored in a JSON file.
type Settings struct {
	Resolution string `json:"resolution"`
	SpriteSize int
	resW       float32
	resH       float32
}

//readSettings reads the settings file, and returns a struct with that data.
func readSettings() Settings {
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
	s.resW = float32(resW)
	if err != nil {
		panic(s)
	}

	resH, err := strconv.Atoi(rs[1])
	s.resH = float32(resH)
	if err != nil {
		panic(err)
	}

	return s
}
