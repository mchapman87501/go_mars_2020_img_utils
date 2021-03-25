package lib

import (
	"image"
	"image/draw"

	mars_lib_image "com.dmoonc/mchapman87501/mars_2020_img_utils/lib/image"
)

// Compositor builds a composite image from constituent tile
// images.
type Compositor struct {
	Bounds     image.Rectangle
	addedAreas []image.Rectangle
	Result     *mars_lib_image.CIELab
}

func NewCompositor(rect image.Rectangle) Compositor {
	return Compositor{
		rect,
		[]image.Rectangle{},
		mars_lib_image.NewCIELab(rect),
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
	result := mars_lib_image.CIELabFromImage(tileImage)

	adjustments := comp.makeValueAdjustmentMap(result, destRect)
	adjustColors(result, adjustments)

	return result
}

func (comp *Compositor) makeValueAdjustmentMap(tileImage *mars_lib_image.CIELab, destRect image.Rectangle) *adjustmentMap {
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

func adjustColors(image *mars_lib_image.CIELab, adjustments *adjustmentMap) {
	hInterp := NewFloat64Interpolator(adjustments.L)
	sInterp := NewFloat64Interpolator(adjustments.A)
	vInterp := NewFloat64Interpolator(adjustments.B)

	for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
		for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
			pix := image.CIELabAt(x, y)
			pix.L = hInterp.Interp(pix.L)
			pix.A = sInterp.Interp(pix.A)
			pix.B = vInterp.Interp(pix.B)
			image.SetCIELab(x, y, pix)
		}
	}
}
