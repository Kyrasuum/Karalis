package rng

import (
	"math/rand"

	"karalis/pkg/lmath"
)

// --- Perlin helpers (classic improved Perlin) ---

func fade(t float32) float32 {
	// 6t^5 - 15t^4 + 10t^3
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float32) float32 {
	return a + t*(b-a)
}

func grad(hash byte, x, y, z float32) float32 {
	h := hash & 15
	var u float32
	if h < 8 {
		u = x
	} else {
		u = y
	}

	var v float32
	if h < 4 {
		v = y
	} else if h == 12 || h == 14 {
		v = x
	} else {
		v = z
	}

	var res float32
	if (h & 1) == 0 {
		res = u
	} else {
		res = -u
	}
	if (h & 2) == 0 {
		res += v
	} else {
		res += -v
	}
	return res
}

// noise returns roughly [-1, 1]
func noise(x, y, z float32, p []byte) float32 {
	// Use floor; wrap to 0..255
	xf := float32(lmath.Floor(x))
	yf := float32(lmath.Floor(y))
	zf := float32(lmath.Floor(z))

	X := int(xf) & 255
	Y := int(yf) & 255
	Z := int(zf) & 255

	x -= xf
	y -= yf
	z -= zf

	u := fade(x)
	v := fade(y)
	w := fade(z)

	A := int(p[X]) + Y
	AA := int(p[A&255]) + Z
	AB := int(p[(A+1)&255]) + Z
	B := int(p[(X+1)&255]) + Y
	BA := int(p[B&255]) + Z
	BB := int(p[(B+1)&255]) + Z

	return lerp(w,
		lerp(v,
			lerp(u, grad(p[AA&255], x, y, z), grad(p[BA&255], x-1, y, z)),
			lerp(u, grad(p[AB&255], x, y-1, z), grad(p[BB&255], x-1, y-1, z)),
		),
		lerp(v,
			lerp(u, grad(p[(AA+1)&255], x, y, z-1), grad(p[(BA+1)&255], x-1, y, z-1)),
			lerp(u, grad(p[(AB+1)&255], x, y-1, z-1), grad(p[(BB+1)&255], x-1, y-1, z-1)),
		),
	)
}

// GeneratePermutation creates a 512-length permutation table, duplicated.
// Deterministic given seed.
func GeneratePermutation(seed int64) []byte {
	r := rand.New(rand.NewSource(seed))

	source := make([]byte, 256)
	for i := 0; i < 256; i++ {
		source[i] = byte(i)
	}

	// Fisher-Yates shuffle
	for i := 255; i > 0; i-- {
		j := r.Intn(i + 1)
		source[i], source[j] = source[j], source[i]
	}

	perm := make([]byte, 512)
	copy(perm[:256], source)
	copy(perm[256:], source)
	return perm
}

// FBM returns sum of octaves of noise at (x,y) with z=0
func FBM(x, y float32, octaves uint8, persistence, lacunarity, scale float32, p []byte) float32 {
	var value float32 = 0
	var amp float32 = 1
	var freq float32 = 1

	// prevent divide-by-zero
	if scale == 0 {
		scale = 1
	}

	for i := uint8(0); i < octaves; i++ {
		value += amp * noise(x*freq/scale, y*freq/scale, 0, p)
		amp *= persistence
		freq *= lacunarity
	}
	return value
}

// WeightedAverageHeight mimics your gaussian-ish weighted average around rounded (x,z).
// Note: your Zig loop variables look a bit inconsistent (dz from -radius to 1) —
// this version uses an integer neighborhood derived from radius.
func WeightedAverageHeight(terrain []float32, x, z, radius float32, size int, alpha float32) (float32, bool) {
	if size <= 0 || len(terrain) < size*size {
		return 0, false
	}
	if radius <= 0 {
		// fall back to nearest
		ix := lmath.Round(x)
		iz := lmath.Round(z)
		if ix < 0 || iz < 0 || ix >= size || iz >= size {
			return 0, false
		}
		return terrain[iz*size+ix], true
	}

	cx := lmath.Round(x)
	cz := lmath.Round(z)

	r := lmath.Ceil(radius)
	var sum float32
	var total float32

	for dz := -r; dz <= r; dz++ {
		for dx := -r; dx <= r; dx++ {
			nx := cx + dx
			nz := cz + dz
			if nx < 0 || nz < 0 || nx >= size || nz >= size {
				continue
			}
			diffx := nx - int(x)
			diffz := nz - int(z)
			dist2 := diffx*diffx + diffz*diffz
			w := lmath.Exp(-alpha * float32(dist2))

			ix := int(nx)
			iz := int(nz)
			h := terrain[iz*size+ix]

			sum += h * w
			total += w
		}
	}

	if total <= 0 {
		return 0, false
	}
	return sum / total, true
}
