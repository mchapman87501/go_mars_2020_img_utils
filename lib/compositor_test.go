package lib

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"testing"

	hsv_image "github.com/mchapman87501/go_mars_2020_img_utils/lib/image"
)

func makeTile(row int, col int, width int, height int) draw.Image {
	rect := image.Rect(0, 0, width, height)
	result := image.NewRGBA(rect)
	for x := 0; x < width; x++ {
		intensity := (x + col*width) % 0xff
		red := (intensity + 50)
		if red > 0xff {
			red = 0xff
		}
		intens8 := uint8(intensity)
		red8 := uint8(red)
		color := color.RGBA{red8, intens8 / 2, intens8, 0xff}
		for y := 0; y < height; y++ {
			result.Set(x, y, color)
		}
	}
	return result
}

func TestCompositorSingleOverlap(t *testing.T) {
	tileWidth := 120
	tileHeight := 64

	tileOverlap := 16

	tilesAcross := 2
	tilesDown := 1

	width := tilesAcross*(tileWidth-tileOverlap) + tileOverlap
	height := tilesDown*(tileHeight-tileOverlap) + tileOverlap
	rect := image.Rect(0, 0, width, height)
	compositor := NewCompositor(rect)

	outDir := "test_data/out/compositor_test/"
	ensureDirExists(outDir)

	for col := 0; col < tilesAcross; col++ {
		x := (tileWidth - tileOverlap) * col
		for row := 0; row < tilesDown; row++ {
			y := (tileHeight - tileOverlap) * row
			tileImage := makeTile(row, col, tileWidth, tileHeight)
			savePNG(tileImage, fmt.Sprintf("%stile_%d_%d.png", outDir, col, row), t)
			// Save its HSV equivalent:
			hsvTileImage := hsv_image.HSVFromImage(tileImage)
			savePNG(hsvTileImage, fmt.Sprintf("%stile_%d_%d_hsv.png", outDir, col, row), t)
			subframeRect := image.Rect(x, y, x+tileWidth, y+tileHeight)
			fmt.Printf("AddImage (%d, %d) at %v\n", row, col, subframeRect)
			compositor.AddImage(tileImage, subframeRect)
		}
	}

	want := tilesAcross * tilesDown
	got := len(compositor.addedAreas)
	if want != got {
		t.Fatal("Expected number of addedAreas:", want, "; actual:", got)
	}

	savePNG(compositor.Result, fmt.Sprintf("%s/single_overlap.png", outDir), t)
	fmt.Println("This is a manual 'test'.")
	fmt.Println("Review files in", outDir, "for correctness.")
}
