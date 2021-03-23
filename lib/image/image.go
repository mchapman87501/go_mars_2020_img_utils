package image

import (
	"image"
	"image/color"

	hsv_color "com.dmoonc/mchapman87501/mars_2020_img_utils/lib/image/color"
)

type HSV struct {
	// Pix, Stride, Rect
	Pix    []float64
	Stride int
	Rect   image.Rectangle
}

func (p *HSV) ColorModel() color.Model { return hsv_color.HSVModel }

func (p *HSV) Bounds() image.Rectangle { return p.Rect }

func (p *HSV) At(x, y int) color.Color {
	return p.HSVAt(x, y)
}

func (p *HSV) HSVAt(x, y int) hsv_color.HSV {
	if !(image.Point{x, y}.In(p.Rect)) {
		return hsv_color.HSV{
			H: 0,
			S: 0,
			V: 0,
		}
	}
	i := p.PixOffset(x, y)
	result := hsv_color.HSV{H: p.Pix[i], S: p.Pix[i+1], V: p.Pix[i+2]}
	return result
}

func (p *HSV) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *HSV) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	c1 := hsv_color.HSVModel.Convert(c).(hsv_color.HSV)

	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3]
	s[0] = c1.H
	s[1] = c1.S
	s[2] = c1.V
}

func (p *HSV) SetHSV(x, y int, c hsv_color.HSV) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3]
	s[0] = c.H
	s[1] = c.S
	s[2] = c.V
}

func (p *HSV) SubImage(rect image.Rectangle) *HSV {
	// This is taken from the implementation of NRGBA's SubImage.
	r := rect.Intersect(p.Rect)
	if r.Empty() {
		return &HSV{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &HSV{
		Pix:    p.Pix[i:], // <- Those who choose to access Pix directly can overrun the image.
		Stride: p.Stride,
		Rect:   r,
	}
}

func (p *HSV) Opaque() bool {
	return true
}

func NewHSV(r image.Rectangle) *HSV {
	area := r.Dx() * r.Dy()
	channels := 3
	bufferSize := area * channels
	pix := make([]float64, bufferSize)
	return &HSV{Pix: pix, Stride: channels * r.Dx(), Rect: r}
}

// Create an HSV from an image.RGBA.
// The image origin will be at 0, 0.
func HSVFromImage(src image.Image) *HSV {
	rect := src.Bounds()
	result := NewHSV(rect)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			result.Set(x, y, src.At(x, y))
		}
	}
	return result
}
