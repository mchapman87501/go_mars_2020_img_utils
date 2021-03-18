package lib

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func savePNG(image image.Image, filename string, t *testing.T) {
	outf, err := os.Create(filename)
	if err != nil {
		t.Fatal("Could not create image file", filename, ":", err)
	}
	defer outf.Close()

	err = png.Encode(outf, image)
	if err != nil {
		t.Fatal("Error encoding image to PNG:", err)
	}
}

func TestGraydient(t *testing.T) {
	// Verify that a constant-tone grayscale image can be demosaiced
	// without crashing.

	width := 255
	height := 255
	imageRect := image.Rect(0, 0, width, height)
	grayImage := image.NewGray(imageRect)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			intensity := (x + y) / 2
			if (x%2 == 0) && (y%2 == 0) {
				intensity += 50 // Boost the red
			}
			if intensity > 255 {
				intensity = 255
			}
			grayImage.Set(x, y, color.Gray{uint8(intensity)})
		}
	}

	savePNG(grayImage, "test_data/out/input_test_graydient.png", t)

	rgbImage, err := Demosaic(grayImage)
	if err != nil {
		t.Fatal("Error de-mosaicing:", err)
	}

	savePNG(rgbImage, "test_data/out/result_test_graydient.png", t)
}
