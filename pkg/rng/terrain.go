package rng

import (
	"math"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// --- Tuning knobs (safe defaults) ---
var (
	SeaLevel      = 0.30 // 0..1, raise to create more water
	SandBand      = 0.03 // 0..1, raise to create more sand
	MountainStart = 0.70 // 0..1, start of rocky
	SnowStart     = 0.86 // 0..1, start of snowy peaks
	WarpStrength  = 0.85 // domain warp magnitude
	WarpFrequency = 0.10 // domain warp noise frequency
	BaseFrequency = 0.20 // overall terrain scale (bigger => more detail)
)

// GenerateHeightmap returns a grayscale heightmap image (R=G=B).
func GenerateHeightmap(width, height int, seed int64) []raylib.Color {
	return GenerateHeightmapTiled(width, height, seed, 0, 0, 1)
}

func GenerateHeightmapTiledWorldSize(width, height int, seed int64, worldOriginX, worldOriginY, worldWidth float64) []raylib.Color {
	worldUnitsPerPixel := worldWidth / float64(width-1) // shared-edge convention
	return GenerateHeightmapTiled(width, height, seed, worldOriginX, worldOriginY, worldUnitsPerPixel)
}

// GenerateHeightmapTiled creates a grayscale heightmap (R=G=B) that is seamless across tiles.
//
// worldOriginX/Y and worldUnitsPerPixel define where this tile samples in WORLD space.
// Adjacent tiles will match perfectly if their origins line up with your mesh tiling convention.
func GenerateHeightmapTiled(width, height int, seed int64, worldOriginX, worldOriginY float64, worldUnitsPerPixel float64) []raylib.Color {
	colors := make([]raylib.Color, width*height)

	// Helpers
	clamp01 := func(x float64) float64 {
		if x < 0 {
			return 0
		}
		if x > 1 {
			return 1
		}
		return x
	}
	smoothstep := func(edge0, edge1, x float64) float64 {
		t := clamp01((x - edge0) / (edge1 - edge0))
		return t * t * (3 - 2*t)
	}

	// Band-limit based on sampling spacing to avoid aliasing when worldUnitsPerPixel grows.
	// Nyquist-ish max frequency in cycles/worldUnit:
	//   maxFreq ≈ 0.5 / worldUnitsPerPixel
	// Safety factor < 1 makes it smoother.
	maxFreq := 0.5 / worldUnitsPerPixel * 0.85

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			// TRUE world-space sample coords (no invMin normalization)
			wxWorld := worldOriginX + float64(x)*worldUnitsPerPixel
			wyWorld := worldOriginY + float64(y)*worldUnitsPerPixel

			// --- Domain warp (band-limited by using WarpFrequency in world units) ---
			wx := perlin2D(wxWorld*WarpFrequency, wyWorld*WarpFrequency, seed+101)
			wy := perlin2D(wxWorld*WarpFrequency, wyWorld*WarpFrequency, seed+202)

			// Warp offsets are dimensionless; scale by a world-space amount so it "means" something.
			// This keeps warp effect stable as you change BaseFrequency.
			warpWorldScale := 1.0 / BaseFrequency // ~ feature size in world units
			sxWorld := wxWorld + wx*WarpStrength*warpWorldScale*0.15
			syWorld := wyWorld + wy*WarpStrength*warpWorldScale*0.15

			// Convert world coords to noise coords using BaseFrequency (cycles/worldUnit)
			sx := sxWorld * BaseFrequency
			sy := syWorld * BaseFrequency

			// Compute how many octaves we are allowed, relative to each base multiplier.
			// Example: cont uses (base*0.20), so allowable octave multiplier is maxFreq/(base*0.20).
			allow := func(mult float64) float64 {
				f := BaseFrequency * mult
				if f <= 0 {
					return 0
				}
				return maxFreq / f
			}

			// --- Continentalness / large-scale landmass ---
			cont := fbmPerlin2D_BandLimited(sx*0.20, sy*0.20, seed+303, 8, 2.0, 0.5, allow(0.20)) // ~[-1,1]
			cont = cont*0.5 + 0.5                                                                 // -> [0,1]
			cont = math.Pow(cont, 1.35)

			// --- Rolling hills ---
			hills := fbmPerlin2D_BandLimited(sx*0.75, sy*0.75, seed+404, 10, 2.05, 0.52, allow(0.75)) // ~[-1,1]
			hills = hills*0.5 + 0.5
			hills = math.Pow(hills, 1.15)

			// --- Mountains / ridged noise ---
			mtn := ridgedFbm2D_BandLimited(sx*1.10, sy*1.10, seed+505, 9, 2.05, 0.55, allow(1.10)) // ~[0,1]
			mtn *= smoothstep(0.35, 0.85, cont)

			// --- Combine fields ---
			h := 0.55*cont + 0.35*hills + 0.55*mtn

			// --- Valley carving (mid-frequency) ---
			val := fbmPerlin2D_BandLimited(sx*1.60, sy*1.60, seed+606, 6, 2.0, 0.5, allow(1.60))
			val = val*0.5 + 0.5
			h -= 0.10 * val

			// Height curve
			h = clamp01(h)
			h = math.Pow(h, 1.30)

			// Sea shaping
			if h < SeaLevel {
				t := h / SeaLevel
				h = SeaLevel * math.Pow(t, 1.55)
			} else {
				t := (h - SeaLevel) / (1.0 - SeaLevel)
				h = SeaLevel + (1.0-SeaLevel)*math.Pow(t, 1.08)
			}

			h = clamp01(h)
			b := uint8(math.Round(h * 255.0))
			colors[y*width+x] = raylib.Color{R: b, G: b, B: b, A: 255}
		}
	}

	return colors
}

