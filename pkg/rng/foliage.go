package rng

import "math/rand"

type Point struct {
	X float32
	Z float32
}

type ForestConfig struct {
	MinDist           float32
	DensityScale      float32
	MoistureThreshold float32
}

// walk a grid in steps of MinDist, jitter inside each cell, then accept if height+moisture pass.
func GenerateFoliage(size int, terrain []float32, cfg ForestConfig, seed int64) []Point {
	if size <= 0 || len(terrain) < size*size {
		return nil
	}
	if cfg.MinDist <= 0 {
		return nil
	}
	if cfg.DensityScale == 0 {
		cfg.DensityScale = 1
	}

	points := make([]Point, 0, size) // rough starter cap

	perm := GeneratePermutation(seed + 1)
	r := rand.New(rand.NewSource(seed))

	fSize := float32(size)

	for y := float32(0); y < fSize; y += cfg.MinDist {
		for x := float32(0); x < fSize; x += cfg.MinDist {
			rx := x + r.Float32()*cfg.MinDist
			rz := y + r.Float32()*cfg.MinDist
			if rx >= fSize || rz >= fSize {
				continue
			}

			ix := int(rx)
			iz := int(rz)
			h := terrain[iz*size+ix]

			// same band as Zig
			if h > 0.15 && h < 0.7 {
				moist := noise(rx/cfg.DensityScale, rz/cfg.DensityScale, 0, perm)
				detail := noise(rx/10.0, rz/10.0, 1, perm) * 0.3
				nm := ((moist + detail) + 1.0) / 2.0

				if nm > cfg.MoistureThreshold {
					points = append(points, Point{X: rx, Z: rz})
				}
			}
		}
	}

	return points
}
