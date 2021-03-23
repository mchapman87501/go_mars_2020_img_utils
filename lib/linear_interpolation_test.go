package lib

import (
	"testing"
)

func BisectLeftTC(interp *Float64Interpolator, x float64, wantXPrev, wantXNext float64, t *testing.T) {
	gotXPrev, gotXNext := interp.bisectLeft(x)
	if (gotXPrev != wantXPrev) || (gotXNext != wantXNext) {
		t.Errorf("Failed BisectLeft.  Given %v, wanted (%v, %v), got (%v, %v)", x, wantXPrev, wantXNext, gotXPrev, gotXNext)
	}

}
func TestBisectLeftOddCount(t *testing.T) {
	yVals := map[float64]float64{
		0.0: 10.0,
		0.5: 5.0,
		1.0: 0.0,
	}

	interp := NewFloat64Interpolator(yVals)

	testCases := []struct {
		x, wantXPrev, wantXNext float64
	}{
		{0.0, 0.0, 0.5},
		{0.25, 0.0, 0.5},
		{0.5, 0.5, 1.0},
		{0.75, 0.5, 1.0},
		{1.0, 0.5, 1.0},
		{42.0, 0.5, 1.0},
	}
	for _, tc := range testCases {
		BisectLeftTC(interp, tc.x, tc.wantXPrev, tc.wantXNext, t)
	}
}

func TestBisectLeftEvenCount(t *testing.T) {
	yVals := map[float64]float64{
		0.0: 10.0,
		0.3: 7.0,
		0.5: 5.0,
		1.0: 0.0,
	}

	interp := NewFloat64Interpolator(yVals)

	testCases := []struct {
		x, wantXPrev, wantXNext float64
	}{
		{0.0, 0.0, 0.3},
		{0.25, 0.0, 0.3},
		{0.3, 0.3, 0.5},
		{0.31, 0.3, 0.5},
		{0.5, 0.5, 1.0},
		{0.75, 0.5, 1.0},
		{1.0, 0.5, 1.0},
		{42.0, 0.5, 1.0},
	}
	for _, tc := range testCases {
		BisectLeftTC(interp, tc.x, tc.wantXPrev, tc.wantXNext, t)
	}
}
