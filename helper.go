package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
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

//Settings struct defines all the variable settings of the game, which are to be stored in a JSON file.
type Settings struct {
	Resolution string `json:"resolution"`
	SpriteSize int
	resW       uint
	resH       uint
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
