package lib

import (
	"fmt"
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
	fmt.Println("Add composite at", subframeRect)
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

	fmt.Println("New image overlap regions:")
	for _, rect := range comp.addedAreas {
		overlap := rect.Intersect(destRect)
		if !overlap.Empty() {
			fmt.Println("   ", overlap)
			// image.Bounds().Min.X corresponds to destRect.Min.X, and so on.
			// Need to convert coords.
			fmt.Printf("  Corresponds to new image (%v, %v)\n", overlap.Min.X-xTransform, overlap.Min.Y-yTransform)
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
	fmt.Println("")

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
	counter := 0
	for y := hsvImage.Bounds().Min.Y; y < hsvImage.Bounds().Max.Y; y++ {
		for x := hsvImage.Bounds().Min.X; x < hsvImage.Bounds().Max.X; x++ {
			pix := hsvImage.HSVAt(x, y)
			vOrig := pix.V
			pix.V = interpolator.Interp(pix.V)
			if pix.V < 0 {
				pix.V = 0.0
			} else if pix.V > 1 {
				pix.V = 1.0
			}
			if counter < hsvImage.Bounds().Dx() {
				if vOrig < pix.V {
					fmt.Printf("Lighten %.4f -> %.4f\n", vOrig, pix.V)
				} else if vOrig > pix.V {
					fmt.Printf("Darken %.4f -> %.4f\n", vOrig, pix.V)
				} else {
					fmt.Println("No change")
				}
			}

			// Debug: desaturate:
			pix.H = 0.0
			pix.S = 0.0
			hsvImage.SetHSV(x, y, pix)
			counter += 1
		}
	}
}
