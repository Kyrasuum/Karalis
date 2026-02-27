package rng

import (
	"math"
	"sort"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	TileWall  = uint8(0)
	TileFloor = uint8(1)
	TileDoor  = uint8(2)
)

type room struct {
	id     int64
	cx, cy int // center (world cell coords)
	w, h   int // size in cells
	x0, y0 int // bounds inclusive
	x1, y1 int
}

// GenerateDungeonTile generates a roguelike-style catacombs tile.
// - Deterministic by (seed, world cell coords)
// - Tileable via expanded margin generation + crop
//
// width,height: tile size in CELLS
// originX,originY: global cell coordinate for local (0,0)
// You probably want odd sizes like 129, 257, etc.
func GenerateDungeonTile(width, height int, seed int64, originX, originY int) []uint8 {
	grid := make([]uint8, width*height)
	for i := range grid {
		grid[i] = TileWall
	}

	// ----- Tunables -----
	const (
		margin   = 32   // expanded gen margin (cells). bigger = safer seams
		roomStep = 18   // lattice spacing for room candidates (bigger => fewer rooms)
		roomProb = 0.60 // chance a candidate becomes a room
		roomMin  = 7    // room min size
		roomMax  = 13   // room max size
		keepout  = 2    // room padding to reduce overlaps

		neighborsK        = 3    // connect each room to K nearest rooms
		extraLoopProb     = 0.25 // extra corridor connections beyond nearest
		corridorHalfWidth = 1    // 1 => corridors ~3 tiles wide with carving

		doorChance = 0.85 // place a door when corridor enters a room
	)

	// Expanded region in WORLD coords
	// We generate rooms/corridors for a bigger area, then crop to the tile.
	ex0 := originX - margin
	ey0 := originY - margin
	ex1 := originX + (width - 1) + margin
	ey1 := originY + (height - 1) + margin

	// Helpers
	idx := func(x, y int) int { return y*width + x }
	inTile := func(wx, wy int) bool {
		return wx >= originX && wx <= originX+width-1 && wy >= originY && wy <= originY+height-1
	}
	toLocal := func(wx, wy int) (lx, ly int) { return wx - originX, wy - originY }
	setLocal := func(lx, ly int, t uint8) {
		if lx < 0 || lx >= width || ly < 0 || ly >= height {
			return
		}
		i := idx(lx, ly)
		// doors should overwrite floor, but not vice versa
		if t == TileDoor {
			grid[i] = TileDoor
			return
		}
		if grid[i] != TileDoor {
			grid[i] = t
		}
	}
	// Carve in world coords (but only writes if inside tile crop)
	carveWorld := func(wx, wy int, t uint8) {
		if !inTile(wx, wy) {
			return
		}
		lx, ly := toLocal(wx, wy)
		setLocal(lx, ly, t)
	}
	clamp := func(v, lo, hi int) int {
		if v < lo {
			return lo
		}
		if v > hi {
			return hi
		}
		return v
	}
	hash01 := func(x, y int, salt int64) float64 {
		h := hash2D(x, y, seed+salt)
		return float64(h&0xFFFFFF) / float64(0xFFFFFF)
	}
	hashI := func(x, y int, salt int64, lo, hi int) int {
		if hi <= lo {
			return lo
		}
		r := hash01(x, y, salt)
		return lo + int(math.Floor(r*float64(hi-lo+1)))
	}

	// ----- 1) Place rooms deterministically in expanded region -----
	rooms := make([]room, 0, 128)
	occupied := make(map[[2]int]struct{}, 4096) // coarse occupancy to reduce overlaps (fast + deterministic)

	// simple coarse occupancy key
	occKey := func(x, y int) [2]int { return [2]int{x >> 2, y >> 2} } // 4x4 buckets

	for cy := ey0; cy <= ey1; cy += roomStep {
		for cx := ex0; cx <= ex1; cx += roomStep {
			// jitter candidate center
			jx := int(math.Round((hash01(cx, cy, 1001) - 0.5) * float64(roomStep) * 0.6))
			jy := int(math.Round((hash01(cx, cy, 1002) - 0.5) * float64(roomStep) * 0.6))
			rcx := cx + jx
			rcy := cy + jy

			if hash01(rcx, rcy, 1003) > roomProb {
				continue
			}

			w := hashI(rcx, rcy, 1004, roomMin, roomMax)
			h := hashI(rcx, rcy, 1005, roomMin, roomMax)
			// force odd sizes for nicer symmetry
			if w%2 == 0 {
				w++
			}
			if h%2 == 0 {
				h++
			}

			x0 := rcx - w/2
			y0 := rcy - h/2
			x1 := rcx + w/2
			y1 := rcy + h/2

			// keep rooms away from expanded region edge a bit
			if x0 < ex0+2 || y0 < ey0+2 || x1 > ex1-2 || y1 > ey1-2 {
				continue
			}

			// overlap rejection via coarse buckets
			ok := true
			for y := y0 - keepout; y <= y1+keepout && ok; y += 2 {
				for x := x0 - keepout; x <= x1+keepout; x += 2 {
					k := occKey(x, y)
					if _, exists := occupied[k]; exists {
						ok = false
						break
					}
				}
			}
			if !ok {
				continue
			}

			id := int64(hash2D(rcx, rcy, seed+4242))
			r := room{id: id, cx: rcx, cy: rcy, w: w, h: h, x0: x0, y0: y0, x1: x1, y1: y1}
			rooms = append(rooms, r)

			// mark occupancy
			for y := y0 - keepout; y <= y1+keepout; y += 2 {
				for x := x0 - keepout; x <= x1+keepout; x += 2 {
					occupied[occKey(x, y)] = struct{}{}
				}
			}

			// carve room floor (cropped to tile)
			for y := y0; y <= y1; y++ {
				for x := x0; x <= x1; x++ {
					carveWorld(x, y, TileFloor)
				}
			}
		}
	}

	if len(rooms) == 0 {
		return grid // all walls
	}

	// ----- 2) Build deterministic connections (nearest neighbors + a light backbone) -----
	type edge struct {
		a, b int // indices into rooms
		d2   int
		key  uint64
	}

	edges := make([]edge, 0, len(rooms)*neighborsK*2)

	// For each room, connect to K nearest rooms (in expanded region)
	for i := range rooms {
		type cand struct{ j, d2 int }
		cands := make([]cand, 0, len(rooms)-1)
		for j := range rooms {
			if i == j {
				continue
			}
			dx := rooms[i].cx - rooms[j].cx
			dy := rooms[i].cy - rooms[j].cy
			d2 := dx*dx + dy*dy
			cands = append(cands, cand{j: j, d2: d2})
		}
		sort.Slice(cands, func(a, b int) bool { return cands[a].d2 < cands[b].d2 })
		k := neighborsK
		if k > len(cands) {
			k = len(cands)
		}
		for n := 0; n < k; n++ {
			j := cands[n].j
			a, b := i, j
			if a > b {
				a, b = b, a
			}
			key := (uint64(a) << 32) | uint64(b)
			edges = append(edges, edge{a: a, b: b, d2: cands[n].d2, key: key})
		}
	}

	// Dedup edges
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].key == edges[j].key {
			return edges[i].d2 < edges[j].d2
		}
		return edges[i].key < edges[j].key
	})
	dedup := edges[:0]
	var last uint64 = ^uint64(0)
	for _, e := range edges {
		if e.key != last {
			dedup = append(dedup, e)
			last = e.key
		}
	}
	edges = dedup

	// Optionally add extra loops stochastically (still deterministic)
	loopEdges := make([]edge, 0, len(edges))
	for _, e := range edges {
		// prefer adding loops for longer edges
		r := hash01(rooms[e.a].cx^rooms[e.b].cx, rooms[e.a].cy^rooms[e.b].cy, 6001)
		if r < extraLoopProb {
			loopEdges = append(loopEdges, e)
		}
	}

	// A minimal backbone: sort by distance and take some edges (MST-ish without full DSU)
	sort.Slice(edges, func(i, j int) bool { return edges[i].d2 < edges[j].d2 })

	// DSU to build a connected backbone graph (this removes “blob randomness”)
	parent := make([]int, len(rooms))
	rank := make([]int, len(rooms))
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	union := func(a, b int) bool {
		ra, rb := find(a), find(b)
		if ra == rb {
			return false
		}
		if rank[ra] < rank[rb] {
			ra, rb = rb, ra
		}
		parent[rb] = ra
		if rank[ra] == rank[rb] {
			rank[ra]++
		}
		return true
	}

	backbone := make([]edge, 0, len(rooms)-1)
	for _, e := range edges {
		if union(e.a, e.b) {
			backbone = append(backbone, e)
			if len(backbone) == len(rooms)-1 {
				break
			}
		}
	}

	// Final connections = backbone + some extra loops
	connections := append([]edge{}, backbone...)
	connections = append(connections, loopEdges...)

	// ----- 3) Carve corridors and doors -----
	// Corridors are carved in the expanded region logically, but only written to the tile crop.
	carveRect := func(x0, y0, x1, y1 int, t uint8) {
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		for y := y0; y <= y1; y++ {
			for x := x0; x <= x1; x++ {
				carveWorld(x, y, t)
			}
		}
	}

	// Returns whether point is inside a room
	pointInRoom := func(r room, x, y int) bool {
		return x >= r.x0 && x <= r.x1 && y >= r.y0 && y <= r.y1
	}

	// Find a “door point” on the room boundary near a target point
	doorOnRoom := func(r room, tx, ty int) (dx, dy int) {
		// clamp target to room bounds, then push to nearest edge
		x := clamp(tx, r.x0, r.x1)
		y := clamp(ty, r.y0, r.y1)
		// choose nearest edge
		left := x - r.x0
		right := r.x1 - x
		top := y - r.y0
		bot := r.y1 - y
		min := left
		edge := 0
		if right < min {
			min = right
			edge = 1
		}
		if top < min {
			min = top
			edge = 2
		}
		if bot < min {
			min = bot
			edge = 3
		}
		switch edge {
		case 0:
			x = r.x0
		case 1:
			x = r.x1
		case 2:
			y = r.y0
		case 3:
			y = r.y1
		}
		return x, y
	}

	for _, e := range connections {
		a := rooms[e.a]
		b := rooms[e.b]

		// Corridor endpoints: door points on each room boundary
		adx, ady := doorOnRoom(a, b.cx, b.cy)
		bdx, bdy := doorOnRoom(b, a.cx, a.cy)

		// Decide bend order deterministically to avoid global diagonal bias
		// (This is one big source of “one diagonal line” artifacts.)
		bendFirstX := hash01(adx^bdx, ady^bdy, 7001) < 0.5

		// Carve corridor with thickness (3-wide by default)
		carveLine := func(x0, y0, x1, y1 int) {
			if x0 == x1 {
				if y0 > y1 {
					y0, y1 = y1, y0
				}
				for y := y0; y <= y1; y++ {
					carveRect(x0-corridorHalfWidth, y, x0+corridorHalfWidth, y, TileFloor)
				}
				return
			}
			if y0 == y1 {
				if x0 > x1 {
					x0, x1 = x1, x0
				}
				for x := x0; x <= x1; x++ {
					carveRect(x, y0-corridorHalfWidth, x, y0+corridorHalfWidth, TileFloor)
				}
				return
			}
		}

		// L corridor with randomized order, plus slight deterministic “jitter” elbow sometimes
		elbowX, elbowY := bdx, ady
		if !bendFirstX {
			elbowX, elbowY = adx, bdy
		}
		// elbow jitter (prevents long repeated diagonals)
		if hash01(elbowX, elbowY, 7002) < 0.35 {
			if bendFirstX {
				elbowY += int(math.Round((hash01(elbowX, elbowY, 7003) - 0.5) * 6.0))
			} else {
				elbowX += int(math.Round((hash01(elbowX, elbowY, 7004) - 0.5) * 6.0))
			}
		}

		// Carve: A door -> elbow -> B door
		// segments are axis-aligned; jitter might make elbow not aligned, so we do 2-step clamp
		carveLine(adx, ady, elbowX, ady)
		carveLine(elbowX, ady, elbowX, elbowY)
		carveLine(elbowX, elbowY, bdx, elbowY)
		carveLine(bdx, elbowY, bdx, bdy)

		// Place doors at boundary points (optional)
		if hash01(adx, ady, 8001) < doorChance {
			carveWorld(adx, ady, TileDoor)
		}
		if hash01(bdx, bdy, 8002) < doorChance {
			carveWorld(bdx, bdy, TileDoor)
		}

		// Ensure corridor doesn't accidentally “fill” room interiors too aggressively:
		// (We already carved room floors; doors sit on boundary.)
		_ = pointInRoom // kept for extension if you want to enforce boundary-only doors
	}

	return grid
}

