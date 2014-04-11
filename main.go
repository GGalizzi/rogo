package main

import (
  sf "bitbucket.org/krepa098/gosfml2"
  "runtime"
)
const SpriteSize = 16

var G *Game

func init() {
  runtime.LockOSThread()
}

func main() {
  G = NewGame()

  G.Texture,_ = sf.NewTextureFromFile("ascii.png",nil)

  G.run()
}
