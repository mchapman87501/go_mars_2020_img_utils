package lib

import (
	"image"

	lib_image "github.com/mchapman87501/go_mars_2020_img_utils/lib/image"
)

// I *think* this is a more traditional exposure matcher than that of
// Compositor.matchColors.  The latter needs to exactly match subimages
// of the same scene.  This just tries to make the exposure (CIE Lab "L")
// histograms of two images look the same.

const numBuckets = int(10001)

type ImageExposure struct {
	BinMinVal float64   // The L value corresponding to the first bin.
	BinScale  float64   // Multiplier maps Lab L range to bucket index range
	Histogram []int     // Histogram of L values.
	CDF       []float64 // Cumulative density function
}

func NewImageExposure(image *lib_image.CIELab) *ImageExposure {
	result := &ImageExposure{
		BinMinVal: 0.0,
		BinScale:  1.0,
		Histogram: make([]int, numBuckets, numBuckets+1),
		CDF:       make([]float64, numBuckets, numBuckets+1),
	}

	// Input image may have an invalid Lab L range.
	// This can happen when building a composite image out of tiles that
	// have widely varying exposure ranges.
	// Map the actual range onto 0...lhBuckets-1
	min := image.Bounds().Min
	max := image.Bounds().Max

	first := true
	labLMin := 0.0
	labLMax := 100.0
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			lab := image.CIELabAt(x, y)
			if first || lab.L < labLMin {
				labLMin = lab.L
			}
			if first || lab.L > labLMax {
				labLMax = lab.L
			}
			first = false
		}
	}
	scale := float64(numBuckets-1) / (labLMax - labLMin)
	result.BinMinVal = labLMin
	result.BinScale = scale
	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			lab := image.CIELabAt(x, y)
			index := int((lab.L - labLMin) * scale)
			result.Histogram[index] += 1
		}
	}

	// Compute the PDFs
	numPixels := image.Bounds().Dx() * image.Bounds().Dy()
	if numPixels > 0 {
		for i := 0; i < numBuckets; i++ {
			pdf := float64(result.Histogram[i]) / float64(numPixels)
			prevCDF := 0.0
			if i > 0 {
				prevCDF = result.CDF[i-1]
			}
			result.CDF[i] = prevCDF + pdf
		}
	}

	return result
}

// Find the cumulative density bin that covers at least fract
// of all image pixels.  I.e., find the Lab L value such that the given
// fraction of pixelx are no brighter than that value.
func (e ImageExposure) cdf(fract float64) float64 {

	for i, cdf := range e.CDF {
		if cdf >= fract {
			return e.labL(i)
		}
	}
	// Should not get here, unless there is a fract domain error.
	if fract <= 0.0 {
		return 0
	}
	// Assume fract >= 1
	return e.labL(numBuckets - 1)
}

func (e ImageExposure) labL(binIndex int) float64 {
	return float64(binIndex)/e.BinScale + e.BinMinVal
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