// GenerateCaveTile returns a Cave tile as a binary grid (0=wall, 1=floor).
// It is seamless across tiles when you use worldOriginCellX/Y that line up on a global grid.
//
// width,height: number of CELLS in this tile (e.g., 128x128)
// seed: deterministic seed
// worldOriginCellX/Y: the global cell coordinate of local (0,0) for this tile
//
// Example tiling:
//
//	tile (tx,ty) origin = (tx*(width-1), ty*(height-1)) if you share borders,
//	or (tx*width, ty*height) if you don't.
func GenerateCaveTile(width, height int, seed int64, worldOriginCellX, worldOriginCellY int) []uint8 {
	out := make([]uint8, width*height)

	// --- Tunables (start here) ---
	// Bigger -> wider corridors / more open space.
	const (
		// Base corridor “band” thickness
		corridorWidth = 0.18 // 0..1-ish, higher => more corridors

		// Room frequency/threshold controls
		roomFreq = 0.035 // rooms blob size (lower => bigger rooms)
		roomCut  = 0.62  // higher => fewer rooms

		// Corridor field frequencies
		corrFreqA = 0.060 // main corridors
		corrFreqB = 0.095 // secondary to increase branching

		// Domain warp to break grid patterns
		warpFreq     = 0.040
		warpStrength = 1.25
	)

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
	smoothstep := func(e0, e1, x float64) float64 {
		t := clamp01((x - e0) / (e1 - e0))
		return t * t * (3 - 2*t)
	}

	// Pass 1: field thresholding into floor/wall (deterministic per world cell)
	for y := 0; y < height; y++ {
		gy := worldOriginCellY + y
		for x := 0; x < width; x++ {
			gx := worldOriginCellX + x

			// World-space continuous coords (cells as units)
			fx := float64(gx)
			fy := float64(gy)

			// Domain warp (two low-freq fields)
			wx := fbmValue2D(fx*warpFreq, fy*warpFreq, seed+101, 3, 2.0, 0.5) // ~[-1,1]
			wy := fbmValue2D(fx*warpFreq, fy*warpFreq, seed+202, 3, 2.0, 0.5)
			fxw := fx + wx*warpStrength
			fyw := fy + wy*warpStrength

			// Corridor fields: “walkable where noise is near 0” (band-pass)
			// Using absolute value gives maze-ish lines that branch when layered.
			cA := fbmValue2D(fxw*corrFreqA, fyw*corrFreqA, seed+303, 5, 2.0, 0.5) // [-1,1]
			cB := fbmValue2D(fxw*corrFreqB, fyw*corrFreqB, seed+404, 4, 2.0, 0.5)

			// Convert to “corridor-ness” [0..1]
			// High when |noise| is small.
			corrA := 1.0 - clamp01(math.Abs(cA)/corridorWidth)
			corrB := 1.0 - clamp01(math.Abs(cB)/(corridorWidth*0.85))

			// Rooms: big blobs from low-freq noise
			rn := fbmValue2D(fx*roomFreq, fy*roomFreq, seed+505, 3, 2.0, 0.5) // [-1,1]
			roomMask := smoothstep(roomCut, roomCut+0.10, rn*0.5+0.5)

			// Add a little “worm” carving guided by a flow direction (more organic branching)
			// This is still deterministic because it’s derived from world coords.
			flow := fbmValue2D(fx*0.030, fy*0.030, seed+606, 2, 2.0, 0.5)
			worm := 1.0 - clamp01(math.Abs(flow)/(corridorWidth*0.9))

			// Combine: corridors + rooms
			floorScore := math.Max(math.Max(corrA, corrB*0.9), worm*0.65)
			floorScore = math.Max(floorScore, roomMask)

			// Threshold (tweak if you want more/less walls)
			if floorScore > 0.22 {
				out[y*width+x] = TileFloor
			} else {
				out[y*width+x] = TileWall
			}
		}
	}

	// Pass 2: small deterministic cleanup for nicer corridors.
	// A couple of cellular-automata style smoothing steps improve readability.
	// IMPORTANT: We must use ONLY local neighborhood data, and since both tiles compute the same border cells,
	// seams still match as long as tiles overlap borders consistently (shared-edge convention).
	out = smoothBinary(out, width, height, 2)

	// Pass 3: optional guaranteed walkable “backbone” roads across the tile to improve cross-tile connectivity.
	// These are deterministic world-space lines that help ensure you don’t get isolated islands per tile.
	carveBackbone(out, width, height, seed, worldOriginCellX, worldOriginCellY)

	return out
}

