package lib

import (
	"fmt"
	"math"

	lib_color "com.dmoonc/mchapman87501/mars_2020_img_utils/lib/image/color"
)

// Use fixed-point values as map keys.
type chanAdjustmentMap map[float64]float64
type cntMap map[float64]float64

type AdjustmentMap struct {
	L, A, B          chanAdjustmentMap
	lCnt, aCnt, bCnt cntMap
}

func NewAdjustmentMap() *AdjustmentMap {
	return &AdjustmentMap{
		L:    make(chanAdjustmentMap),
		A:    make(chanAdjustmentMap),
		B:    make(chanAdjustmentMap),
		lCnt: make(cntMap),
		aCnt: make(cntMap),
		bCnt: make(cntMap),
	}
}

func (am *AdjustmentMap) AddSample(srcPix, targetPix lib_color.CIELab) {
	srcL := srcPix.L
	am.L[srcL] += targetPix.L
	am.lCnt[srcL] += 1

	srcA := srcPix.A
	am.A[srcA] += targetPix.A
	am.aCnt[srcA] += 1

	srcB := srcPix.B
	am.B[srcB] += targetPix.B
	am.bCnt[srcB] += 1
}

func (am *AdjustmentMap) Complete() {
	completeChan := func(cam chanAdjustmentMap, counts cntMap, chanName string) {
		for k := range cam {
			cam[k] /= counts[k]
		}
	}
	completeChan(am.L, am.lCnt, "L")
	completeChan(am.A, am.aCnt, "a")
	completeChan(am.B, am.bCnt, "b")

	am.addExtrema()
}

func (am *AdjustmentMap) addExtrema() {
	getChanExtrema := func(cam chanAdjustmentMap) (minIn, minOut, maxIn, maxOut float64) {
		first := true
		for k, v := range cam {
			if first {
				minIn = k
				maxIn = k
				minOut = v
				maxOut = v
				first = false
			} else {
				if minIn > k {
					minIn = k
					minOut = v
				}
				if maxIn < k {
					maxIn = k
					maxOut = v
				}
			}
		}
		return
	}

	addChanExtrema := func(cam chanAdjustmentMap, minVal, maxVal float64) {
		extrema := []float64{minVal, maxVal}

		// Extrapolate the extreme values using the minimum and maximum adjustment values?
		minIn, minOut, maxIn, maxOut := getChanExtrema(cam)
		dOut := maxOut - minOut
		dIn := maxIn - minIn

		if dIn > 1.0e-3 {
			for _, inVal := range extrema {
				_, ok := cam[inVal]
				if !ok {
					// Value is not yet mapped.
					outVal := (inVal-minIn)*dOut/dIn + minOut
					if math.IsNaN(outVal) {
						fmt.Println("Ooops.", minIn, minOut, maxIn, maxOut)
					}
					cam[inVal] = outVal
				}
			}
		}
	}

	addChanExtrema(am.L, 0.0, 100.0)
	addChanExtrema(am.A, -128.0, 127.0)
	addChanExtrema(am.B, -128.0, 127.0)
}
