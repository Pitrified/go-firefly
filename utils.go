package firefly

import (
	"math"
	"math/rand"
)

// A message to be sent to the world when a firefly wants to change cell.
type ChangeCellReq struct {
	f        *Firefly
	from, to *Cell
}

// Cached values of cos/sin in [0,360) degrees.
var cCos, cSin map[int16]float32
var cacheCSdone = false

// Populate the cos/sin cache.
func cacheCosSin() {
	if cacheCSdone {
		return
	}
	cCos = make(map[int16]float32)
	cSin = make(map[int16]float32)
	for o := int16(0); o < 360; o++ {
		cCos[o] = float32(math.Cos(float64(o) * math.Pi / 180))
		cSin[o] = float32(math.Sin(float64(o) * math.Pi / 180))
	}
	cacheCSdone = true
}

// Returns a valid orientation in degrees in the [0, 360) interval.
func ValidateOri(o int16) int16 {
	for o < 0 {
		o += 360
	}
	for o >= 360 {
		o -= 360
	}
	return o
}

// Returns an int16 in the requested range, including extremes.
func RandRangeInt16(min, max int) int16 {
	return int16(rand.Intn(max+1-min) + min)
}

// Returns an int in the requested range, including extremes.
func RandRangeInt(min, max int) int {
	return rand.Intn(max+1-min) + min
}

// Absolute value for float32
func AbsFloat32(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}