// Optional helper to visualize as a grayscale heightmap image (wall=0 floor=255).
func ToGrayscaleColors(grid []uint8, width, height int) []raylib.Color {
	cols := make([]raylib.Color, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v := uint8(100)
			if grid[y*width+x] == TileFloor {
				v = 200
			}
			cols[y*width+x] = raylib.Color{R: v, G: v, B: v, A: 255}
		}
	}
	return cols
}

/* ------------------- Smoothing / connectivity helpers ------------------- */

func smoothBinary(src []uint8, w, h, steps int) []uint8 {
	dst := make([]uint8, len(src))
	copy(dst, src)

	// Rule: if a wall has many floor neighbors, open it; if a floor has many wall neighbors, close it.
	// This is mild smoothing, not full cave-generation.
	for s := 0; s < steps; s++ {
		tmp := make([]uint8, len(dst))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				floors := 0
				for oy := -1; oy <= 1; oy++ {
					yy := y + oy
					if yy < 0 || yy >= h {
						continue
					}
					for ox := -1; ox <= 1; ox++ {
						xx := x + ox
						if xx < 0 || xx >= w || (ox == 0 && oy == 0) {
							continue
						}
						if dst[yy*w+xx] == TileFloor {
							floors++
						}
					}
				}
				cur := dst[y*w+x]
				if cur == TileWall {
					// Open narrow choke points a bit
					if floors >= 6 {
						tmp[y*w+x] = TileFloor
					} else {
						tmp[y*w+x] = TileWall
					}
				} else {
					// Close tiny specks
					if floors <= 2 {
						tmp[y*w+x] = TileWall
					} else {
						tmp[y*w+x] = TileFloor
					}
				}
			}
		}
		dst = tmp
	}
	return dst
}

