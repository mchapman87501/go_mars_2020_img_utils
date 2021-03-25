package color

import (
	"fmt"
	"image/color"
	"math"
)

// CIEXYZ represents the CIE 1931 XYZ color space.
// For background on the experiments that produced the
// color space definition (as well as CIE RGB), see
// https://en.wikipedia.org/wiki/CIE_1931_color_space

type CIEXYZ struct {
	X, Y, Z float64
	// TODO consider also storing alpha, as a "pass-through" value.
}

// Get a CIE XYZ color as RGBA.  This assumes
// sRGB ("Standard" RGB).
func (c CIEXYZ) RGBA() (r, g, b, a uint32) {
	r, g, b = cieXYZToRGB(c.X, c.Y, c.Z)
	a = 0xffff
	return
}

func cieXYZModel(c color.Color) color.Color {
	if _, ok := c.(CIEXYZ); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	x, y, z := rgbToCIEXYZ(r, g, b)
	return CIEXYZ{X: x, Y: y, Z: z}
}

var CIEXYZModel color.Model = color.ModelFunc(cieXYZModel)

// For the conversion used here see
// https://en.wikipedia.org/wiki/SRGB

func gammaExpanded(u float64) float64 {
	if u <= 0.04045 {
		return 25.0 * u / 323.0
	}
	uScaled := (u + 0.055) / 1.055
	return math.Pow(uScaled, 2.4)
}

func rgbToCIEXYZ(r, g, b uint32) (x, y, z float64) {
	nr := gammaExpanded(norm(r))
	ng := gammaExpanded(norm(g))
	nb := gammaExpanded(norm(b))

	x = 0.41239080*nr + 0.35758434*ng + 0.18048079*nb
	y = 0.21263901*nr + 0.71516868*ng + 0.07219232*nb
	z = 0.01933082*nr + 0.11919478*ng + 0.95053215*nb

	return
}

func gammaCompressed(u float64) float64 {
	result := 0.0

	if u <= 0.0031308 {
		result = 12.92 * u
	} else {
		result = 1.055*math.Pow(u, 1/2.4) - 0.055
	}
	if result < 0.0 {
		result = 0.0
	} else if result > 1.0 {
		result = 1.0
	}
	return result
}

func cieXYZToRGB(x, y, z float64) (r, g, b uint32) {
	nr := 3.24096994*x + -1.53738318*y + -0.49861076*z
	ng := -0.96924364*x + 1.8759675*y + 0.04155506*z
	nb := 0.05563008*x + -0.20397696*y + 1.05697151*z

	gr := gammaCompressed(nr)
	gg := gammaCompressed(ng)
	gb := gammaCompressed(nb)
	// fmt.Printf("xyztorgb(%v, %v, %v) -> (%v, %v, %v) -> (%v, %v, %v)\n", x, y, z, nr, ng, nb, gr, gg, gb)

	r = denorm(clip(gr))
	g = denorm(clip(gg))
	b = denorm(clip(gb))
	return
}

func clip(rgbComp float64) float64 {
	if rgbComp < 0.0 {
		fmt.Println("OOG < 0:", rgbComp)
		return 0.0
	}
	if rgbComp > 1.0 {
		fmt.Println("OOG > 1:", rgbComp)
		return 1.0
	}
	return rgbComp
}
