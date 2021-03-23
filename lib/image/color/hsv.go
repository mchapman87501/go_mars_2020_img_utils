package hsv_color

import (
	"image/color"
	"math"
)

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

// Convert "normalized" RGB, with all components in 0.0 ... 1.0, to HSV.
// Returns values in the range 0.0 ... 1.0
func NormRGBToHSV(rn, gn, bn float64) (float64, float64, float64) {
	// This derives from colorsys in the Python standard library.
	// "Comp" means "color Component"
	maxComp := max(rn, gn, bn)
	minComp := min(rn, gn, bn)

	// Normalized value, 0...1
	v := maxComp
	if minComp == maxComp {
		return 0, 0, v
	}

	// Normalized saturation, 0...1
	s := (maxComp - minComp) / maxComp

	compRange := maxComp - minComp
	// How far from maximum value toward minimum value are r, g, b?
	rComplement := (maxComp - rn) / compRange
	gComplement := (maxComp - gn) / compRange
	bComplement := (maxComp - bn) / compRange

	// Normalized hue, 0...1
	var h float64
	if rn == maxComp {
		h = bComplement - gComplement
	} else if gn == maxComp {
		h = 2 + rComplement - bComplement
	} else {
		h = 4 + gComplement - rComplement
	}
	if h < 0.0 {
		h += 6.0
	}
	h = math.Mod((h / 6.0), 1.0)

	return h, s, v
}

// Convert RGB color to HSV.
// Returns values in the range 0.0 ... 1.0
func RGBToHSV(r, g, b uint32) (float64, float64, float64) {
	return NormRGBToHSV(norm(r), norm(g), norm(b))
}

func RGB8ToHSV(r, g, b uint8) (float64, float64, float64) {
	// This is blindly copied from image/color/color.go.
	r32 := uint32(r)
	r32 |= r32 << 8
	g32 := uint32(g)
	g32 |= g32 << 8
	b32 := uint32(b)
	b32 |= b32 << 8
	return RGBToHSV(r32, g32, b32)
}

// Convert HSV color to "normalized" RGB,
// with each result color component in 0.0 ... 1.0
// Round-trip RGB->HSV->RGB does not always succeed.
func HSVToNormRGB(h, s, v float64) (float64, float64, float64) {
	if s == 0 {
		return v, v, v
	}
	h6 := h * 6.0
	i := int(math.Floor(h6))
	f := h6 - float64(i) // nearest hue sextant
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	switch i % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	case 5:
		return v, p, q
	}
	// Should not be able to get here...
	return 0, 0, 0
}

// Convert HSV color to RGB.
// h, s, v are each in 0.0 ... 1.0.
// Despite being uint32, the returned values are each in 0 ... 0xffff
// I think this matches the expected behavior for Go's image/color code.
// The algorithm is from Python's colorsys.
func HSVToRGB(h, s, v float64) (uint32, uint32, uint32) {
	rn, gn, bn := HSVToNormRGB(h, s, v)
	return denorm(rn), denorm(gn), denorm(bn)
}

type HSV struct {
	H, S, V float64 // TODO consider a fixed-point repr, e.g., int32
}

// Conform to color.Color interface:
func (c HSV) RGBA() (r, g, b, a uint32) {
	r8, g8, b8 := HSVToRGB(c.H, c.S, c.V)
	r = uint32(r8)
	// Blindly copying from image/color/color.go - '|=' instead of '='?
	r |= r << 8
	g = uint32(g8)
	g |= g << 8
	b = uint32(b8)
	b |= b << 8
	a = math.MaxUint32
	return
}

// Conform to color.ColorModel:
func hsvModel(c color.Color) color.Color {
	if _, ok := c.(HSV); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	h, s, v := RGBToHSV(r, g, b)
	return HSV{h, s, v}
}

var HSVModel color.Model = color.ModelFunc(hsvModel)
