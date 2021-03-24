package hsv_color

import (
	"fmt"
	"testing"
)

// Are two values approximately equal?
func eq(a, b float64) bool {
	aRound := int(a*100.0 + 0.5)
	bRound := int(b*100.0 + 0.5)
	result := aRound == bRound
	if !result {
		fmt.Println(a, "!=", b, "(", aRound, "!=", bRound, ")")
	}
	return result
}

func hsvStr(h, s, v float64) string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", h, s, v)
}

func rgb2hsvTC(r, g, b uint32, hWant, sWant, vWant float64, t *testing.T) {
	hGot, sGot, vGot := RGBToHSV(r, g, b)
	if !(eq(hWant, hGot) && eq(sWant, sGot) && eq(vWant, vGot)) {
		rgbStr := fmt.Sprintf("(%d, %d, %d)", r, g, b)
		wantStr := hsvStr(hWant, sWant, vWant)
		gotStr := hsvStr(hGot, sGot, vGot)
		t.Errorf("RGBToHSV failed for %s: want %s, got %s", rgbStr, wantStr, gotStr)
	}
}

func TestRGBToHSV(t *testing.T) {
	testCases := []struct {
		r, g, b             uint32
		hWant, sWant, vWant float64
	}{
		// Test cases generated using Python's colorsys
		{0x0, 0x0, 0x0, 0.000, 0.000, 0.000},
		{0x0, 0x0, 0x3333, 0.667, 1.000, 0.200},
		{0x0, 0x0, 0x6666, 0.667, 1.000, 0.400},
		{0x0, 0x0, 0x9999, 0.667, 1.000, 0.600},
		{0x0, 0x0, 0xcccc, 0.667, 1.000, 0.800},
		{0x0, 0x0, 0xffff, 0.667, 1.000, 1.000},
		{0x0, 0x3333, 0x0, 0.333, 1.000, 0.200},
		{0x0, 0x3333, 0x3333, 0.500, 1.000, 0.200},
		{0x0, 0x3333, 0x6666, 0.583, 1.000, 0.400},
		{0x0, 0x3333, 0x9999, 0.611, 1.000, 0.600},
		{0x0, 0x3333, 0xcccc, 0.625, 1.000, 0.800},
		{0x0, 0x3333, 0xffff, 0.633, 1.000, 1.000},
		{0x0, 0x6666, 0x0, 0.333, 1.000, 0.400},
		{0x0, 0x6666, 0x3333, 0.417, 1.000, 0.400},
		{0x0, 0x6666, 0x6666, 0.500, 1.000, 0.400},
		{0x0, 0x6666, 0x9999, 0.556, 1.000, 0.600},
		{0x0, 0x6666, 0xcccc, 0.583, 1.000, 0.800},
		{0x0, 0x6666, 0xffff, 0.600, 1.000, 1.000},
		{0x0, 0x9999, 0x0, 0.333, 1.000, 0.600},
		{0x0, 0x9999, 0x3333, 0.389, 1.000, 0.600},
		{0x0, 0x9999, 0x6666, 0.444, 1.000, 0.600},
		{0x0, 0x9999, 0x9999, 0.500, 1.000, 0.600},
		{0x0, 0x9999, 0xcccc, 0.542, 1.000, 0.800},
		{0x0, 0x9999, 0xffff, 0.567, 1.000, 1.000},
		{0x0, 0xcccc, 0x0, 0.333, 1.000, 0.800},
		{0x0, 0xcccc, 0x3333, 0.375, 1.000, 0.800},
		{0x0, 0xcccc, 0x6666, 0.417, 1.000, 0.800},
		{0x0, 0xcccc, 0x9999, 0.458, 1.000, 0.800},
		{0x0, 0xcccc, 0xcccc, 0.500, 1.000, 0.800},
		{0x0, 0xcccc, 0xffff, 0.533, 1.000, 1.000},
		{0x0, 0xffff, 0x0, 0.333, 1.000, 1.000},
		{0x0, 0xffff, 0x3333, 0.367, 1.000, 1.000},
		{0x0, 0xffff, 0x6666, 0.400, 1.000, 1.000},
		{0x0, 0xffff, 0x9999, 0.433, 1.000, 1.000},
		{0x0, 0xffff, 0xcccc, 0.467, 1.000, 1.000},
		{0x0, 0xffff, 0xffff, 0.500, 1.000, 1.000},
		{0x3333, 0x0, 0x0, 0.000, 1.000, 0.200},
		{0x3333, 0x0, 0x3333, 0.833, 1.000, 0.200},
		{0x3333, 0x0, 0x6666, 0.750, 1.000, 0.400},
		{0x3333, 0x0, 0x9999, 0.722, 1.000, 0.600},
		{0x3333, 0x0, 0xcccc, 0.708, 1.000, 0.800},
		{0x3333, 0x0, 0xffff, 0.700, 1.000, 1.000},
		{0x3333, 0x3333, 0x0, 0.167, 1.000, 0.200},
		{0x3333, 0x3333, 0x3333, 0.000, 0.000, 0.200},
		{0x3333, 0x3333, 0x6666, 0.667, 0.500, 0.400},
		{0x3333, 0x3333, 0x9999, 0.667, 0.667, 0.600},
		{0x3333, 0x3333, 0xcccc, 0.667, 0.750, 0.800},
		{0x3333, 0x3333, 0xffff, 0.667, 0.800, 1.000},
		{0x3333, 0x6666, 0x0, 0.250, 1.000, 0.400},
		{0x3333, 0x6666, 0x3333, 0.333, 0.500, 0.400},
		{0x3333, 0x6666, 0x6666, 0.500, 0.500, 0.400},
		{0x3333, 0x6666, 0x9999, 0.583, 0.667, 0.600},
		{0x3333, 0x6666, 0xcccc, 0.611, 0.750, 0.800},
		{0x3333, 0x6666, 0xffff, 0.625, 0.800, 1.000},
		{0x3333, 0x9999, 0x0, 0.278, 1.000, 0.600},
		{0x3333, 0x9999, 0x3333, 0.333, 0.667, 0.600},
		{0x3333, 0x9999, 0x6666, 0.417, 0.667, 0.600},
		{0x3333, 0x9999, 0x9999, 0.500, 0.667, 0.600},
		{0x3333, 0x9999, 0xcccc, 0.556, 0.750, 0.800},
		{0x3333, 0x9999, 0xffff, 0.583, 0.800, 1.000},
		{0x3333, 0xcccc, 0x0, 0.292, 1.000, 0.800},
		{0x3333, 0xcccc, 0x3333, 0.333, 0.750, 0.800},
		{0x3333, 0xcccc, 0x6666, 0.389, 0.750, 0.800},
		{0x3333, 0xcccc, 0x9999, 0.444, 0.750, 0.800},
		{0x3333, 0xcccc, 0xcccc, 0.500, 0.750, 0.800},
		{0x3333, 0xcccc, 0xffff, 0.542, 0.800, 1.000},
		{0x3333, 0xffff, 0x0, 0.300, 1.000, 1.000},
		{0x3333, 0xffff, 0x3333, 0.333, 0.800, 1.000},
		{0x3333, 0xffff, 0x6666, 0.375, 0.800, 1.000},
		{0x3333, 0xffff, 0x9999, 0.417, 0.800, 1.000},
		{0x3333, 0xffff, 0xcccc, 0.458, 0.800, 1.000},
		{0x3333, 0xffff, 0xffff, 0.500, 0.800, 1.000},
		{0x6666, 0x0, 0x0, 0.000, 1.000, 0.400},
		{0x6666, 0x0, 0x3333, 0.917, 1.000, 0.400},
		{0x6666, 0x0, 0x6666, 0.833, 1.000, 0.400},
		{0x6666, 0x0, 0x9999, 0.778, 1.000, 0.600},
		{0x6666, 0x0, 0xcccc, 0.750, 1.000, 0.800},
		{0x6666, 0x0, 0xffff, 0.733, 1.000, 1.000},
		{0x6666, 0x3333, 0x0, 0.083, 1.000, 0.400},
		{0x6666, 0x3333, 0x3333, 0.000, 0.500, 0.400},
		{0x6666, 0x3333, 0x6666, 0.833, 0.500, 0.400},
		{0x6666, 0x3333, 0x9999, 0.750, 0.667, 0.600},
		{0x6666, 0x3333, 0xcccc, 0.722, 0.750, 0.800},
		{0x6666, 0x3333, 0xffff, 0.708, 0.800, 1.000},
		{0x6666, 0x6666, 0x0, 0.167, 1.000, 0.400},
		{0x6666, 0x6666, 0x3333, 0.167, 0.500, 0.400},
		{0x6666, 0x6666, 0x6666, 0.000, 0.000, 0.400},
		{0x6666, 0x6666, 0x9999, 0.667, 0.333, 0.600},
		{0x6666, 0x6666, 0xcccc, 0.667, 0.500, 0.800},
		{0x6666, 0x6666, 0xffff, 0.667, 0.600, 1.000},
		{0x6666, 0x9999, 0x0, 0.222, 1.000, 0.600},
		{0x6666, 0x9999, 0x3333, 0.250, 0.667, 0.600},
		{0x6666, 0x9999, 0x6666, 0.333, 0.333, 0.600},
		{0x6666, 0x9999, 0x9999, 0.500, 0.333, 0.600},
		{0x6666, 0x9999, 0xcccc, 0.583, 0.500, 0.800},
		{0x6666, 0x9999, 0xffff, 0.611, 0.600, 1.000},
		{0x6666, 0xcccc, 0x0, 0.250, 1.000, 0.800},
		{0x6666, 0xcccc, 0x3333, 0.278, 0.750, 0.800},
		{0x6666, 0xcccc, 0x6666, 0.333, 0.500, 0.800},
		{0x6666, 0xcccc, 0x9999, 0.417, 0.500, 0.800},
		{0x6666, 0xcccc, 0xcccc, 0.500, 0.500, 0.800},
		{0x6666, 0xcccc, 0xffff, 0.556, 0.600, 1.000},
		{0x6666, 0xffff, 0x0, 0.267, 1.000, 1.000},
		{0x6666, 0xffff, 0x3333, 0.292, 0.800, 1.000},
		{0x6666, 0xffff, 0x6666, 0.333, 0.600, 1.000},
		{0x6666, 0xffff, 0x9999, 0.389, 0.600, 1.000},
		{0x6666, 0xffff, 0xcccc, 0.444, 0.600, 1.000},
		{0x6666, 0xffff, 0xffff, 0.500, 0.600, 1.000},
		{0x9999, 0x0, 0x0, 0.000, 1.000, 0.600},
		{0x9999, 0x0, 0x3333, 0.944, 1.000, 0.600},
		{0x9999, 0x0, 0x6666, 0.889, 1.000, 0.600},
		{0x9999, 0x0, 0x9999, 0.833, 1.000, 0.600},
		{0x9999, 0x0, 0xcccc, 0.792, 1.000, 0.800},
		{0x9999, 0x0, 0xffff, 0.767, 1.000, 1.000},
		{0x9999, 0x3333, 0x0, 0.056, 1.000, 0.600},
		{0x9999, 0x3333, 0x3333, 0.000, 0.667, 0.600},
		{0x9999, 0x3333, 0x6666, 0.917, 0.667, 0.600},
		{0x9999, 0x3333, 0x9999, 0.833, 0.667, 0.600},
		{0x9999, 0x3333, 0xcccc, 0.778, 0.750, 0.800},
		{0x9999, 0x3333, 0xffff, 0.750, 0.800, 1.000},
		{0x9999, 0x6666, 0x0, 0.111, 1.000, 0.600},
		{0x9999, 0x6666, 0x3333, 0.083, 0.667, 0.600},
		{0x9999, 0x6666, 0x6666, 0.000, 0.333, 0.600},
		{0x9999, 0x6666, 0x9999, 0.833, 0.333, 0.600},
		{0x9999, 0x6666, 0xcccc, 0.750, 0.500, 0.800},
		{0x9999, 0x6666, 0xffff, 0.722, 0.600, 1.000},
		{0x9999, 0x9999, 0x0, 0.167, 1.000, 0.600},
		{0x9999, 0x9999, 0x3333, 0.167, 0.667, 0.600},
		{0x9999, 0x9999, 0x6666, 0.167, 0.333, 0.600},
		{0x9999, 0x9999, 0x9999, 0.000, 0.000, 0.600},
		{0x9999, 0x9999, 0xcccc, 0.667, 0.250, 0.800},
		{0x9999, 0x9999, 0xffff, 0.667, 0.400, 1.000},
		{0x9999, 0xcccc, 0x0, 0.208, 1.000, 0.800},
		{0x9999, 0xcccc, 0x3333, 0.222, 0.750, 0.800},
		{0x9999, 0xcccc, 0x6666, 0.250, 0.500, 0.800},
		{0x9999, 0xcccc, 0x9999, 0.333, 0.250, 0.800},
		{0x9999, 0xcccc, 0xcccc, 0.500, 0.250, 0.800},
		{0x9999, 0xcccc, 0xffff, 0.583, 0.400, 1.000},
		{0x9999, 0xffff, 0x0, 0.233, 1.000, 1.000},
		{0x9999, 0xffff, 0x3333, 0.250, 0.800, 1.000},
		{0x9999, 0xffff, 0x6666, 0.278, 0.600, 1.000},
		{0x9999, 0xffff, 0x9999, 0.333, 0.400, 1.000},
		{0x9999, 0xffff, 0xcccc, 0.417, 0.400, 1.000},
		{0x9999, 0xffff, 0xffff, 0.500, 0.400, 1.000},
		{0xcccc, 0x0, 0x0, 0.000, 1.000, 0.800},
		{0xcccc, 0x0, 0x3333, 0.958, 1.000, 0.800},
		{0xcccc, 0x0, 0x6666, 0.917, 1.000, 0.800},
		{0xcccc, 0x0, 0x9999, 0.875, 1.000, 0.800},
		{0xcccc, 0x0, 0xcccc, 0.833, 1.000, 0.800},
		{0xcccc, 0x0, 0xffff, 0.800, 1.000, 1.000},
		{0xcccc, 0x3333, 0x0, 0.042, 1.000, 0.800},
		{0xcccc, 0x3333, 0x3333, 0.000, 0.750, 0.800},
		{0xcccc, 0x3333, 0x6666, 0.944, 0.750, 0.800},
		{0xcccc, 0x3333, 0x9999, 0.889, 0.750, 0.800},
		{0xcccc, 0x3333, 0xcccc, 0.833, 0.750, 0.800},
		{0xcccc, 0x3333, 0xffff, 0.792, 0.800, 1.000},
		{0xcccc, 0x6666, 0x0, 0.083, 1.000, 0.800},
		{0xcccc, 0x6666, 0x3333, 0.056, 0.750, 0.800},
		{0xcccc, 0x6666, 0x6666, 0.000, 0.500, 0.800},
		{0xcccc, 0x6666, 0x9999, 0.917, 0.500, 0.800},
		{0xcccc, 0x6666, 0xcccc, 0.833, 0.500, 0.800},
		{0xcccc, 0x6666, 0xffff, 0.778, 0.600, 1.000},
		{0xcccc, 0x9999, 0x0, 0.125, 1.000, 0.800},
		{0xcccc, 0x9999, 0x3333, 0.111, 0.750, 0.800},
		{0xcccc, 0x9999, 0x6666, 0.083, 0.500, 0.800},
		{0xcccc, 0x9999, 0x9999, 0.000, 0.250, 0.800},
		{0xcccc, 0x9999, 0xcccc, 0.833, 0.250, 0.800},
		{0xcccc, 0x9999, 0xffff, 0.750, 0.400, 1.000},
		{0xcccc, 0xcccc, 0x0, 0.167, 1.000, 0.800},
		{0xcccc, 0xcccc, 0x3333, 0.167, 0.750, 0.800},
		{0xcccc, 0xcccc, 0x6666, 0.167, 0.500, 0.800},
		{0xcccc, 0xcccc, 0x9999, 0.167, 0.250, 0.800},
		{0xcccc, 0xcccc, 0xcccc, 0.000, 0.000, 0.800},
		{0xcccc, 0xcccc, 0xffff, 0.667, 0.200, 1.000},
		{0xcccc, 0xffff, 0x0, 0.200, 1.000, 1.000},
		{0xcccc, 0xffff, 0x3333, 0.208, 0.800, 1.000},
		{0xcccc, 0xffff, 0x6666, 0.222, 0.600, 1.000},
		{0xcccc, 0xffff, 0x9999, 0.250, 0.400, 1.000},
		{0xcccc, 0xffff, 0xcccc, 0.333, 0.200, 1.000},
		{0xcccc, 0xffff, 0xffff, 0.500, 0.200, 1.000},
		{0xffff, 0x0, 0x0, 0.000, 1.000, 1.000},
		{0xffff, 0x0, 0x3333, 0.967, 1.000, 1.000},
		{0xffff, 0x0, 0x6666, 0.933, 1.000, 1.000},
		{0xffff, 0x0, 0x9999, 0.900, 1.000, 1.000},
		{0xffff, 0x0, 0xcccc, 0.867, 1.000, 1.000},
		{0xffff, 0x0, 0xffff, 0.833, 1.000, 1.000},
		{0xffff, 0x3333, 0x0, 0.033, 1.000, 1.000},
		{0xffff, 0x3333, 0x3333, 0.000, 0.800, 1.000},
		{0xffff, 0x3333, 0x6666, 0.958, 0.800, 1.000},
		{0xffff, 0x3333, 0x9999, 0.917, 0.800, 1.000},
		{0xffff, 0x3333, 0xcccc, 0.875, 0.800, 1.000},
		{0xffff, 0x3333, 0xffff, 0.833, 0.800, 1.000},
		{0xffff, 0x6666, 0x0, 0.067, 1.000, 1.000},
		{0xffff, 0x6666, 0x3333, 0.042, 0.800, 1.000},
		{0xffff, 0x6666, 0x6666, 0.000, 0.600, 1.000},
		{0xffff, 0x6666, 0x9999, 0.944, 0.600, 1.000},
		{0xffff, 0x6666, 0xcccc, 0.889, 0.600, 1.000},
		{0xffff, 0x6666, 0xffff, 0.833, 0.600, 1.000},
		{0xffff, 0x9999, 0x0, 0.100, 1.000, 1.000},
		{0xffff, 0x9999, 0x3333, 0.083, 0.800, 1.000},
		{0xffff, 0x9999, 0x6666, 0.056, 0.600, 1.000},
		{0xffff, 0x9999, 0x9999, 0.000, 0.400, 1.000},
		{0xffff, 0x9999, 0xcccc, 0.917, 0.400, 1.000},
		{0xffff, 0x9999, 0xffff, 0.833, 0.400, 1.000},
		{0xffff, 0xcccc, 0x0, 0.133, 1.000, 1.000},
		{0xffff, 0xcccc, 0x3333, 0.125, 0.800, 1.000},
		{0xffff, 0xcccc, 0x6666, 0.111, 0.600, 1.000},
		{0xffff, 0xcccc, 0x9999, 0.083, 0.400, 1.000},
		{0xffff, 0xcccc, 0xcccc, 0.000, 0.200, 1.000},
		{0xffff, 0xcccc, 0xffff, 0.833, 0.200, 1.000},
		{0xffff, 0xffff, 0x0, 0.167, 1.000, 1.000},
		{0xffff, 0xffff, 0x3333, 0.167, 0.800, 1.000},
		{0xffff, 0xffff, 0x6666, 0.167, 0.600, 1.000},
		{0xffff, 0xffff, 0x9999, 0.167, 0.400, 1.000},
		{0xffff, 0xffff, 0xcccc, 0.167, 0.200, 1.000},
		{0xffff, 0xffff, 0xffff, 0.000, 0.000, 1.000},
	}
	for _, tc := range testCases {
		rgb2hsvTC(tc.r, tc.g, tc.b, tc.hWant, tc.sWant, tc.vWant, t)
	}
}

