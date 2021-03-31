package lib

import (
	"image"
	"image/draw"
	"math"

	lib_image "github.com/mchapman87501/go_mars_2020_img_utils/lib/image"
	lib_color "github.com/mchapman87501/go_mars_2020_img_utils/lib/image/color"
)

// Compositor builds a composite image from constituent tile
// images.
type Compositor struct {
	Bounds     image.Rectangle
	addedAreas []image.Rectangle
	Result     *lib_image.CIELab
}

func NewCompositor(rect image.Rectangle) Compositor {
	return Compositor{
		rect,
		[]image.Rectangle{},
		lib_image.NewCIELab(rect),
	}
}

// Add a new image.  Adjust its colors as necessary to match
// any overlapping image data that has already been composited.
func (comp *Compositor) AddImage(image image.Image, subframeRect image.Rectangle) {
	adjustedImage := comp.matchColors(image, subframeRect)
	srcPoint := adjustedImage.Bounds().Min
	draw.Src.Draw(comp.Result, subframeRect, adjustedImage, srcPoint)
	comp.addedAreas = append(comp.addedAreas, subframeRect)
}

func (comp *Compositor) matchColors(tileImage image.Image, destRect image.Rectangle) image.Image {
	result := lib_image.CIELabFromImage(tileImage)

	adjustments := comp.makeValueAdjustmentMap(result, destRect)
	adjustColors(result, adjustments)

	return result
}

func (comp *Compositor) makeValueAdjustmentMap(tileImage *lib_image.CIELab, destRect image.Rectangle) *AdjustmentMap {
	// What is the tile image's origin in composite image coordinates?
	tileOrigin := tileImage.Bounds().Min
	// 'translate' the tile to destRect.
	tileOffset := destRect.Bounds().Min.Sub(tileOrigin)

	// Build an adjustment mapping, for all overlapping addedAreas, that
	// maps from result's pixel V channel to that of the corresponding
	// comp.Result pixel.
	result := NewAdjustmentMap()

	for _, rect := range comp.addedAreas {
		overlap := rect.Intersect(destRect)
		if !overlap.Empty() {
			for x := overlap.Min.X; x < overlap.Max.X; x++ {
				for y := overlap.Min.Y; y < overlap.Max.Y; y++ {
					srcPix := tileImage.CIELabAt(x-tileOffset.X, y-tileOffset.Y)
					targetPix := comp.Result.CIELabAt(x, y)
					result.AddSample(srcPix, targetPix)
				}
			}
		}
	}
	result.Complete()
	return result
}

func adjustColors(image *lib_image.CIELab, adjustments *AdjustmentMap) {
	lInterp := NewFloat64Interpolator(adjustments.L)
	aInterp := NewFloat64Interpolator(adjustments.A)
	bInterp := NewFloat64Interpolator(adjustments.B)

	for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
		for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
			pix := image.CIELabAt(x, y)
			pix.L = lInterp.Interp(pix.L)
			pix.A = aInterp.Interp(pix.A)
			pix.B = bInterp.Interp(pix.B)
			image.SetCIELab(x, y, pix)
		}
	}
}

// Compress the dynamic range of the result image to make it more likely that
// it will fit within the sRGB color gamut.
// sRGB bounding cube, from image/color/print_srgb_gamut.py, is roughly
// {'labL': [0.0, 99.99998453333127], 'laba': [-86.1829494051608, 98.23532017664644], 'labb': [-107.86546414496824, 94.47731817969378]}

type LabBounds struct {
	Min lib_color.CIELab
	Max lib_color.CIELab
}

func (comp *Compositor) CompressDynamicRange() {
	compressAsNeeded(getDynamicRange(comp.Result), comp.Result)
}