// ColorizeHeightmapTiled builds a colormap from a grayscale heightmap (R=G=B), but:
// - Seamless across tiles (uses world-space coords for variation).
// - Can output a higher-res texture (texW/texH) than the heightmap resolution.
//
// heights: grayscale colors from GenerateHeightmapTiled (len = hmW*hmH)
// worldOriginX/Y + worldUnitsPerPixel MUST match the ones used to generate the heightmap tile.
//
// Convention: assumes the heightmap spans world size:
//
//	worldWidth  = (hmW-1) * worldUnitsPerPixel
//	worldHeight = (hmH-1) * worldUnitsPerPixel
//
// and the output texture covers that same world area.
func ColorizeHeightmapTiled(
	heights []raylib.Color,
	hmW, hmH int,
	seed int64,
	worldOriginX, worldOriginY float64,
	worldWidth float64,
	texW, texH int,
) []raylib.Color {
	worldUnitsPerPixel := worldWidth / float64(hmW-1) // shared-edge convention
	out := make([]raylib.Color, texW*texH)

	clamp01 := func(x float64) float64 {
		if x < 0 {
			return 0
		}
		if x > 1 {
			return 1
		}
		return x
	}
	smoothstep := func(edge0, edge1, x float64) float64 {
		t := clamp01((x - edge0) / (edge1 - edge0))
		return t * t * (3 - 2*t)
	}
	lerp := func(a, b, t float64) float64 { return a + (b-a)*t }
	lerpColor := func(a, b raylib.Color, t float64) raylib.Color {
		t = clamp01(t)
		return raylib.Color{
			R: uint8(math.Round(lerp(float64(a.R), float64(b.R), t))),
			G: uint8(math.Round(lerp(float64(a.G), float64(b.G), t))),
			B: uint8(math.Round(lerp(float64(a.B), float64(b.B), t))),
			A: 255,
		}
	}

	// --- Palette (tweak as you like) ---
	deepWater := raylib.Color{R: 15, G: 40, B: 70, A: 255}
	shallowWater := raylib.Color{R: 25, G: 70, B: 100, A: 255}
	sand := raylib.Color{R: 194, G: 178, B: 128, A: 255}
	grassLow := raylib.Color{R: 45, G: 110, B: 55, A: 255}
	grassHigh := raylib.Color{R: 70, G: 140, B: 70, A: 255}
	rockLow := raylib.Color{R: 110, G: 110, B: 110, A: 255}
	rockHigh := raylib.Color{R: 170, G: 170, B: 170, A: 255}
	snowLow := raylib.Color{R: 205, G: 205, B: 205, A: 255}
	snowHigh := raylib.Color{R: 240, G: 240, B: 240, A: 255}

	// Output texture spans the same world area as the heightmap tile.
	worldW := float64(hmW-1) * worldUnitsPerPixel
	worldH := float64(hmH-1) * worldUnitsPerPixel

	// Avoid divide by zero if texW/texH == 1
	texUnitsPerPixelX := worldW
	texUnitsPerPixelY := worldH
	if texW > 1 {
		texUnitsPerPixelX = worldW / float64(texW-1)
	}
	if texH > 1 {
		texUnitsPerPixelY = worldH / float64(texH-1)
	}

	// Bilinear sampler for heights slice, returning [0..1]
	sampleHeight01 := func(fx, fy float64) float64 {
		// Clamp to valid range
		if fx < 0 {
			fx = 0
		}
		if fy < 0 {
			fy = 0
		}
		maxX := float64(hmW - 1)
		maxY := float64(hmH - 1)
		if fx > maxX {
			fx = maxX
		}
		if fy > maxY {
			fy = maxY
		}

		x0 := int(math.Floor(fx))
		y0 := int(math.Floor(fy))
		x1 := x0 + 1
		y1 := y0 + 1
		if x1 >= hmW {
			x1 = hmW - 1
		}
		if y1 >= hmH {
			y1 = hmH - 1
		}

		tx := fx - float64(x0)
		ty := fy - float64(y0)

		h00 := float64(heights[y0*hmW+x0].R) / 255.0
		h10 := float64(heights[y0*hmW+x1].R) / 255.0
		h01 := float64(heights[y1*hmW+x0].R) / 255.0
		h11 := float64(heights[y1*hmW+x1].R) / 255.0

		a := lerp(h00, h10, tx)
		b := lerp(h01, h11, tx)
		return lerp(a, b, ty)
	}

	// World-space variation (seamless across tiles)
	// Low frequencies so it reads like “texture detail” not noise.
	const (
		varFreq1 = 0.10 // cycles per world unit (adjust down if too speckly)
		varFreq2 = 0.04
	)

	for ty := 0; ty < texH; ty++ {
		for tx := 0; tx < texW; tx++ {

			// World-space coordinate for this texture pixel
			wx := worldOriginX + float64(tx)*texUnitsPerPixelX
			wy := worldOriginY + float64(ty)*texUnitsPerPixelY

			// Map world coord back into heightmap sample space
			hx := (wx - worldOriginX) / worldUnitsPerPixel // 0..hmW-1
			hy := (wy - worldOriginY) / worldUnitsPerPixel // 0..hmH-1

			h := sampleHeight01(hx, hy) // [0..1]

			// Variation noise: stable in world space, stable across tiles
			v1 := perlin2D(wx*varFreq1, wy*varFreq1, seed+9001) // ~[-1,1]
			v2 := perlin2D(wx*varFreq2, wy*varFreq2, seed+9002) // ~[-1,1]
			v := (0.6*v1 + 0.4*v2) * 0.06                       // roughly [-0.06..0.06]

			var c raylib.Color

			// Water
			if h <= SeaLevel {
				t := smoothstep(0.0, SeaLevel, h)
				c = lerpColor(deepWater, shallowWater, t)
			} else if h <= SeaLevel+SandBand {
				// Beach blend water->sand
				t := smoothstep(SeaLevel, SeaLevel+SandBand, h)
				c = lerpColor(shallowWater, sand, t)
			} else if h < MountainStart {
				// Grass zone
				t := smoothstep(SeaLevel+SandBand, MountainStart, h)
				c = lerpColor(grassLow, grassHigh, t)
			} else if h < SnowStart {
				// Rock zone
				t := smoothstep(MountainStart, SnowStart, h)
				c = lerpColor(rockLow, rockHigh, t)
			} else {
				// Snow zone
				t := smoothstep(SnowStart, 1.0, h)
				c = lerpColor(snowLow, snowHigh, t)
			}

			// Apply subtle brightness variation mostly to land (and lightly to sand)
			if h > SeaLevel*0.98 {
				r := clamp01(float64(c.R)/255.0 + v)
				g := clamp01(float64(c.G)/255.0 + v)
				b := clamp01(float64(c.B)/255.0 + v)
				c = raylib.Color{R: uint8(r * 255), G: uint8(g * 255), B: uint8(b * 255), A: 255}
			}

			out[ty*texW+tx] = c
		}
	}

	return out
}

