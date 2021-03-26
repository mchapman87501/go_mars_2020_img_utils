package lib

import (
	"image"
	"image/draw"

	lib_image "com.dmoonc/mchapman87501/mars_2020_img_utils/lib/image"
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

type chanAdjustmentMap map[float64]float64
type cntMap map[float64]float64

type adjustmentMap struct {
	L, A, B chanAdjustmentMap
}

func (comp *Compositor) matchColors(tileImage image.Image, destRect image.Rectangle) image.Image {
	result := lib_image.CIELabFromImage(tileImage)

	adjustments := comp.makeValueAdjustmentMap(result, destRect)
	adjustColors(result, adjustments)

	return result
}

func (comp *Compositor) makeValueAdjustmentMap(tileImage *lib_image.CIELab, destRect image.Rectangle) *adjustmentMap {
	xTransform := destRect.Min.X - tileImage.Bounds().Min.X
	yTransform := destRect.Min.Y - tileImage.Bounds().Min.Y

	// Build an adjustment mapping, for all overlapping addedAreas, that
	// maps from result's pixel V channel to that of the corresponding
	// comp.Result pixel.
	labL := make(chanAdjustmentMap)
	laba := make(chanAdjustmentMap)
	labb := make(chanAdjustmentMap)

	lCounts := make(cntMap)
	aCounts := make(cntMap)
	bCounts := make(cntMap)

	for _, rect := range comp.addedAreas {
		overlap := rect.Intersect(destRect)
		if !overlap.Empty() {
			// image.Bounds().Min.X corresponds to destRect.Min.X, and so on.
			// Need to convert coords.
			for x := overlap.Min.X; x < overlap.Max.X; x++ {
				for y := overlap.Min.Y; y < overlap.Max.Y; y++ {
					targetPix := comp.Result.CIELabAt(x, y)
					srcPix := tileImage.CIELabAt(x-xTransform, y-yTransform)

					// Multiple srcPix-els may have the same V channel value,
					// and each one may map to a different targetPix V
					// channel value.
					// Use the average: sum
					// all mappings, then divide by the number of mappings.
					labL[srcPix.L] += targetPix.L
					laba[srcPix.A] += targetPix.A
					labb[srcPix.B] += targetPix.B

					lCounts[srcPix.L] += 1
					aCounts[srcPix.A] += 1
					bCounts[srcPix.B] += 1
				}
			}
		}
	}

	// Overlapping regions may not cover the full gamut of channel values.
	// Add a default 1:1 mapping for extreme values, to aid interpolation.
	addExtrema(labL, lCounts, 0.0, 100.0)
	addExtrema(laba, aCounts, -128.0, 127.0)
	addExtrema(labb, bCounts, -128.0, 127.0)

	normalizeCAM(labL, lCounts)
	normalizeCAM(laba, aCounts)
	normalizeCAM(labb, bCounts)

	return &adjustmentMap{L: labL, A: laba, B: labb}
}

func addExtrema(cam chanAdjustmentMap, counts cntMap, minVal, maxVal float64) {
	// Overlapping regions may not cover the full gamut of channel values.
	// Add a default 1:1 mapping for extreme values, to aid interpolation.
	extrema := []float64{minVal, maxVal}
	for _, v := range extrema {
		_, ok := cam[v]
		if !ok {
			cam[v] = v
			counts[v] = 1
		}
	}
}

func normalizeCAM(cam chanAdjustmentMap, counts cntMap) {
	for k := range cam {
		cam[k] /= counts[k]
	}
}

func adjustColors(image *lib_image.CIELab, adjustments *adjustmentMap) {
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

// Try to "enhance" the result image.  This may mean adjusting L histograms,
// etc.
func (comp *Compositor) AutoEnhance() {

}

// Compress the dynamic range of the result image to make it more likely that
// it will fit within the sRGB color gamut.
// sRGB bounding cube, from image/color/print_srgb_gamut.py, is roughly
// {'labL': [0.0, 99.99998453333127], 'laba': [-86.1829494051608, 98.23532017664644], 'labb': [-107.86546414496824, 94.47731817969378]}

// type LabBounds struct {
// 	Min lib_color.CIELab
// 	Max lib_color.CIELab
// }

// func (comp *Compositor) CompressDynamicRange() {
// compressAsNeeded(getDynamicRange(comp.Result), comp.Result)
// }

// func getDynamicRange(image *lib_image.CIELab) LabBounds {
// 	min := lib_color.CIELab{100.0, 127.0, 127.0}
// 	max := lib_color.CIELab{0.0, -128.0, -128.0}

// 	for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
// 		for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
// 			pix := image.CIELabAt(x, y)

// 			if pix.L > max.L {
// 				max.L = pix.L
// 			}
// 			if pix.L < min.L {
// 				min.L = pix.L
// 			}

// 			if pix.A > max.A {
// 				max.A = pix.A
// 			}
// 			if pix.A < min.A {
// 				min.A = pix.A
// 			}

// 			if pix.B > max.B {
// 				max.B = pix.B
// 			}
// 			if pix.B < min.B {
// 				min.B = pix.B
// 			}
// 		}
// 	}

// 	return LabBounds{Min: min, Max: max}
// }

// type channelGetter func(p lib_color.CIELab) float64
// type channelSetter func(p lib_color.CIELab, v float64)

// type scaler struct {
// 	MinIn, MinOut, Scale float64

// 	getChan channelGetter
// 	setChan channelSetter
// }

// func newScaler(
// 	minIn, minOut, maxIn, maxOut float64,
// 	getChan channelGetter, setChan channelSetter) scaler {
// 	// Is the input already in bounds?
// 	if (minIn >= minOut) && (maxIn <= maxOut) {
// 		return scaler{minIn, minIn, 1.0, nil, nil}
// 	}

// 	dIn := maxIn - minIn
// 	dOut := maxOut - minOut
// 	if dIn <= 0.0 {
// 		return scaler{minIn, minIn, 1.0, nil, nil}
// 	}

// 	scale := dOut / dIn
// 	return scaler{minIn, minOut, scale, getChan, setChan}
// }

// func (s *scaler) UpdatePix(p lib_color.CIELab) {
// 	v := s.getChan(p)
// 	vNew := (v-s.MinIn)*s.Scale + s.MinOut
// 	s.setChan(p, vNew)
// }

// func (s *scaler) isIdentity() bool {
// 	dMin := math.Abs(s.MinOut - s.MinIn)

// 	return (dMin <= 1.0e-6) && (math.Abs(s.Scale-1.0) <= 1.0e-6)
// }

// func compressAsNeeded(imageRange LabBounds, image *lib_image.CIELab) {
// fmt.Println("Image band ranges:", imageRange)
// // {'labL': [0.0, 99.99998453333127], 'laba': [-86.1829494051608, 98.23532017664644], 'labb': [-107.86546414496824, 94.47731817969378]}
// minRGB := lib_color.CIELab{L: 0.0, A: -86.0, B: -107.0}
// maxRGB := lib_color.CIELab{L: 100.0, A: 98.0, B: 94.0}

// scalers := []scaler{}

// // getL := func(p lib_color.CIELab) float64 { return p.L }
// // setL := func(p lib_color.CIELab, v float64) { p.L = v }
// // lScaler := newScaler(imageRange.Min.L, minRGB.L, imageRange.Max.L, maxRGB.L, getL, setL)
// // if !lScaler.isIdentity() {
// // 	fmt.Println("Compress Lab L")
// // 	scalers = append(scalers, lScaler)
// // }

// getA := func(p lib_color.CIELab) float64 { return p.A }
// setA := func(p lib_color.CIELab, v float64) { p.A = v }
// aScaler := newScaler(imageRange.Min.A, minRGB.A, imageRange.Max.A, maxRGB.A, getA, setA)
// if !aScaler.isIdentity() {
// 	fmt.Println("Compress Lab a")
// 	scalers = append(scalers, aScaler)
// }

// getB := func(p lib_color.CIELab) float64 { return p.B }
// setB := func(p lib_color.CIELab, v float64) { p.B = v }
// bScaler := newScaler(imageRange.Min.B, minRGB.B, imageRange.Max.B, maxRGB.B, getB, setB)
// if !bScaler.isIdentity() {
// 	fmt.Println("Compress Lab a")
// 	scalers = append(scalers, bScaler)
// }
// compressChannels(image, scalers)
// }

// func compressChannels(image *lib_image.CIELab, scalers []scaler) {
// 	if len(scalers) > 0 {
// 		for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
// 			for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
// 				pix := image.CIELabAt(x, y)
// 				for _, s := range scalers {
// 					s.UpdatePix(pix)
// 				}
// 				image.SetCIELab(x, y, pix)
// 			}
// 		}
// 	}
// }
