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

// Add a new image.  Adjust its contrast range as necessary to match
// any overlapping image data that has already been composited.
func (comp *Compositor) AddImage(image image.Image, subframeRect image.Rectangle) {
	adjustedImage := comp.matchValue(image, subframeRect)
	// First draft: don't worry about contrast matching.  Just create the composite image.
	srcPoint := adjustedImage.Bounds().Min
	draw.Src.Draw(comp.Result, subframeRect, adjustedImage, srcPoint)
	comp.addedAreas = append(comp.addedAreas, subframeRect)
}

type valueAdjustmentMap map[float64]float64

func (comp *Compositor) matchValue(image image.Image, destRect image.Rectangle) image.Image {
	result := hsv_image.HSVFromImage(image)

	vam := comp.makeValueAdjustmentMap(result, destRect)
	adjustValues(result, vam)

	return result
}

func (comp *Compositor) makeValueAdjustmentMap(hsvImage *hsv_image.HSV, destRect image.Rectangle) valueAdjustmentMap {
	xTransform := destRect.Min.X - hsvImage.Bounds().Min.X
	yTransform := destRect.Min.Y - hsvImage.Bounds().Min.Y

	// Build an adjustment mapping, for all overlapping addedAreas, that
	// maps from result's pixel V channel to that of the corresponding
	// comp.Result pixel.
	result := make(valueAdjustmentMap)
	counts := make(map[float64]float64)

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
					result[srcPix.V] += targetPix.V
					counts[srcPix.V] += 1
				}
			}
		}
	}

	// Overlapping regions may not cover the full gamut of Values.
	// Add a default 1:1 mapping for extreme values, to aid interpolation.
	extrema := []float64{0.0, 1.0}
	for _, v := range extrema {
		_, ok := result[v]
		if !ok {
			result[v] = v
			counts[v] = 1
		}
	}

	for k := range result {
		result[k] /= counts[k]
	}

	return result
}

func adjustValues(hsvImage *hsv_image.HSV, mapping valueAdjustmentMap) {
	interpolator := NewFloat64Interpolator(mapping)
	for y := hsvImage.Bounds().Min.Y; y < hsvImage.Bounds().Max.Y; y++ {
		for x := hsvImage.Bounds().Min.X; x < hsvImage.Bounds().Max.X; x++ {
			pix := hsvImage.HSVAt(x, y)
			pix.V = interpolator.Interp(pix.V)
			if pix.V < 0 {
				pix.V = 0.0
			} else if pix.V > 1 {
				pix.V = 1.0
			}
			hsvImage.SetHSV(x, y, pix)
		}
	}
}

// Ensure the range of Value component values lies within 0.0 ... 1.0
func (comp *Compositor) CompressDynamicRange() {
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