// Deterministic “roads” that run across a tile to improve cross-tile connectivity.
// These are based on world cell coords so adjacent tiles line up.
func carveBackbone(grid []uint8, w, h int, seed int64, originX, originY int) {
	// Carve a couple of meandering lines across the tile.
	// If you don’t want this, remove this call.
	for y := 0; y < h; y++ {
		gy := originY + y
		// A wavy x centerline
		n := fbmValue2D(float64(originX)*0.02, float64(gy)*0.02, seed+777, 3, 2.0, 0.5) // [-1,1]
		center := int(float64(w)/2 + n*float64(w)*0.18)
		for dx := -1; dx <= 1; dx++ {
			x := center + dx
			if x >= 0 && x < w {
				grid[y*w+x] = TileFloor
			}
		}
	}
	for x := 0; x < w; x++ {
		gx := originX + x
		// A wavy y centerline
		n := fbmValue2D(float64(gx)*0.02, float64(originY)*0.02, seed+888, 3, 2.0, 0.5) // [-1,1]
		center := int(float64(h)/2 + n*float64(h)*0.18)
		for dy := -1; dy <= 1; dy++ {
			y := center + dy
			if y >= 0 && y < h {
				grid[y*w+x] = TileFloor
			}
		}
	}
}

/* ----------------------- Deterministic value noise ---------------------- */