func TestHSVToRGB(t *testing.T) {
	testCases := []struct {
		h, s, v             float64
		rWant, gWant, bWant uint32
	}{
		{0.000, 0.000, 0.000, 0x0, 0x0, 0x0},
		{0.000, 0.000, 0.250, 0x3fffffff, 0x3fffffff, 0x3fffffff},
		{0.000, 0.000, 0.500, 0x7fffffff, 0x7fffffff, 0x7fffffff},
		{0.000, 0.000, 0.750, 0xbfffffff, 0xbfffffff, 0xbfffffff},
		{0.000, 0.000, 1.000, 0xffff, 0xffff, 0xffff},
		{0.000, 0.250, 0.000, 0x0, 0x0, 0x0},
		{0.000, 0.250, 0.250, 0x3fffffff, 0x2fffffff, 0x2fffffff},
		{0.000, 0.250, 0.500, 0x7fffffff, 0x5fffffff, 0x5fffffff},
		{0.000, 0.250, 0.750, 0xbfffffff, 0x8fffffff, 0x8fffffff},
		{0.000, 0.250, 1.000, 0xffff, 0xbfffffff, 0xbfffffff},
		{0.000, 0.500, 0.000, 0x0, 0x0, 0x0},
		{0.000, 0.500, 0.250, 0x3fffffff, 0x1fffffff, 0x1fffffff},
		{0.000, 0.500, 0.500, 0x7fffffff, 0x3fffffff, 0x3fffffff},
		{0.000, 0.500, 0.750, 0xbfffffff, 0x5fffffff, 0x5fffffff},
		{0.000, 0.500, 1.000, 0xffff, 0x7fffffff, 0x7fffffff},
		{0.000, 0.750, 0.000, 0x0, 0x0, 0x0},
		{0.000, 0.750, 0.250, 0x3fffffff, 0xfffffff, 0xfffffff},
		{0.000, 0.750, 0.500, 0x7fffffff, 0x1fffffff, 0x1fffffff},
		{0.000, 0.750, 0.750, 0xbfffffff, 0x2fffffff, 0x2fffffff},
		{0.000, 0.750, 1.000, 0xffff, 0x3fffffff, 0x3fffffff},
		{0.000, 1.000, 0.000, 0x0, 0x0, 0x0},
		{0.000, 1.000, 0.250, 0x3fffffff, 0x0, 0x0},
		{0.000, 1.000, 0.500, 0x7fffffff, 0x0, 0x0},
		{0.000, 1.000, 0.750, 0xbfffffff, 0x0, 0x0},
		{0.000, 1.000, 1.000, 0xffff, 0x0, 0x0},
		{0.250, 0.000, 0.000, 0x0, 0x0, 0x0},
		{0.250, 0.000, 0.250, 0x3fffffff, 0x3fffffff, 0x3fffffff},
		{0.250, 0.000, 0.500, 0x7fffffff, 0x7fffffff, 0x7fffffff},
		{0.250, 0.000, 0.750, 0xbfffffff, 0xbfffffff, 0xbfffffff},
		{0.250, 0.000, 1.000, 0xffff, 0xffff, 0xffff},
		{0.250, 0.250, 0.000, 0x0, 0x0, 0x0},
		{0.250, 0.250, 0.250, 0x37ffffff, 0x3fffffff, 0x2fffffff},
		{0.250, 0.250, 0.500, 0x6fffffff, 0x7fffffff, 0x5fffffff},
		{0.250, 0.250, 0.750, 0xa7ffffff, 0xbfffffff, 0x8fffffff},
		{0.250, 0.250, 1.000, 0xdfffffff, 0xffff, 0xbfffffff},
		{0.250, 0.500, 0.000, 0x0, 0x0, 0x0},
		{0.250, 0.500, 0.250, 0x2fffffff, 0x3fffffff, 0x1fffffff},
		{0.250, 0.500, 0.500, 0x5fffffff, 0x7fffffff, 0x3fffffff},
		{0.250, 0.500, 0.750, 0x8fffffff, 0xbfffffff, 0x5fffffff},
		{0.250, 0.500, 1.000, 0xbfffffff, 0xffff, 0x7fffffff},
		{0.250, 0.750, 0.000, 0x0, 0x0, 0x0},
		{0.250, 0.750, 0.250, 0x27ffffff, 0x3fffffff, 0xfffffff},
		{0.250, 0.750, 0.500, 0x4fffffff, 0x7fffffff, 0x1fffffff},
		{0.250, 0.750, 0.750, 0x77ffffff, 0xbfffffff, 0x2fffffff},
		{0.250, 0.750, 1.000, 0x9fffffff, 0xffff, 0x3fffffff},
		{0.250, 1.000, 0.000, 0x0, 0x0, 0x0},
		{0.250, 1.000, 0.250, 0x1fffffff, 0x3fffffff, 0x0},
		{0.250, 1.000, 0.500, 0x3fffffff, 0x7fffffff, 0x0},
		{0.250, 1.000, 0.750, 0x5fffffff, 0xbfffffff, 0x0},
		{0.250, 1.000, 1.000, 0x7fffffff, 0xffff, 0x0},
		{0.500, 0.000, 0.000, 0x0, 0x0, 0x0},
		{0.500, 0.000, 0.250, 0x3fffffff, 0x3fffffff, 0x3fffffff},
		{0.500, 0.000, 0.500, 0x7fffffff, 0x7fffffff, 0x7fffffff},
		{0.500, 0.000, 0.750, 0xbfffffff, 0xbfffffff, 0xbfffffff},
		{0.500, 0.000, 1.000, 0xffff, 0xffff, 0xffff},
		{0.500, 0.250, 0.000, 0x0, 0x0, 0x0},
		{0.500, 0.250, 0.250, 0x2fffffff, 0x3fffffff, 0x3fffffff},
		{0.500, 0.250, 0.500, 0x5fffffff, 0x7fffffff, 0x7fffffff},
		{0.500, 0.250, 0.750, 0x8fffffff, 0xbfffffff, 0xbfffffff},
		{0.500, 0.250, 1.000, 0xbfffffff, 0xffff, 0xffff},
		{0.500, 0.500, 0.000, 0x0, 0x0, 0x0},
		{0.500, 0.500, 0.250, 0x1fffffff, 0x3fffffff, 0x3fffffff},
		{0.500, 0.500, 0.500, 0x3fffffff, 0x7fffffff, 0x7fffffff},
		{0.500, 0.500, 0.750, 0x5fffffff, 0xbfffffff, 0xbfffffff},
		{0.500, 0.500, 1.000, 0x7fffffff, 0xffff, 0xffff},
		{0.500, 0.750, 0.000, 0x0, 0x0, 0x0},
		{0.500, 0.750, 0.250, 0xfffffff, 0x3fffffff, 0x3fffffff},
		{0.500, 0.750, 0.500, 0x1fffffff, 0x7fffffff, 0x7fffffff},
		{0.500, 0.750, 0.750, 0x2fffffff, 0xbfffffff, 0xbfffffff},
		{0.500, 0.750, 1.000, 0x3fffffff, 0xffff, 0xffff},
		{0.500, 1.000, 0.000, 0x0, 0x0, 0x0},
		{0.500, 1.000, 0.250, 0x0, 0x3fffffff, 0x3fffffff},
		{0.500, 1.000, 0.500, 0x0, 0x7fffffff, 0x7fffffff},
		{0.500, 1.000, 0.750, 0x0, 0xbfffffff, 0xbfffffff},
		{0.500, 1.000, 1.000, 0x0, 0xffff, 0xffff},
		{0.750, 0.000, 0.000, 0x0, 0x0, 0x0},
		{0.750, 0.000, 0.250, 0x3fffffff, 0x3fffffff, 0x3fffffff},
		{0.750, 0.000, 0.500, 0x7fffffff, 0x7fffffff, 0x7fffffff},
		{0.750, 0.000, 0.750, 0xbfffffff, 0xbfffffff, 0xbfffffff},
		{0.750, 0.000, 1.000, 0xffff, 0xffff, 0xffff},
		{0.750, 0.250, 0.000, 0x0, 0x0, 0x0},
		{0.750, 0.250, 0.250, 0x37ffffff, 0x2fffffff, 0x3fffffff},
		{0.750, 0.250, 0.500, 0x6fffffff, 0x5fffffff, 0x7fffffff},
		{0.750, 0.250, 0.750, 0xa7ffffff, 0x8fffffff, 0xbfffffff},
		{0.750, 0.250, 1.000, 0xdfffffff, 0xbfffffff, 0xffff},
		{0.750, 0.500, 0.000, 0x0, 0x0, 0x0},
		{0.750, 0.500, 0.250, 0x2fffffff, 0x1fffffff, 0x3fffffff},
		{0.750, 0.500, 0.500, 0x5fffffff, 0x3fffffff, 0x7fffffff},
		{0.750, 0.500, 0.750, 0x8fffffff, 0x5fffffff, 0xbfffffff},
		{0.750, 0.500, 1.000, 0xbfffffff, 0x7fffffff, 0xffff},
		{0.750, 0.750, 0.000, 0x0, 0x0, 0x0},
		{0.750, 0.750, 0.250, 0x27ffffff, 0xfffffff, 0x3fffffff},
		{0.750, 0.750, 0.500, 0x4fffffff, 0x1fffffff, 0x7fffffff},
		{0.750, 0.750, 0.750, 0x77ffffff, 0x2fffffff, 0xbfffffff},
		{0.750, 0.750, 1.000, 0x9fffffff, 0x3fffffff, 0xffff},
		{0.750, 1.000, 0.000, 0x0, 0x0, 0x0},
		{0.750, 1.000, 0.250, 0x1fffffff, 0x0, 0x3fffffff},
		{0.750, 1.000, 0.500, 0x3fffffff, 0x0, 0x7fffffff},
		{0.750, 1.000, 0.750, 0x5fffffff, 0x0, 0xbfffffff},
		{0.750, 1.000, 1.000, 0x7fffffff, 0x0, 0xffff},
		{1.000, 0.000, 0.000, 0x0, 0x0, 0x0},
		{1.000, 0.000, 0.250, 0x3fffffff, 0x3fffffff, 0x3fffffff},
		{1.000, 0.000, 0.500, 0x7fffffff, 0x7fffffff, 0x7fffffff},
		{1.000, 0.000, 0.750, 0xbfffffff, 0xbfffffff, 0xbfffffff},
		{1.000, 0.000, 1.000, 0xffff, 0xffff, 0xffff},
		{1.000, 0.250, 0.000, 0x0, 0x0, 0x0},
		{1.000, 0.250, 0.250, 0x3fffffff, 0x2fffffff, 0x2fffffff},
		{1.000, 0.250, 0.500, 0x7fffffff, 0x5fffffff, 0x5fffffff},
		{1.000, 0.250, 0.750, 0xbfffffff, 0x8fffffff, 0x8fffffff},
		{1.000, 0.250, 1.000, 0xffff, 0xbfffffff, 0xbfffffff},
		{1.000, 0.500, 0.000, 0x0, 0x0, 0x0},
		{1.000, 0.500, 0.250, 0x3fffffff, 0x1fffffff, 0x1fffffff},
		{1.000, 0.500, 0.500, 0x7fffffff, 0x3fffffff, 0x3fffffff},
		{1.000, 0.500, 0.750, 0xbfffffff, 0x5fffffff, 0x5fffffff},
		{1.000, 0.500, 1.000, 0xffff, 0x7fffffff, 0x7fffffff},
		{1.000, 0.750, 0.000, 0x0, 0x0, 0x0},
		{1.000, 0.750, 0.250, 0x3fffffff, 0xfffffff, 0xfffffff},
		{1.000, 0.750, 0.500, 0x7fffffff, 0x1fffffff, 0x1fffffff},
		{1.000, 0.750, 0.750, 0xbfffffff, 0x2fffffff, 0x2fffffff},
		{1.000, 0.750, 1.000, 0xffff, 0x3fffffff, 0x3fffffff},
		{1.000, 1.000, 0.000, 0x0, 0x0, 0x0},
		{1.000, 1.000, 0.250, 0x3fffffff, 0x0, 0x0},
		{1.000, 1.000, 0.500, 0x7fffffff, 0x0, 0x0},
		{1.000, 1.000, 0.750, 0xbfffffff, 0x0, 0x0},
		{1.000, 1.000, 1.000, 0xffff, 0x0, 0x0},
	}

	for _, tc := range testCases {
		rGot, gGot, bGot := HSVToRGB(tc.h, tc.s, tc.v)

		if !((rGot == tc.rWant) && (gGot == tc.gWant) && (bGot == tc.bWant)) {
			inStr := hsvStr(tc.h, tc.s, tc.v)
			t.Errorf("HSVToRGB failed for %v: want (%x, %x, %x), got (%x, %x, %x)", inStr, tc.rWant, tc.gWant, tc.bWant, rGot, gGot, bGot)
		}
	}
}

