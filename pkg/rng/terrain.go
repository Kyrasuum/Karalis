package rng

import (
	"image/color"

	"karalis/pkg/lmath"
)

type TerrainConfig struct {
	Seed        int64
	Size        int       // map is Size x Size
	BaseMap     []float32 // optional; if nil -> no base
	NoiseWeight float32
	Octaves     uint8
	Persistence float32
	Lacunarity  float32
	Scale       float32
	WorldStartX float32
	WorldStartZ float32
	WorldStep   float32
}

// GenerateTerrain returns a Size*Size heightmap in row-major order: [z*Size + x]
func GenerateTerrain(cfg TerrainConfig) []float32 {
	size := cfg.Size
	if size <= 0 {
		return nil
	}

	out := make([]float32, size*size)
	perm := GeneratePermutation(cfg.Seed)

	// compute max amplitude so FBM can be normalized similarly to your Zig
	var maxAmp float32 = 0
	{
		var a float32 = 1
		for i := uint8(0); i < cfg.Octaves; i++ {
			maxAmp += a
			a *= cfg.Persistence
		}
		if maxAmp == 0 {
			maxAmp = 1
		}
	}

	step := cfg.WorldStep
	if step == 0 {
		step = 1
	}

	for z := 0; z < size; z++ {
		wz := cfg.WorldStartZ + float32(z)*step
		for x := 0; x < size; x++ {
			wx := cfg.WorldStartX + float32(x)*step

			raw := FBM(wx, wz, cfg.Octaves, cfg.Persistence, cfg.Lacunarity, cfg.Scale, perm)
			normed := (raw / maxAmp) * cfg.NoiseWeight

			idx := z*size + x
			if cfg.BaseMap != nil && len(cfg.BaseMap) >= size*size {
				out[idx] = cfg.BaseMap[idx] + normed
			} else {
				out[idx] = normed
			}
		}
	}

	return out
}

func CreateBowlMap(size int) []float32 {
	if size <= 0 {
		return nil
	}
	out := make([]float32, size*size)

	center := float32(size) / 2
	maxDist := lmath.Sqrt(center * center * 2.0)

	for z := 0; z < size; z++ {
		for x := 0; x < size; x++ {
			dx := float32(x) - center
			dz := float32(z) - center
			dist := lmath.Sqrt(dx*dx + dz*dz)
			// pow((dist/maxDist), 2) + 0.2
			t := dist / maxDist
			out[z*size+x] = t*t + 0.2
		}
	}
	return out
}

// BilinearHeight samples within [0,size) x [0,size); returns ok=false if outside
func BilinearHeight(terrain []float32, x, z float32, size int) (h float32, ok bool) {
	if size <= 0 || len(terrain) < size*size {
		return 0, false
	}
	if x < 0 || z < 0 || x >= float32(size) || z >= float32(size) {
		return 0, false
	}

	x0 := float32(lmath.Floor(x))
	z0 := float32(lmath.Floor(z))
	tx := x - x0
	tz := z - z0

	ix0 := int(x0)
	iz0 := int(z0)
	ix1 := ix0 + 1
	iz1 := iz0 + 1
	if ix1 >= size {
		ix1 = size - 1
	}
	if iz1 >= size {
		iz1 = size - 1
	}

	h00 := terrain[iz0*size+ix0]
	h10 := terrain[iz0*size+ix1]
	h01 := terrain[iz1*size+ix0]
	h11 := terrain[iz1*size+ix1]

	row0 := h00 + tx*(h10-h00)
	row1 := h01 + tx*(h11-h01)
	return row0 + tz*(row1-row0), true
}

// --- Coloring / texture buffer ---
type Color = color.RGBA

func hsvToRGBA(h, s, v float32) Color {
	c := v * s
	x := c * (1 - lmath.Abs(lmath.Mod(h/60, 2)) - 1)
	m := v - c

	seg := lmath.Mod(h/60, 6)
	var r, g, b float32
	switch seg {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3:
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	toU8 := func(f float32) uint8 {
		if f <= 0 {
			return 0
		}
		if f >= 1 {
			return 255
		}
		return uint8(f * 255)
	}

	return Color{R: toU8(r + m), G: toU8(g + m), B: toU8(b + m), A: 255}
}

func terrainColor(height, waterNormed, sMod float32) Color {
	v := height - waterNormed

	var h, s, val float32
	switch {
	case v < -0.4:
		h, s, val = 40, 0.5, 0.0
	case v < 0.05:
		h, s, val = 40, 0.5, 0.8+v*2.0
	case v < 0.15:
		h, s, val = 95, 0.55+v, 0.625-v*0.5
	case v < 0.35:
		h, s, val = 110, 0.6+(v-0.15)*2.0, 0.5-(v-0.15)*0.7
	case v < 0.55:
		h, s, val = 30, 0.6-(v-0.35)*1.5, 0.3+(v-0.35)*1.5
	case v < 1.55:
		h, s, val = 210, 0.1, 1.55-v
	default:
		h, s, val = 210, 0.1, 0.0
	}

	s = lmath.Clamp(s-sMod, 0, 1)
	return hsvToRGBA(h, s, val)
}

// WriteTextureBuffer fills buffer (len >= w*h) where w = (baseSize-1)*textureScale.
func WriteTextureBuffer(buffer []Color, terrain []float32, baseSize int, textureScale float32, waterLevel float32, cubeHeight float32) (w, h int) {
	if baseSize <= 1 || textureScale <= 0 || cubeHeight == 0 {
		return 0, 0
	}
	w = int(float32(baseSize-1) * textureScale)
	h = w
	if w <= 0 || len(buffer) < w*h {
		return 0, 0
	}

	waterNormed := waterLevel / cubeHeight
	invScale := 1 / textureScale
	fBase := float32(baseSize - 1)

	const contrast float32 = 0.5
	const maxEffect float32 = 0.1

	for z := 0; z < h; z++ {
		fz := (float32(z) + 0.5) * invScale
		if fz > fBase {
			fz = fBase
		}
		for x := 0; x < w; x++ {
			fx := (float32(x) + 0.5) * invScale
			if fx > fBase {
				fx = fBase
			}

			hVal, _ := BilinearHeight(terrain, fx, fz, baseSize)

			hLeft := hVal
			if fx >= invScale {
				if v, ok := BilinearHeight(terrain, fx-invScale, fz, baseSize); ok {
					hLeft = v
				}
			}
			hUp := hVal
			if fz >= invScale {
				if v, ok := BilinearHeight(terrain, fx, fz-invScale, baseSize); ok {
					hUp = v
				}
			}

			dx := hVal - hLeft
			dz := hVal - hUp
			slopeRaw := dx + dz

			sMod := lmath.Clamp(slopeRaw*contrast, -maxEffect, maxEffect)
			buffer[z*w+x] = terrainColor(hVal, waterNormed, sMod)
		}
	}

	return w, h
}
