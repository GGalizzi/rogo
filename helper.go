package main

import sf "bitbucket.org/krepa098/gosfml2"

type Spriter interface {
	GetSprite() *sf.Sprite
}

func SetSprite(obj Spriter, x, y int) {
	obj.GetSprite().SetTextureRect(sf.IntRect{ReadSettings().SpriteSize * x, ReadSettings().SpriteSize * y, ReadSettings().SpriteSize, ReadSettings().SpriteSize})
}