func TestRGB8ToHSV(t *testing.T) {
	testCases := []struct {
		r, g, b             uint8
		hWant, sWant, vWant float64
	}{
		{50, 0, 0, 0.0, 1.0, 50.0 / 255.0},
		{255, 255, 255, 0.0, 0.0, 1.0},
	}
	for _, tc := range testCases {
		hGot, sGot, vGot := RGB8ToHSV(tc.r, tc.g, tc.b)
		if !(eq(tc.hWant, hGot) && eq(tc.sWant, sGot) && eq(tc.vWant, vGot)) {
			t.Errorf("RGB8ToHSV failed for (%v, %v, %v): want (%v, %v, %v), got (%v, %v, %v)", tc.r, tc.g, tc.b, tc.hWant, tc.sWant, tc.vWant, hGot, sGot, vGot)
		}
	}
}

func TestRGBToHSVRoundTrip(t *testing.T) {
	var intensity uint32
	for intensity = 0; intensity <= 0xff; intensity++ {
		red := intensity + 50
		if red > 0xff {
			red = 0xff
		}
		green := intensity / 4
		blue := intensity

		h, s, v := RGBToHSV(red, green, blue)
		r, g, b := HSVToRGB(h, s, v)
		if !((r == red) && (g == green) && (b == blue)) {
			t.Errorf(
				"Round-trip failed: (%v, %v, %v) -> (%v, %v, %v)",
				red, green, blue, r, g, b)
		}
	}
}

