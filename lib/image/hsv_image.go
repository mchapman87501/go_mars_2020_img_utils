package image

import (
	"image"
	"image/color"

	hsv_color "com.dmoonc/mchapman87501/mars_2020_img_utils/lib/image/color"
)

type HSVImage struct {
	// Pix, Stride, Rect
	Pix    []float64
	Stride int
	Rect   image.Rectangle
}

func (p *HSVImage) ColorModel() color.Model { return hsv_color.HSVModel }

func (p *HSVImage) Bounds() image.Rectangle { return p.Rect }

func (p *HSVImage) At(x, y int) color.Color {
	return p.HSVAt(x, y)
}

func (p *HSVImage) HSVAt(x, y int) hsv_color.HSV {
	if !(image.Point{x, y}.In(p.Rect)) {
		return hsv_color.HSV{
			H: 0,
			S: 0,
			V: 0,
		}
	}
	i := p.PixOffset(x, y)
	return hsv_color.HSV{p.Pix[i], p.Pix[i+1], p.Pix[i+2]}
}

func (p *HSVImage) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
}

func (p *HSVImage) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := hsv_color.HSVModel.Convert(c).(hsv_color.HSV)
	s := p.Pix[i : i+3 : i+3]
	s[0] = c1.H
	s[1] = c1.S
	s[2] = c1.V
}

func (p *HSVImage) SetHSVAt(x, y int, c hsv_color.HSV) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+3 : i+3]
	s[0] = c.H
	s[1] = c.S
	s[2] = c.V
}

func (p *HSVImage) SubImage(rect image.Rectangle) *HSVImage {
	// This is taken from the implementation of NRGBA's SubImage.
	r := rect.Intersect(p.Rect)
	if r.Empty() {
		return &HSVImage{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &HSVImage{
		Pix:    p.Pix[i:], // <- Those who choose to access Pix directly can overrun the image.
		Stride: p.Stride,
		Rect:   r,
	}
}

func (p *HSVImage) Opaque() bool {
	return true
}

// Create an HSVImage from an image.RGBA.
// The image origin will be at 0, 0.
func HSVImageFromRGBA(rgba *image.RGBA) HSVImage {
	width := rgba.Rect.Dx()
	height := rgba.Rect.Dy()
	pixelSize := 3 // H,S,V
	numPixels := width * height
	rect := rgba.Rect
	result := HSVImage{make([]float64, numPixels*pixelSize), width * pixelSize, rect}

	srcOffset := 0
	destOffset := 0
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			srcSlice := rgba.Pix[srcOffset : srcOffset+4]
			destSlice := result.Pix[destOffset : destOffset+3]
			destSlice[0], destSlice[1], destSlice[2] = hsv_color.RGB8ToHSV(srcSlice[0], srcSlice[1], srcSlice[2])

			srcOffset += 4
			destOffset += 3
		}
	}
	return result
}
