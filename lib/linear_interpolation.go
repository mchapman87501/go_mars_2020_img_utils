package lib

import (
	"sort"
)

type Float64Interpolator struct {
	// Mapping from known input values to known output values
	// Uses scaled input values
	yVals map[int]float64
	// Ordered keys from values:
	xVals []int
}

func toFixed(fval float64) int {
	return int(fval * 1000000)
}

func NewFloat64Interpolator(yVals map[float64]float64) *Float64Interpolator {
	xVals := make([]int, 0, len(yVals))
	myYVals := make(map[int]float64, len(yVals))
	for k, v := range yVals {
		fixed := toFixed(k)
		xVals = append(xVals, fixed)
		myYVals[fixed] = v
	}
	sort.Ints(xVals)

	return &Float64Interpolator{yVals: myYVals, xVals: xVals}
}

func (interp *Float64Interpolator) Interp(x float64) float64 {
	if len(interp.xVals) <= 0 {
		return x
	}
	fixed := toFixed(x)
	result, ok := interp.yVals[fixed]
	if ok {
		return result
	}
	xPrev, xNext := interp.bisectLeft(fixed)
	xOffset := fixed - xPrev
	fract := float64(xOffset) / float64(xNext-xPrev)
	yPrev := interp.yVals[xPrev]
	yNext := interp.yVals[xNext]

	// Cache the new result.
	result = yPrev + fract*(yNext-yPrev)
	interp.yVals[fixed] = result
	return result
}

func (interp *Float64Interpolator) bisectLeft(fixed int) (xPrev, xNext int) {
	// This is derived from Python's bisect.py.
	lo := 0
	hi := len(interp.xVals)
	for lo < hi {
		mid := int(lo+hi) / 2
		if interp.xVals[mid] > fixed {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	nextIndex := lo
	if nextIndex <= 0 {
		nextIndex = 1
	} else if nextIndex >= len(interp.xVals) {
		nextIndex = len(interp.xVals) - 1
	}
	xPrev = interp.xVals[nextIndex-1]
	xNext = interp.xVals[nextIndex]
	return
}
