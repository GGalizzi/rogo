package main

import (
	"runtime"
)

const SpriteSize = 16

var G *Game

func init() {
	runtime.LockOSThread()
}

func main() {
	G = NewGame()

	G.run()
}
