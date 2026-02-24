package rng

import (
	"errors"
	"math"
	"math/rand"
)

type PityRng struct {
	r *rand.Rand

	chance   float64
	ramp     float64
	failures int
}

// --- Calibration ---
//
// We choose C so that the long-run proc rate equals chance.
// In steady-state, proc rate = 1 / E[rolls per success].
// We compute E as:
//
//	P(success on i) = (prod_{j< i} (1 - chance(j))) * chance(i)
//	E = sum_i i * P(success on i)
//
// Then we binary-search C in (0,1].
func calibrate(chance float64) float64 {
	// Monotonic: larger C => more frequent success => higher proc rate.
	lo, hi := 0.0, 1.0

	// Binary search tolerance: good enough for gameplay systems.
	for iter := 0; iter < 80; iter++ {
		mid := (lo + hi) / 2
		rate := procRate(mid)
		if rate < chance {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

func procRate(ramp float64) float64 {
	if ramp <= 0 {
		return 0
	}

	// Determine the maximum attempt index we need to consider.
	// Attempts are 1-based: i=1 uses failures=0.
	// First i where C*i >= 1 => chance hits 1
	limit := int(math.Ceil(1.0 / ramp))
	if limit < 1 {
		limit = 1
	}

	survive := 1.0 // probability we have not succeeded before attempt i
	E := 0.0

	for i := 1; i <= limit; i++ {
		ch := ramp * float64(i)
		if ch > 1 {
			ch = 1
		}
		if i == limit {
			ch = 1 // explicit force
		}

		pSuccessAtI := survive * ch
		E += float64(i) * pSuccessAtI
		survive *= (1 - ch)

		// Once ch=1, survive becomes 0 and remaining terms are 0.
		if ch >= 1 {
			break
		}
	}

	// proc rate = successes / rolls = 1 / E[rolls per success]
	if E <= 0 {
		return 0
	}
	return 1.0 / E
}

func NewPityRng(chance float64) (*PityRng, error) {
	if chance < 0 || chance > 1 || math.IsNaN(chance) {
		return nil, errors.New("chance must be in [0,1]")
	}
	seed := rand.Int63()

	r := rand.New(rand.NewSource(seed))

	prd := &PityRng{
		r:      r,
		chance: chance,
	}

	switch {
	case chance == 0:
		prd.ramp = 0
	case chance == 1:
		prd.ramp = 1
	default:
		prd.ramp = calibrate(chance)
	}

	return prd, nil
}

func (p *PityRng) Next() bool {
	// Fast-path edges
	if p.chance <= 0 {
		p.failures++
		return false
	}
	if p.chance >= 1 {
		p.failures = 0
		return true
	}

	ch := p.ChanceNow()
	if p.r.Float64() < ch {
		p.failures = 0
		return true
	}

	p.failures++
	return false
}

// ChanceNow returns the current proc chance for the *next* roll (based on current failures).
func (p *PityRng) ChanceNow() float64 {
	// Linear ramp
	ch := p.ramp * float64(p.failures+1)
	if ch > 1 {
		return 1
	}
	if ch < 0 {
		return 0
	}
	return ch
}

func (p *PityRng) Failures() int {
	return p.failures
}

func (p *PityRng) Reset() {
	p.failures = 0
}

func (p *PityRng) Reseed(seed int64) {
	p.r.Seed(seed)
	p.failures = 0
}