// Value noise (smooth grid interpolation) in [-1,1]
func valueNoise2D(x, y float64, seed int64) float64 {
	x0 := fastFloor(x)
	y0 := fastFloor(y)
	x1 := x0 + 1
	y1 := y0 + 1

	xf := x - float64(x0)
	yf := y - float64(y0)

	u := fade(xf)
	v := fade(yf)

	// Corner values in [-1,1]
	v00 := hashVal2D(x0, y0, seed)
	v10 := hashVal2D(x1, y0, seed)
	v01 := hashVal2D(x0, y1, seed)
	v11 := hashVal2D(x1, y1, seed)

	a := lerp(v00, v10, u)
	b := lerp(v01, v11, u)
	return lerp(a, b, v)
}

func fbmValue2D(x, y float64, seed int64, octaves int, lacunarity, gain float64) float64 {
	amp := 1.0
	freq := 1.0
	sum := 0.0
	norm := 0.0

	for i := 0; i < octaves; i++ {
		sum += amp * valueNoise2D(x*freq, y*freq, seed+int64(i)*99991)
		norm += amp
		amp *= gain
		freq *= lacunarity
	}
	if norm > 0 {
		sum /= norm
	}
	return sum
}

// Hash cell -> deterministic float in [-1,1]
func hashVal2D(x, y int, seed int64) float64 {
	h := hash2D(x, y, seed)
	// 0..1
	f := float64(h&0xFFFFFF) / float64(0xFFFFFF)
	return f*2.0 - 1.0
}
