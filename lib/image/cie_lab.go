package image

import (
	"image"
	"image/color"

	lib_color "github.com/mchapman87501/go_mars_2020_img_utils/lib/image/color"
)

type CIELab struct {
	// Pix, Stride, Rect
	Pix    []float64
	Stride int
	Rect   image.Rectangle
}

func (p *CIELab) ColorModel() color.Model { return lib_color.CIELabModel }

func (p *CIELab) Bounds() image.Rectangle { return p.Rect }

func (p *CIELab) At(x, y int) color.Color {
	return p.CIELabAt(x, y)
}

func (p *CIELab) CIELabAt(x, y int) lib_color.CIELab {
	if !(image.Point{x, y}.In(p.Rect)) {
		return lib_color.CIELab{
			L: 0,
			A: 0,
			B: 0,
		}
	}
	i := p.PixOffset(x, y)
	result := lib_color.CIELab{L: p.Pix[i], A: p.Pix[i+1], B: p.Pix[i+2]}
	return result
}

func (p *CIELab) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *CIELab) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	c1 := lib_color.CIELabModel.Convert(c).(lib_color.CIELab)

	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3]
	s[0] = c1.L
	s[1] = c1.A
	s[2] = c1.B
}

func (p *CIELab) SetCIELab(x, y int, c lib_color.CIELab) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3]
	s[0] = c.L
	s[1] = c.A
	s[2] = c.B
}

func (p *CIELab) SubImage(rect image.Rectangle) *CIELab {
	// This is taken from the implementation of NRGBA's SubImage.
	r := rect.Intersect(p.Rect)
	if r.Empty() {
		return &CIELab{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &CIELab{
		Pix:    p.Pix[i:], // <- Those who choose to access Pix directly can overrun the image.
		Stride: p.Stride,
		Rect:   r,
	}
}

func (p *CIELab) Opaque() bool {
	return true
}

func NewCIELab(r image.Rectangle) *CIELab {
	area := r.Dx() * r.Dy()
	channels := 3 // L, a, b
	bufferSize := area * channels
	pix := make([]float64, bufferSize)
	return &CIELab{Pix: pix, Stride: channels * r.Dx(), Rect: r}
}

// Create a CIELab from an image.RGBA.
// The image origin will be at 0, 0.
func CIELabFromImage(src image.Image) *CIELab {
	rect := src.Bounds()
	result := NewCIELab(rect)

	offset := 0
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			labL, laba, labb := lib_color.RGBToCIELab(r, g, b)
			result.Pix[offset] = labL
			result.Pix[offset+1] = laba
			result.Pix[offset+2] = labb
			offset += 3
		}
	}
	return result
}
