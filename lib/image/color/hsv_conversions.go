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

// Convert a value in the range 0...0xffffffff to a value in 0.0...1.0:
func norm(v uint32) float64 {
	return float64(v) / float64(0xffffffff)
}

// Convert a value in 0.0 ... 1.0 to a value in 0.0 ... 0xffffffff
func denorm(v float64) uint32 {
	return uint32(v * 0xffffffff)
}

// Convert RGB color to HSV.
// Returns values in the range 0.0 ... 1.0
func RGBToHSV(r, g, b uint32) (float64, float64, float64) {
	rn := norm(r)
	gn := norm(g)
	bn := norm(b)

	// This derives from colorsys in the Python standard library.
	maxComp := max(rn, gn, bn)
	minComp := min(rn, gn, bn)

	// Normalized value, 0...1
	vComp := maxComp
	if minComp == maxComp {
		return 0, 0, vComp
	}

	// Normalized saturation, 0...1
	sComp := (maxComp - minComp) / maxComp

	compRange := maxComp - minComp
	rMargin := (maxComp - rn) / compRange
	gMargin := (maxComp - gn) / compRange
	bMargin := (maxComp - bn) / compRange

	// Normalized hue, 0...1
	var hComp float64
	if rn == maxComp {
		hComp = bMargin - gMargin
	} else if gn == maxComp {
		hComp = 2 + rMargin - bMargin
	} else {
		hComp = 4 + gMargin - rMargin
	}
	if hComp < 0.0 {
		hComp += 6.0
	}
	hComp = math.Mod((hComp / 6.0), 1.0)

	return hComp, sComp, vComp
}

func RGB8ToHSV(r, g, b uint8) (float64, float64, float64) {
	r32 := uint32(r)
	// This is blindly copied from image/color/color.go.
	r32 |= r32 << 8
	g32 := uint32(g)
	g32 |= g32 << 8
	b32 := uint32(b)
	b32 |= b32 << 8
	return RGBToHSV(r32, g32, b32)
}

// Convert HSV color to RGB.
// h, s, v are each in 0.0 ... 1.0.
// The returned values are each in 0 ... 0xffffffff
func HSVToRGB(h, s, v float64) (uint32, uint32, uint32) {
	v32 := denorm(v)
	if s == 0 {
		return v32, v32, v32
	}
	i := math.Floor(h * 6.0)
	f := (h * 6.0) - float64(i) // nearest hue sextant
	p := denorm(v * (1.0 - s))
	q := denorm(v * (1.0 - s*f))
	t := denorm(v * (1.0 - s*(1.0-f)))

	switch uint32(math.Mod(i, 6.0)) {
	case 0:
		return v32, t, p
	case 1:
		return q, v32, p
	case 2:
		return p, v32, t
	case 3:
		return p, q, v32
	case 4:
		return t, p, v32
	case 5:
		return v32, p, q
	}
	// Should not be able to get here...
	return 0, 0, 0
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
