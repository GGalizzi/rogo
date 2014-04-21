package main

import (
	"math"
	"math/rand"
)

type PerlinNoise []int

func NewPerlin(seed uint) PerlinNoise {

	randList := rand.Perm(256)

	p := PerlinNoise(randList)

	p = append(p, p...)

	return p
}

func (p PerlinNoise) noise(x, y, z float64) float64 {
	X := int(int(math.Floor(x)) & 255)
	Y := int(int(math.Floor(y)) & 255)
	Z := int(int(math.Floor(z)) & 255)

	x -= math.Floor(x)
	y -= math.Floor(y)
	z -= math.Floor(z)

	u := p.fade(x)
	v := p.fade(y)
	w := p.fade(z)

	A := p[X] + Y
	AA := p[A] + Z
	AB := p[A+1] + Z
	B := p[X+1] + Y
	BA := p[B] + Z
	BB := p[B+1] + Z

	//res := p.lerp(w, p.lerp(v, p.lerp(u, p.grad(p[AA],x,y,z), p.grad(p[BA], x-1, y, z)), p.lerp(u, p.grad(p[AB],x,y-1,z), p.grad(j

	res := p.lerp(w, p.lerp(v, p.lerp(u, p.grad(p[AA], x, y, z), p.grad(p[BA], x-1, y, z)), p.lerp(u, p.grad(p[AB], x, y-1, z), p.grad(p[BB], x-1, y-1, z))), p.lerp(v, p.lerp(u, p.grad(p[AA+1], x, y, z-1), p.grad(p[BA+1], x-1, y, z-1)), p.lerp(u, p.grad(p[AB+1], x, y-1, z-1), p.grad(p[BB+1], x-1, y-1, z-1))))

	return ((res + 1.0) / 2.0)
}

func (p *PerlinNoise) fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func (p *PerlinNoise) lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func (p *PerlinNoise) grad(hash int, x, y, z float64) float64 {
	h := hash & 15

	var u float64
	if h < 8 {
		u = x
	} else {
		u = y
	}

	var v float64
	if v < 4 {
		v = y
	} else {
		if h == 12 || h == 14 {
			v = x
		} else {
			v = z
		}
	}

	var ret1 float64

	if (h & 1) == 0 {
		ret1 = u
	} else {
		ret1 = -u
	}

	var ret2 float64

	if (h & 2) == 0 {
		ret2 = v
	} else {
		ret2 = -v
	}

	return ret1 + ret2
}
