package lib

import (
	"errors"
	"image"
	"image/color"
)

// Offsets of channels within RGBA pixels:
const r = 0
const g = 1
const b = 2

// filter pattern, addressable as [y][x]
var pattern = [][]int{
	{r, g},
	{g, b},
}

// Simple-minded, I am.  Generate indices of "neighbors" of a value at a given index.
func getIndices(index int, minIndex int, maxIndex int) []int {
	result := []int{}

	for currIndex := index - 1; currIndex <= (index + 1); currIndex++ {
		if (minIndex <= currIndex) && (currIndex < maxIndex) {
			result = append(result, currIndex)
		}
	}
	return result
}

// Given a bayer RGGB filter pattern overlaid on a sensor,
// what filter color (channel) lies over a given pixel?
func getChannel(x int, y int) int {
	return pattern[y%2][x%2]
}

func pixCompAvg(sum int, count int) uint8 {
	divisor := 1
	if count > 1 {
		divisor = count
	}
	return uint8(sum / divisor)
}

// Demosaic an RGGB bayer image.
func Demosaic(bayerImage *image.Gray) (image.Image, error) {
	// The most obvious way to do this is via convolution kernels,
	// but I'm finding it tedious to declare those.
	bounds := bayerImage.Bounds()
	result := image.NewRGBA(bounds)

	for xImage := bounds.Min.X; xImage < bounds.Max.X; xImage++ {
		xIndices := getIndices(xImage, bounds.Min.X, bounds.Max.X)
		for yImage := bounds.Min.Y; yImage < bounds.Max.Y; yImage++ {
			yIndices := getIndices(yImage, bounds.Min.Y, bounds.Max.Y)

			sums := []int{0, 0, 0}
			count := []int{0, 0, 0}
			for _, x := range xIndices {
				for _, y := range yIndices {
					channel := getChannel(x-bounds.Min.X, y-bounds.Min.Y)
					sums[channel] += int(bayerImage.GrayAt(x, y).Y)
					count[channel] += 1
				}
			}

			result.Set(xImage, yImage, color.RGBA{
				// Avoid divide-by-zero:
				pixCompAvg(sums[r], count[r]),
				pixCompAvg(sums[g], count[g]),
				pixCompAvg(sums[b], count[b]),
				0xff,
			})
		}
	}
	return result, nil
}

// Demosaic a grayscale image that was stored as RGBA.
func DemosaicRGBGray(bayerImage image.Image) (image.Image, error) {

	rgbaImage, ok := (bayerImage).(*image.RGBA)
	if !ok {
		return bayerImage, errors.New("bayerImage must be an RGBA image")
	}

	grayBayer := image.NewGray(rgbaImage.Bounds())

	// It should suffice to copy out any channel - all channels should
	// have the same values.
	iDest := 0
	for iSrc := 0; iSrc < len(rgbaImage.Pix); iSrc += 4 {
		grayBayer.Pix[iDest] = rgbaImage.Pix[iSrc]
		iDest += 1
	}
	return Demosaic(grayBayer)
}
