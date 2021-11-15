package main

import "github.com/Pitrified/go-firefly"

// MaxFloat32 returns the maximum value between the float32 parameters.
func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

// Get an arbitrary key from a map.
//
// https://stackoverflow.com/q/23482786/2237151
func getMapKey(m map[int]*firefly.Firefly) int {
	for k := range m {
		return k
	}
	return -1
}
