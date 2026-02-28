package rng

import "math"

// GenerateCityHeightTile returns height values [0..255] (0 = ground).
// Deterministic + tileable like your terrain: sample is function of (worldX, worldY, seed).
func GenerateCityTile(
	width, height int,
	seed int64,
	worldOriginX, worldOriginY float64,
	worldUnitsPerPixel float64,
) []uint8 {

	out := make([]uint8, width*height)

	// --- Tunables (world units) ---
	const (
		blockSizeWorld   = 2000.0 // size of a city block
		streetWidthWorld = 3.0    // streets between blocks

		lotsPerBlock = 4
		lotInsetFrac = 0.14

		buildProb = 0.80

		// Height range in 0..255 space (kept low to avoid spikes)
		// Think of this as "meters" only after you scale it in your terrain mesh.
		minH = 18.0  // ~low buildings
		maxH = 120.0 // cap to avoid spires

		// Quantize to a small number of discrete heights (more city-like)
		heightLevels = 10 // e.g. 6..16
	)

	interior := blockSizeWorld - streetWidthWorld
	if interior <= 1.0 {
		panic("blockSizeWorld must be > streetWidthWorld + 1")
	}
	lotSize := interior / float64(lotsPerBlock)

	lerp := func(a, b, t float64) float64 { return a + (b-a)*t }

	hash01 := func(ix, iy int, salt int64) float64 {
		h := hash2D(ix, iy, seed+salt)
		return float64(h&0xFFFFFF) / float64(0xFFFFFF)
	}
	quantize := func(v float64, levels int) float64 {
		if levels <= 1 {
			return v
		}
		q := math.Round(v*float64(levels-1)) / float64(levels-1)
		return q
	}

	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			wx := worldOriginX + float64(px)*worldUnitsPerPixel
			wy := worldOriginY + float64(py)*worldUnitsPerPixel

			bx := int(math.Floor(wx / blockSizeWorld))
			by := int(math.Floor(wy / blockSizeWorld))

			lx := wx - float64(bx)*blockSizeWorld
			ly := wy - float64(by)*blockSizeWorld

			// Streets = ground
			if lx >= interior || ly >= interior {
				out[py*width+px] = 0
				continue
			}

			// Lot selection
			lotX := int(lx / lotSize)
			lotY := int(ly / lotSize)

			llx := lx - float64(lotX)*lotSize
			lly := ly - float64(lotY)*lotSize

			inset := lotSize * lotInsetFrac
			if llx < inset || llx > lotSize-inset || lly < inset || lly > lotSize-inset {
				out[py*width+px] = 0
				continue
			}

			// Deterministic lot key
			lotKeyX := bx*1000 + lotX
			lotKeyY := by*1000 + lotY

			if hash01(lotKeyX, lotKeyY, 1001) > buildProb {
				out[py*width+px] = 0
				continue
			}

			// Building footprint
			fw := lerp(0.60, 0.92, hash01(lotKeyX, lotKeyY, 1002))
			fh := lerp(0.60, 0.92, hash01(lotKeyX, lotKeyY, 1003))

			innerW := lotSize - 2*inset
			innerH := lotSize - 2*inset
			bW := innerW * fw
			bH := innerH * fh

			offX := (hash01(lotKeyX, lotKeyY, 1004) - 0.5) * (innerW - bW) * 0.25
			offY := (hash01(lotKeyX, lotKeyY, 1005) - 0.5) * (innerH - bH) * 0.25

			cx := inset + innerW*0.5 + offX
			cy := inset + innerH*0.5 + offY

			x0 := cx - bW*0.5
			x1 := cx + bW*0.5
			y0 := cy - bH*0.5
			y1 := cy + bH*0.5

			if llx < x0 || llx > x1 || lly < y0 || lly > y1 {
				out[py*width+px] = 0
				continue
			}

			// Building height: biased towards lower heights, then quantized
			r := hash01(lotKeyX, lotKeyY, 1006)
			r = math.Pow(r, 1.9) // more low buildings than tall

			h01 := quantize(r, heightLevels)
			hv := lerp(minH, maxH, h01)

			// Clamp + write
			if hv < 0 {
				hv = 0
			}
			if hv > 255 {
				hv = 255
			}
			out[py*width+px] = uint8(math.Round(hv))
		}
	}

	return out
}