func getDynamicRange(image *lib_image.CIELab) LabBounds {
	min := lib_color.CIELab{L: 100.0, A: 127.0, B: 127.0}
	max := lib_color.CIELab{L: 0.0, A: -128.0, B: -128.0}

	for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
		for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
			pix := image.CIELabAt(x, y)

			if pix.A > max.A {
				max.A = pix.A
			}
			if pix.A < min.A {
				min.A = pix.A
			}

			if pix.B > max.B {
				max.B = pix.B
			}
			if pix.B < min.B {
				min.B = pix.B
			}
		}
	}

	// Special case for Lab L:  Adjust exposure so that some large fraction
	// of pixels are within gamut.
	exposure := NewImageExposure(image)
	min.L = exposure.cdf(0.0)
	max.L = exposure.cdf(0.95)
	return LabBounds{Min: min, Max: max}
}

type channelGetter func(p *lib_color.CIELab) float64
type channelSetter func(p *lib_color.CIELab, v float64)

type scaler struct {
	MinIn, MinOut, Scale float64

	getChan channelGetter
	setChan channelSetter
}

func newScaler(
	minIn, minOut, maxIn, maxOut float64,
	getChan channelGetter, setChan channelSetter) scaler {
	// Is the input already in bounds?
	if (minIn >= minOut) && (maxIn <= maxOut) {
		return scaler{minIn, minIn, 1.0, nil, nil}
	}

	dIn := maxIn - minIn
	dOut := maxOut - minOut
	if dIn <= 0.0 {
		return scaler{minIn, minIn, 1.0, nil, nil}
	}

	scale := dOut / dIn
	return scaler{minIn, minOut, scale, getChan, setChan}
}

func (s *scaler) UpdatePix(p *lib_color.CIELab) {
	v := s.getChan(p)
	vNew := (v-s.MinIn)*s.Scale + s.MinOut
	s.setChan(p, vNew)
}

func (s *scaler) isIdentity() bool {
	dMin := math.Abs(s.MinOut - s.MinIn)

	return (dMin <= 1.0e-6) && (math.Abs(s.Scale-1.0) <= 1.0e-6)
}

// Compress dynamic range as needed, to fit sRGB gamut.
// NOTE that this naive implementation is sensitive to outliers.
// Probably better to use a method, for L channel at least,
// that ensures some fraction f of pixels are unclipped.
func compressAsNeeded(imageRange LabBounds, image *lib_image.CIELab) {
	// {'labL': [0.0, 99.99998453333127], 'laba': [-86.1829494051608, 98.23532017664644], 'labb': [-107.86546414496824, 94.47731817969378]}
	minRGB := lib_color.CIELab{L: 0.0, A: -86.0, B: -107.0}
	maxRGB := lib_color.CIELab{L: 100.0, A: 98.0, B: 94.0}

	scalers := []scaler{}

	getL := func(p *lib_color.CIELab) float64 { return p.L }
	setL := func(p *lib_color.CIELab, v float64) { p.L = v }
	lScaler := newScaler(imageRange.Min.L, minRGB.L, imageRange.Max.L, maxRGB.L, getL, setL)
	if !lScaler.isIdentity() {
		scalers = append(scalers, lScaler)
	}

	getA := func(p *lib_color.CIELab) float64 { return p.A }
	setA := func(p *lib_color.CIELab, v float64) { p.A = v }
	aScaler := newScaler(imageRange.Min.A, minRGB.A, imageRange.Max.A, maxRGB.A, getA, setA)
	if !aScaler.isIdentity() {
		scalers = append(scalers, aScaler)
	}

	getB := func(p *lib_color.CIELab) float64 { return p.B }
	setB := func(p *lib_color.CIELab, v float64) { p.B = v }
	bScaler := newScaler(imageRange.Min.B, minRGB.B, imageRange.Max.B, maxRGB.B, getB, setB)
	if !bScaler.isIdentity() {
		scalers = append(scalers, bScaler)
	}
	compressChannels(image, scalers)
}

func compressChannels(image *lib_image.CIELab, scalers []scaler) {
	if len(scalers) > 0 {
		for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
			for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
				pix := image.CIELabAt(x, y)
				for _, s := range scalers {
					s.UpdatePix(&pix)
				}
				image.SetCIELab(x, y, pix)
			}
		}
	}
}
