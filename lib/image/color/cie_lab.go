package color

import (
	"fmt"
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
	fmt.Printf("rgb(0x%x, 0x%x, 0x%x) -> XYZ(%v, %v, %v)\n", r, g, b, x, y, z)
	labL, laba, labb = cieXYZToLabD65(x, y, z)
	fmt.Printf("     -> Lab(%v, %v, %v)\n", labL, laba, labb)
	return
}

const cieE = 216.0 / 24389.0

const cieDelta float64 = 6.0 / 29.0
const cieD3 = cieDelta * cieDelta * cieDelta

// Caveat: these match python colormath, rather than Wikipedia.
const d65IllumX = 0.95047
const d65IllumY = 1.0
const d65IllumZ = 1.08883

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
	fmt.Printf("     -> xyzIllum(%v, %v, %v)\n", xIllum, yIllum, zIllum)

	xScaled := scaleXYZToLabD65(xIllum)
	yScaled := scaleXYZToLabD65(yIllum)
	zScaled := scaleXYZToLabD65(zIllum)

	fmt.Printf("     -> xyzScaled(%v, %v, %v)\n", xScaled, yScaled, zScaled)

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
