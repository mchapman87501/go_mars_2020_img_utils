package lib

import (
	"sort"
)

type Float64Interpolator struct {
	// Mapping from known input values to known output values
	yVals map[float64]float64
	// Ordered keys from values:
	xVals []float64
}

func NewFloat64Interpolator(yVals map[float64]float64) *Float64Interpolator {
	var xVals []float64 = make([]float64, 0, len(yVals))
	for k := range yVals {
		xVals = append(xVals, k)
	}
	sort.Float64s(xVals)

	return &Float64Interpolator{yVals: yVals, xVals: xVals}
}

func (interp *Float64Interpolator) Interp(x float64) float64 {
	xPrev, xNext := interp.bisectLeft(x)
	xOffset := x - xPrev
	fract := xOffset / (xNext - xPrev)
	yPrev := interp.yVals[xPrev]
	yNext := interp.yVals[xNext]
	return yPrev + fract*(yNext-yPrev)
}

func (interp *Float64Interpolator) bisectLeft(x float64) (xPrev, xNext float64) {
	// This is derived from Python's bisect.py.
	lo := 0
	hi := len(interp.xVals)
	for lo < hi {
		mid := (lo + hi) / 2
		if interp.xVals[mid] > x {
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
