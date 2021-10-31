package xrand

import (
	"math"
	"math/bits"
)

// https://prng.di.unimi.it/xoshiro256plus.c

type Xoshiro256p struct {
	s [4]uint64
}

func NewXoshiro256p(seed int64) *Xoshiro256p {
	x := Xoshiro256p{}
	x.Seed(seed)
	return &x
}

func (x *Xoshiro256p) Seed(seed int64) {
	s := NewSplitMix64(seed)
	x.s[0] = s.Uint64()
	x.s[1] = s.Uint64()
	x.s[2] = s.Uint64()
	x.s[3] = s.Uint64()
}

func (x *Xoshiro256p) Uint64() uint64 {
	result := x.s[0] + x.s[3]
	t := x.s[1] << 17
	x.s[2] ^= x.s[0]
	x.s[3] ^= x.s[1]
	x.s[1] ^= x.s[2]
	x.s[0] ^= x.s[3]

	x.s[2] ^= t
	x.s[3] = bits.RotateLeft64(x.s[3], 45)

	return result
}

func (x *Xoshiro256p) Int64() int64 {
	return unsafeUint64ToInt64(x.Uint64())
}

func (x *Xoshiro256p) Int63() int64 {
	return int64(x.Uint64() & (1<<63 - 1))
}

func (x *Xoshiro256p) Float64() float64 {
	return math.Float64frombits(0x3ff<<52|x.Uint64()>>12) - 1.0
}

// Call Uint64() * 2^128
func (x *Xoshiro256p) Jump() {
	var (
		jump = [...]uint64{
			0x180ec6d33cfd0aba, 0xd5a61266f0c9392c, 0xa9582618e03fc9aa, 0x39abdc4529b1661c,
		}
	)

	s := [...]uint64{0, 0, 0, 0}
	for i := range jump {
		for b := 0; b < 64; b++ {
			if (jump[i] & 1 << b) != 0 {
				s[0] ^= x.s[0]
				s[1] ^= x.s[1]
				s[2] ^= x.s[2]
				s[3] ^= x.s[3]
			}
		}
		x.Uint64()
	}

	x.s[0] = s[0]
	x.s[1] = s[1]
	x.s[2] = s[2]
	x.s[3] = s[3]
}

// Call Uint64() * 2^192
func (x *Xoshiro256p) LongJump() {
	var (
		jump = [...]uint64{
			0x76e15d3efefdcbbf, 0xc5004e441c522fb3, 0x77710069854ee241, 0x39109bb02acbe635,
		}
	)

	s := [...]uint64{0, 0, 0, 0}
	for i := range jump {
		for b := 0; b < 64; b++ {
			if (jump[i] & 1 << b) != 0 {
				s[0] ^= x.s[0]
				s[1] ^= x.s[1]
				s[2] ^= x.s[2]
				s[3] ^= x.s[3]
			}
		}
		x.Uint64()
	}

	x.s[0] = s[0]
	x.s[1] = s[1]
	x.s[2] = s[2]
	x.s[3] = s[3]
}
