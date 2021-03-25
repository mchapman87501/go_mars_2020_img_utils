package color

import (
	"image/color"
	"math"
)

// CIELab representss the CIE L*a*b* colorspace, 2° observer, D65 illuminant.
// See https://en.wikipedia.org/wiki/CIELAB_color_space#Forward_transformation
// Apologies for the uppercase member names ;)
type CIELab struct {
	L, A, B float64
	// TODO consider providing a pass-through alpha component.
}

// Get a CIE Lab color as RGBA.  This assumes
// sRGB ("Standard" RGB).
func (c CIELab) RGBA() (r, g, b, a uint32) {
	r, g, b = cieLabToRGB(c.L, c.A, c.B)
	a = 0xffff
	return
}

func cieLabModel(c color.Color) color.Color {
	if _, ok := c.(CIELab); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	labL, laba, labb := rgbToCIELab(r, g, b)
	return CIELab{L: labL, A: laba, B: labb}
}

var CIELabModel color.Model = color.ModelFunc(cieLabModel)

func cieLabToRGB(labL, laba, labb float64) (r, g, b uint32) {
	x, y, z := cieLabD65ToXYZ(labL, laba, labb)
	r, g, b = cieXYZToRGB(x, y, z)
	return
}

func rgbToCIELab(r, g, b uint32) (labL, laba, labb float64) {
	x, y, z := rgbToCIEXYZ(r, g, b)
	labL, laba, labb = cieXYZToLabD65(x, y, z)
	return
}

const cieE = 216.0 / 24389.0

const cieDelta float64 = 6.0 / 29.0
const cieD3 = cieDelta * cieDelta * cieDelta

// Caveat: these match python colormath, rather than Wikipedia.
const d65IllumX = 0.950489
const d65IllumY = 1.0
const d65IllumZ = 1.08884

func cieLabD65ToXYZ(labL, laba, labb float64) (x, y, z float64) {
	ty := (labL + 16.0) / 116.0
	tx := laba/500.0 + ty
	tz := ty - labb/200.0

	x = d65IllumX * scaleLabD65ToXYZ(tx)
	y = d65IllumY * scaleLabD65ToXYZ(ty)
	z = d65IllumZ * scaleLabD65ToXYZ(tz)
	return
}

func cieXYZToLabD65(x, y, z float64) (labL, laba, labb float64) {
	xIllum := x / d65IllumX
	yIllum := y / d65IllumY
	zIllum := z / d65IllumZ

	xScaled := scaleXYZToLabD65(xIllum)
	yScaled := scaleXYZToLabD65(yIllum)
	zScaled := scaleXYZToLabD65(zIllum)

	labL = 116.0*yScaled - 16.0
	laba = 500.0 * (xScaled - yScaled)
	labb = 200.0 * (yScaled - zScaled)
	return
}

func scaleLabD65ToXYZ(component float64) float64 {
	if component > cieDelta {
		return component * component * component
	}
	return 3.0 * cieD3 * (component - 4.0/29.0)
}

func scaleXYZToLabD65(component float64) float64 {
	if component > cieE {
		return math.Pow(component, 1.0/3.0)
	}
	return (7.787 * component) + (16.0 / 116.0)
}

// Get the CIE 2000 ΔE between 2 CIELab colors.
// A return value < 2.3 indicates a human probably would
// not notice any difference in the two colors.
func DeltaECIE2000(lab1, lab2 CIELab) float64 {
	return math.Sqrt(DeltaESqrCIE2000(lab1, lab2))
}

func DeltaESqrCIE2000(lab1, lab2 CIELab) float64 {
	// I have no idea whether this is correct, or where it
	// may fail with NaNs :)
	avgLp := (lab1.L + lab2.L) / 2.0

	b1Sqr := lab1.B * lab1.B
	b2Sqr := lab2.B * lab2.B
	c1 := math.Sqrt(lab1.A*lab1.A + b1Sqr)
	c2 := math.Sqrt(lab2.A*lab2.A + b2Sqr)
	avgC := (c1 + c2) / 2.0
	avgC7 := math.Pow(avgC, 7.0)

	pow25_7 := math.Pow(25.0, 7.0)

	g := (1.0 - math.Sqrt(avgC7/(avgC7+pow25_7))) / 2.0

	a1p := lab1.A * (1.0 + g)
	a2p := lab2.A * (1.0 + g)

	c1p := math.Sqrt(a1p*a1p + b1Sqr)
	c2p := math.Sqrt(a2p*a2p + b2Sqr)
	avgCp := (c1p + c2p) / 2.0

	h1p := posRadToDeg(math.Atan2(lab1.B, a1p))
	h2p := posRadToDeg(math.Atan2(lab2.B, a2p))
	avgHp := (h1p + h2p) / 2.0
	if math.Abs(h1p-h2p) > 180.0 {
		avgHp += 180.0
	}

	t := (1.0 -
		0.17*math.Cos(degToRad(avgHp-30.0)) +
		0.24*math.Cos(degToRad(2.0*avgHp)) +
		0.32*math.Cos(degToRad(3.0*avgHp+6.0)) -
		0.2*math.Cos(degToRad(4.0*avgHp-63.0)))

	deltaLp := lab2.L - lab1.L
	deltaCp := c2p - c1p

	deltaHp := h2p - h1p
	if math.Abs(deltaHp) > 180.0 {
		if h2p <= h1p {
			deltaHp += 360.0
		} else {
			deltaHp -= 360.0
		}
	}
	deltaHp = 2.0 * math.Sqrt(c1p*c2p) * math.Sin(degToRad(deltaHp)/2.0)

	avgLPBias := avgLp - 50.0
	aLPBSqr := avgLPBias * avgLPBias

	sL := 1.0 + ((0.015 * aLPBSqr) / math.Sqrt(20.0+aLPBSqr))
	sC := 1.0 + 0.045*avgCp
	sH := 1.0 + 0.015*avgCp*t

	aHp := (avgHp - 275.0) / 25.0
	aHpSqr := aHp * aHp
	deltaRo := 30.0 * math.Exp(-aHpSqr)

	avgCp7 := math.Pow(avgCp, 7.0)

	rC := 2.0 * math.Sqrt(avgCp7/(avgCp7+pow25_7))
	rT := -rC * math.Sin(2.0*degToRad(deltaRo))

	kL := 1.0
	kC := 1.0
	kH := 1.0

	t0 := deltaLp / (sL * kL)
	t1 := deltaCp / (sC * kC)
	t2 := deltaHp / (sH * kH)
	return (t0 * t0) + (t1 * t1) + (t2 * t2) + rT*t1*t2
}

func posRadToDeg(radians float64) float64 {
	result := radToDeg(radians)
	if result < 0.0 {
		return result + 360.0
	}
	return result
}

func radToDeg(radians float64) float64 {
	return 180.0 * radians / math.Pi
}

func degToRad(degrees float64) float64 {
	return math.Pi * degrees / 180.0
}