func TestNormRGBToHSVRoundTrip(t *testing.T) {
	for intensity := 0.0; intensity <= 1.0; intensity += 0.10 {
		red := intensity + (50.0 / 255.0)
		if red > 1.0 {
			red = 1.0
		}
		green := intensity / 4.0
		blue := intensity

		h, s, v := NormRGBToHSV(red, green, blue)
		r, g, b := HSVToNormRGB(h, s, v)
		if !(eq(r, red) && eq(g, green) && eq(b, blue)) {
			t.Errorf(
				"Norm round-trip failed: (%v, %v, %v) -> (%v, %v, %v) -> (%v, %v, %v)",
				red, green, blue, h, s, v, r, g, b)
		}
	}
}

func TestNormRGBToHSVGrayRT(t *testing.T) {
	for intensity := 0.0; intensity <= 1.0; intensity += 0.01 {
		h, s, v := NormRGBToHSV(intensity, intensity, intensity)
		r, g, b := HSVToNormRGB(h, s, v)
		if !(eq(r, intensity) && eq(g, intensity) && eq(b, intensity)) {
			t.Errorf(
				"Norm round-trip failed: (%v, %v, %v) -> (%v, %v, %v) -> (%v, %v, %v)",
				intensity, intensity, intensity, h, s, v, r, g, b)
		}
	}
}

func TestColorHSVRGBA(t *testing.T) {
	// Test the HSV.RGBA function
	for intensity := 0.0; intensity <= 1.0; intensity += 0.02 {
		expected := uint32(intensity * 0xffff)
		hsv := HSV{H: 0.0, S: 0.0, V: intensity}
		r, g, b, a := hsv.RGBA()
		if !((r == expected) && (g == expected) && (b == expected)) {
			t.Errorf("Round trip failed for %v; expected rgb=%v, got (%v, %v, %v)", intensity, expected, r, g, b)
		}
		if a != 0xffff {
			t.Errorf("Expected full opaque alpha 0xffff, got %v", a)
		}
	}
}