func ColorizeHeightmap(heights []raylib.Color, width, height int, seed int64) []raylib.Color {
	return ColorizeHeightmapTiled(heights, height, width, seed, 0, 0, 1, height, width)
}

/* ------------------------- Band-limited FBM/Ridged ------------------------- */

// maxOctaveFreqMultiplier limits the internal fBm freq multiplier (starts at 1.0 and multiplies by lacunarity).
// Example: if maxOctaveFreqMultiplier == 4.0, octaves stop once freq > 4.
func fbmPerlin2D_BandLimited(x, y float64, seed int64, octaves int, lacunarity, gain float64, maxOctaveFreqMultiplier float64) float64 {
	amp := 1.0
	freq := 1.0
	sum := 0.0
	norm := 0.0

	for i := 0; i < octaves; i++ {
		if maxOctaveFreqMultiplier > 0 && freq > maxOctaveFreqMultiplier {
			break
		}
		sum += amp * perlin2D(x*freq, y*freq, seed+int64(i)*99991)
		norm += amp
		amp *= gain
		freq *= lacunarity
	}
	if norm > 0 {
		sum /= norm
	}
	return sum
}

func ridgedFbm2D_BandLimited(x, y float64, seed int64, octaves int, lacunarity, gain float64, maxOctaveFreqMultiplier float64) float64 {
	amp := 1.0
	freq := 1.0
	sum := 0.0
	norm := 0.0

	for i := 0; i < octaves; i++ {
		if maxOctaveFreqMultiplier > 0 && freq > maxOctaveFreqMultiplier {
			break
		}
		n := perlin2D(x*freq, y*freq, seed+int64(i)*131071) // [-1,1]
		n = 1.0 - math.Abs(n)                               // [0,1]
		n = n * n                                           // sharpen
		sum += n * amp
		norm += amp
		amp *= gain
		freq *= lacunarity
	}
	if norm > 0 {
		sum /= norm
	}
	if sum < 0 {
		return 0
	}
	if sum > 1 {
		return 1
	}
	return sum
}

