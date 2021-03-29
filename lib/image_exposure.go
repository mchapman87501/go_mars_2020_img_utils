package lib

import (
	"image"

	lib_image "dmoonc.com/mchapman87501/mars_2020_img_utils/lib/image"
)

// I *think* this is a more traditional exposure matcher than that of
// Compositor.matchColors.  The latter needs to exactly match subimages
// of the same scene.  This just tries to make the exposure (CIE Lab "L")
// histograms of two images look the same.

const numBuckets = int(10001)

type ImageExposure struct {
	Histogram []int     // Histogram of L values.
	PDF       []float64 // Probability density function
	CDF       []float64 // Cumulative density function
}

func NewImageExposure(image *lib_image.CIELab) *ImageExposure {
	result := &ImageExposure{
		Histogram: make([]int, numBuckets, numBuckets+1),
		PDF:       make([]float64, numBuckets, numBuckets+1),
		CDF:       make([]float64, numBuckets, numBuckets+1),
	}

	// map CIE Lab L range, 0...100, onto 0...lhBuckets-1
	scale := float64(numBuckets-1) / 100.0

	min := image.Bounds().Min
	max := image.Bounds().Max
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			lab := image.CIELabAt(x, y)
			index := int(lab.L * scale)
			result.Histogram[index] += 1
		}
	}

	// Compute the PDFs
	numPixels := image.Bounds().Dx() * image.Bounds().Dy()
	if numPixels > 0 {
		for i := 0; i < numBuckets; i++ {
			result.PDF[i] = float64(result.Histogram[i]) / float64(numPixels)
			prevCDF := 0.0
			if i > 0 {
				prevCDF = result.CDF[i-1]
			}
			result.CDF[i] = prevCDF + result.PDF[i]
		}
	}

	return result
}

// Get a copy of a ref image whose exposure is matched to that of a target image.
func MatchExposure(ref image.Image, target image.Image) image.Image {
	refLab := lib_image.CIELabFromImage(ref)
	refExposure := NewImageExposure(refLab)
	targetExposure := NewImageExposure(lib_image.CIELabFromImage(target))

	// map CIE Lab L range, 0...100, onto 0...lhBuckets-1
	scale := float64(numBuckets-1) / 100.0

	// Build a mapping from refExposure to targetExposure.
	labLMap := make(map[int]float64)
	iTarg := 0
	for iRef := 0; iRef < numBuckets; iRef++ {
		refCDF := refExposure.CDF[iRef]
		for (iTarg < numBuckets-1) && (targetExposure.CDF[iTarg] < refCDF) {
			iTarg += 1
		}
		labLMap[iRef] = float64(iTarg) / scale
	}

	result := lib_image.NewCIELab(ref.Bounds())
	min := result.Bounds().Min
	max := result.Bounds().Max
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			src := refLab.CIELabAt(x, y)
			key := int(src.L * scale)
			src.L = labLMap[key]
			result.SetCIELab(x, y, src)
		}
	}
	return result
}
