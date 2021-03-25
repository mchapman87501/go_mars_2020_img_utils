package color

import (
	"image/color"
	"math"
)

type HSV struct {
	H, S, V float64 // TODO consider a fixed-point repr, e.g., int32
}

// Conform to color.Color interface:
func (c HSV) RGBA() (r, g, b, a uint32) {
	r, g, b = hsvToRGB(c.H, c.S, c.V)
	a = 0xffff
	return
}

// Conform to color.ColorModel:
func hsvModel(c color.Color) color.Color {
	if _, ok := c.(HSV); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	h, s, v := rgbToHSV(r, g, b)
	return HSV{h, s, v}
}

var HSVModel color.Model = color.ModelFunc(hsvModel)

func max(v ...float64) float64 {
	if len(v) <= 0 {
		return 0.0 // Should return an error?
	}
	result := v[0]
	for _, vCurr := range v {
		if vCurr > result {
			result = vCurr
		}
	}
	return result
}

func min(v ...float64) float64 {
	if len(v) <= 0 {
		return 0.0 // Should return an error?
	}
	result := v[0]
	for _, vCurr := range v {
		if vCurr < result {
			result = vCurr
		}
	}
	return result
}

// Convert a value in the range 0...0xffff to a value in 0.0...1.0:
func norm(v uint32) float64 {
	return float64(v) / float64(0xffff)
}

// Convert a value in 0.0 ... 1.0 to a value in 0.0 ... 0xffff
func denorm(v float64) uint32 {
	return uint32(v * 0xffff)
}

func rgbToHue(rn, gn, bn, minComp, maxComp float64) float64 {
	const hueSextant = 1.0 / 6.0
	chroma := maxComp - minComp

	// These exact comparisons should be ok given that minComp and maxComp
	// are each set to one of rn, gn or bn.
	if chroma <= 0.0 {
		return 0.0
	} else if rn == maxComp {
		result := hueSextant * (gn - bn) / chroma
		if result < 0.0 {
			result += 1
		}
		return result
	} else if gn == maxComp {
		return hueSextant * ((bn-rn)/chroma + 2.0)
	} else {
		// bn == maxComp
		return hueSextant * ((rn-gn)/chroma + 4.0)
	}
}

// Convert "normalized" RGB, with all components in 0.0 ... 1.0, to HSV.
// Returns values in the range 0.0 ... 1.0
func normRGBToHSV(rn, gn, bn float64) (float64, float64, float64) {
	// This derives from Wikipedia.
	maxComp := max(rn, gn, bn)
	minComp := min(rn, gn, bn)

	h := rgbToHue(rn, gn, bn, minComp, maxComp)
	s := 0.0
	if maxComp > 0.0 {
		s = (maxComp - minComp) / maxComp
	}
	v := maxComp
	return h, s, v
}

// Convert RGB color to HSV.
// Returns values in the range 0.0 ... 1.0
func rgbToHSV(r, g, b uint32) (float64, float64, float64) {
	return normRGBToHSV(norm(r), norm(g), norm(b))
}

// Convert HSV color to "normalized" RGB,
// with each result color component in 0.0 ... 1.0
// Round-trip RGB->HSV->RGB does not always succeed.
func hsvToNormRGB(h, s, v float64) (float64, float64, float64) {
	if s == 0 {
		return v, v, v
	}

	// This is from Wikipedia.
	chroma := s * v
	h6 := h * 6.0
	x := chroma * (1.0 - math.Abs(math.Mod(h6, 2.0)-1.0))

	var r1, g1, b1 float64
	if (0.0 <= h6) && (h6 <= 1) {
		r1 = chroma
		g1 = x
		b1 = 0
	} else if (1 < h6) && (h6 <= 2) {
		r1 = x
		g1 = chroma
		b1 = 0
	} else if (2 < h6) && (h6 <= 3) {
		r1 = 0
		g1 = chroma
		b1 = x
	} else if (3 < h6) && (h6 <= 4) {
		r1 = 0
		g1 = x
		b1 = chroma
	} else if (4 < h6) && (h6 <= 5) {
		r1 = x
		g1 = 0
		b1 = chroma
	} else {
		r1 = chroma
		g1 = 0
		b1 = x
	}

	m := v - chroma
	return r1 + m, g1 + m, b1 + m
}

// Convert HSV color to RGB.
// h, s, v are each in 0.0 ... 1.0.
// Despite being uint32, the returned values are each in 0 ... 0xffff
// I think this matches the expected behavior for Go's image/color code.
// The algorithm is from Python's colorsys.
func hsvToRGB(h, s, v float64) (uint32, uint32, uint32) {
	rn, gn, bn := hsvToNormRGB(h, s, v)
	return denorm(rn), denorm(gn), denorm(bn)
}
