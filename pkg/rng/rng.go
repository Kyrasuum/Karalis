package rng

import (
	"errors"
	"math"
	"math/rand"

	"karalis/pkg/lmath"
)

type PityRng struct {
	r *rand.Rand

	target float64
	fails  float64
	offset int
}

func NewPityRng(chance float64) (*PityRng, error) {
	if chance < 0 || math.IsNaN(chance) {
		return nil, errors.New("chance must be in [0,]")
	}
	seed := rand.Int63()
	r := rand.New(rand.NewSource(seed))

	offset := lmath.Floor(chance)

	prd := &PityRng{
		r:      r,
		target: chance - float64(offset),
		fails:  0,
		offset: offset,
	}

	return prd, nil
}

func (p *PityRng) Next() int {
	if p == nil {
		return 0
	}
	roll := p.r.Float64()
	chance := p.ChanceNow()
	p.fails++
	if roll < chance {
		if chance < 0.5 {
			expected := float64(lmath.Ceil(1 / p.target))
			p.fails -= expected
		} else {
			expected := 1.0 / p.target
			p.fails -= expected
		}
		return 1 + p.offset
	}
	return p.offset
}

func (p *PityRng) ChanceNow() float64 {
	if p == nil {
		return 0.0
	}
	if p.target == 0 {
		return 0.0
	}
	if p.target == 1 {
		return 1
	}

	expected := float64(lmath.Ceil(1 / p.target))
	if p.fails < 0 {
		return p.target * (expected - lmath.Abs(p.fails)) / expected
	}

	// return p.target * float64(p.fails+1)
	return 1 / (expected - p.fails)
}

func (p *PityRng) Reset() {
	if p == nil {
		return
	}
	p.fails = 0
}

func (p *PityRng) Reseed(seed int64) {
	if p == nil {
		return
	}
	p.r.Seed(seed)
	p.Reset()
}
