package color

import "math"

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

// Convert a value in the range 0...255 to a value in 0.0...1.0:
func norm(v uint8) float64 {
	return float64(v) / 255.0
}

// Convert a value in 0.0 ... 1.0 to a value in 0.0 ... 255.0
func denorm(v float64) uint8 {
	return uint8(v * 255.0)
}

// Convert RGB color to HSV.
// Returns values in the range 0.0 ... 1.0
func RGBToHSV(r, g, b uint8) (float64, float64, float64) {
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

// Convert HSV color to RGB.
// h, s, v are each in 0.0 ... 1.0.
// The returned values are each in 0 ... 255
func HSVToRGB(h, s, v float64) (uint8, uint8, uint8) {
	v8 := denorm(v)
	if s == 0 {
		return v8, v8, v8
	}
	i := math.Floor(h * 6.0)
	f := (h * 6.0) - float64(i) // nearest hue sextant
	p := denorm(v * (1.0 - s))
	q := denorm(v * (1.0 - s*f))
	t := denorm(v * (1.0 - s*(1.0-f)))

	switch uint8(math.Mod(i, 6.0)) {
	case 0:
		return v8, t, p
	case 1:
		return q, v8, p
	case 2:
		return p, v8, t
	case 3:
		return p, q, v8
	case 4:
		return t, p, v8
	case 5:
		return v8, p, q
	}
	// Should not be able to get here...
	return 0, 0, 0
}
