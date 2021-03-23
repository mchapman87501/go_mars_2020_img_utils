package lib

import (
	"math"
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

func InterpTC(interp *Float64Interpolator, x float64, want float64, t *testing.T) {
	got := interp.Interp(x)
	// Test for approximate equality.
	const epsilon = 1.0e-5

	diff := math.Abs(got - want)
	if diff > epsilon {
		t.Errorf("interp(%v), wanted %v, got %v", x, want, got)
	}
}

func TestInterpOddCount(t *testing.T) {
	yVals := map[float64]float64{
		0.0: 10.0,
		0.5: 5.0,
		1.0: 0.0,
	}

	interp := NewFloat64Interpolator(yVals)

	testCases := []struct {
		x, want float64
	}{
		{0.0, 10.0},
		{0.25, 7.5},
		{0.5, 5.0},
		{0.75, 2.5},
		{1.0, 0.0},
		// Linear extrapolation:
		{42.0, -410.0},
	}
	for _, tc := range testCases {
		InterpTC(interp, tc.x, tc.want, t)
	}
}

func TestInterpEvenCount(t *testing.T) {
	yVals := map[float64]float64{
		0.0: 0.0,
		0.3: 3.0,
		0.5: 5.0,
		1.0: 10.0,
	}

	interp := NewFloat64Interpolator(yVals)

	testCases := []struct {
		x, want float64
	}{
		{-0.5, -5.0},
		{0.0, 0.0},
		{0.25, 2.5},
		{0.5, 5.0},
		{0.75, 7.5},
		{1.0, 10.0},
		{42.0, 420.0},
	}
	for _, tc := range testCases {
		InterpTC(interp, tc.x, tc.want, t)
	}
}