/* ------------------------- Noise implementation (unchanged) ------------------------- */

// 2D Perlin noise in range approximately [-1, 1]
func perlin2D(x, y float64, seed int64) float64 {
	x0 := fastFloor(x)
	y0 := fastFloor(y)
	x1 := x0 + 1
	y1 := y0 + 1

	xf := x - float64(x0)
	yf := y - float64(y0)

	u := fade(xf)
	v := fade(yf)

	n00 := gradDot2D(x0, y0, x, y, seed)
	n10 := gradDot2D(x1, y0, x, y, seed)
	n01 := gradDot2D(x0, y1, x, y, seed)
	n11 := gradDot2D(x1, y1, x, y, seed)

	nx0 := lerp(n00, n10, u)
	nx1 := lerp(n01, n11, u)
	return lerp(nx0, nx1, v)
}

func gradDot2D(ix, iy int, x, y float64, seed int64) float64 {
	h := hash2D(ix, iy, seed)

	var gx, gy float64
	switch h & 7 {
	case 0:
		gx, gy = 1, 0
	case 1:
		gx, gy = -1, 0
	case 2:
		gx, gy = 0, 1
	case 3:
		gx, gy = 0, -1
	case 4:
		gx, gy = 0.70710678, 0.70710678
	case 5:
		gx, gy = -0.70710678, 0.70710678
	case 6:
		gx, gy = 0.70710678, -0.70710678
	default:
		gx, gy = -0.70710678, -0.70710678
	}

	dx := x - float64(ix)
	dy := y - float64(iy)
	return gx*dx + gy*dy
}

func hash2D(x, y int, seed int64) uint32 {
	h := uint32(seed) ^ 0x9E3779B9
	h ^= uint32(x) * 0x85EBCA6B
	h ^= uint32(y) * 0xC2B2AE35
	h ^= h >> 16
	h *= 0x7FEB352D
	h ^= h >> 15
	h *= 0x846CA68B
	h ^= h >> 16
	return h
}

func fade(t float64) float64       { return t * t * t * (t*(t*6-15) + 10) }
func lerp(a, b, t float64) float64 { return a + t*(b-a) }

func fastFloor(x float64) int {
	i := int(x)
	if float64(i) > x {
		return i - 1
	}
	return i
}
