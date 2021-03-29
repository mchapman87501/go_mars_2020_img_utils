package lib

import (
	"image"
	"image/png"
	"os"
)

func SavePNG(image image.Image, filename string) error {
	outf, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outf.Close()

	return png.Encode(outf, image)
}
