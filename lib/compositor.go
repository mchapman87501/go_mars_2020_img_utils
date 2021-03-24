package lib

import (
	"image"
	"image/draw"

	hsv_image "com.dmoonc/mchapman87501/mars_2020_img_utils/lib/image"
)

// Compositor builds a composite image from constituent tile
// images.
type Compositor struct {
	Bounds     image.Rectangle
	addedAreas []image.Rectangle
	Result     *hsv_image.HSV
}

func NewCompositor(rect image.Rectangle) Compositor {
	return Compositor{
		rect,
		[]image.Rectangle{},
		hsv_image.NewHSV(rect),
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
	H, S, V chanAdjustmentMap
}

func (comp *Compositor) matchColors(image image.Image, destRect image.Rectangle) image.Image {
	result := hsv_image.HSVFromImage(image)

	adjustments := comp.makeValueAdjustmentMap(result, destRect)
	adjustColors(result, adjustments)

	return result
}

func (comp *Compositor) makeValueAdjustmentMap(hsvImage *hsv_image.HSV, destRect image.Rectangle) *adjustmentMap {
	xTransform := destRect.Min.X - hsvImage.Bounds().Min.X
	yTransform := destRect.Min.Y - hsvImage.Bounds().Min.Y

	// Build an adjustment mapping, for all overlapping addedAreas, that
	// maps from result's pixel V channel to that of the corresponding
	// comp.Result pixel.
	h := make(chanAdjustmentMap)
	s := make(chanAdjustmentMap)
	v := make(chanAdjustmentMap)

	hCounts := make(cntMap)
	sCounts := make(cntMap)
	vCounts := make(cntMap)

	for _, rect := range comp.addedAreas {
		overlap := rect.Intersect(destRect)
		if !overlap.Empty() {
			// image.Bounds().Min.X corresponds to destRect.Min.X, and so on.
			// Need to convert coords.
			for x := overlap.Min.X; x < overlap.Max.X; x++ {
				for y := overlap.Min.Y; y < overlap.Max.Y; y++ {
					targetPix := comp.Result.HSVAt(x, y)
					srcPix := hsvImage.HSVAt(x-xTransform, y-yTransform)

					// Multiple srcPix-els may have the same V channel value,
					// and each one may map to a different targetPix V
					// channel value.
					// Use the average: sum
					// all mappings, then divide by the number of mappings.
					h[srcPix.H] += targetPix.H
					s[srcPix.S] += targetPix.S
					v[srcPix.V] += targetPix.V

					hCounts[srcPix.H] += 1
					sCounts[srcPix.S] += 1
					vCounts[srcPix.V] += 1
				}
			}
		}
	}

	// Overlapping regions may not cover the full gamut of channel values.
	// Add a default 1:1 mapping for extreme values, to aid interpolation.
	addExtrema(h, hCounts)
	addExtrema(s, sCounts)
	addExtrema(v, vCounts)

	normalizeCAM(h, hCounts)
	normalizeCAM(s, sCounts)
	normalizeCAM(v, vCounts)

	return &adjustmentMap{H: h, S: s, V: v}
}

func addExtrema(cam chanAdjustmentMap, counts cntMap) {
	// Overlapping regions may not cover the full gamut of channel values.
	// Add a default 1:1 mapping for extreme values, to aid interpolation.
	extrema := []float64{0.0, 1.0}
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

func adjustColors(hsvImage *hsv_image.HSV, adjustments *adjustmentMap) {
	hInterp := NewFloat64Interpolator(adjustments.H)
	sInterp := NewFloat64Interpolator(adjustments.S)
	vInterp := NewFloat64Interpolator(adjustments.V)

	for y := hsvImage.Bounds().Min.Y; y < hsvImage.Bounds().Max.Y; y++ {
		for x := hsvImage.Bounds().Min.X; x < hsvImage.Bounds().Max.X; x++ {
			pix := hsvImage.HSVAt(x, y)
			pix.H = hInterp.Interp(pix.H)
			pix.S = sInterp.Interp(pix.S)
			pix.V = vInterp.Interp(pix.V)
			hsvImage.SetHSV(x, y, pix)
		}
	}
}

// Ensure the range of Value component values lies within 0.0 ... 1.0
func (comp *Compositor) CompressDynamicRange() {
	// I think this is completely un-necessary...
	pixelStride := 3
	minValue := 1.0
	maxValue := 0.0
	for i := 2; i < len(comp.Result.Pix); i += pixelStride {
		currValue := comp.Result.Pix[i]
		if minValue > currValue {
			minValue = currValue
		}
		if maxValue < currValue {
			maxValue = currValue
		}
	}

	if (minValue < 0.0) || (maxValue > 1.0) {
		dValue := maxValue - minValue
		scale := 1.0 / dValue
		for i := 2; i < len(comp.Result.Pix); i += pixelStride {
			comp.Result.Pix[i] = (comp.Result.Pix[i] - minValue) * scale
		}
	}
}
