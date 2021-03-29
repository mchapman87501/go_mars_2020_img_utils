package lib

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"sort"
	"testing"

	lib_image "dmoonc.com/mchapman87501/mars_2020_img_utils/lib/image"
)

const troublesomeTilesOutDir = "test_data/out/compositor_test/"
const pythonScript = troublesomeTilesOutDir + "plot_adjustments.py"

var pyfile *os.File

func init() {
	os.MkdirAll(troublesomeTilesOutDir, 0755)

	var err error
	pyfile, err = os.Create(pythonScript)
	if err != nil {
		log.Fatalf("Could not open %v for writing: %v\n", pythonScript, err)
	}
}

func pythonln(a ...interface{}) (n int, err error) {
	return pyfile.WriteString(fmt.Sprintln(a...))
}

func pythonf(format string, a ...interface{}) (n int, err error) {
	return pyfile.WriteString(fmt.Sprintf(format, a...))
}

func TestTroublesomeTiles(t *testing.T) {
	outDir := troublesomeTilesOutDir

	tileDir := "test_data/in/compositor_test/tiles/"
	readTile := func(tileName string) image.Image {
		pathname := tileDir + tileName + ".png"
		inf, err := os.Open(pathname)
		if err != nil {
			t.Fatalf("Could not open tile image %s: %v\n", pathname, err)
		}
		result, err := png.Decode(inf)
		if err != nil {
			t.Fatalf("Could not decode PNG %s: %v\n", pathname, err)
		}

		return result
	}

	imageNames := []string{
		"demosaic_NLE_0009_0667755063_497ECM_N0030000NCAM00102_01_0LLJ",
		"demosaic_NLE_0009_0667755063_497ECM_N0030000NCAM00102_05_0LLJ",
		"demosaic_NLE_0009_0667755063_497ECM_N0030000NCAM00102_09_0LLJ",
	}
	// width, height, overlap: 644, 484, 8
	imageRects := []image.Rectangle{
		image.Rect(1, 1, 645, 485),
		image.Rect(1, 477, 645, 961),
		image.Rect(1, 953, 645, 1437),
	}

	fullRect := image.Rect(0, 0, 0, 0)
	for _, r := range imageRects {
		fullRect = fullRect.Union(r)
	}
	compositor := NewCompositor(fullRect)

	for i, name := range imageNames {
		tile := readTile(name)
		labTile := lib_image.CIELabFromImage(tile)
		am := compositor.makeValueAdjustmentMap(labTile, imageRects[i])
		printVAM(am, i)
		printTargLVals(&compositor, i, imageRects[i], t)
		compositor.AddImage(tile, imageRects[i])
	}
	compositor.CompressDynamicRange()
	savePNG(compositor.Result, outDir+"troublesome_tiling.png", t)

	// Don't hang the test when it cannot be run interactively.
	fmt.Println("This is a manual 'test'.")
	fmt.Println("Review the files in", troublesomeTilesOutDir)
	fmt.Println("Also run", pythonScript)
}

func printTargLVals(comp *Compositor, index int, r image.Rectangle, t *testing.T) {
	if len(comp.addedAreas) <= 0 {
		return
	}
	pythonf(`import matplotlib.pyplot as plt
import numpy as np
l_targ_%d = np.array([
`, index)
	origin := comp.Result.Bounds().Min
	for _, rect := range comp.addedAreas {
		pythonln("    # Composite image region:", rect)
		overlap := rect.Intersect(r)
		pythonln("    # Overlap:", overlap)

		if overlap.Empty() {
			pythonln("    # EMPTY")
		} else {
			for x := overlap.Min.X; x < overlap.Max.X; x++ {
				for y := overlap.Min.Y; y < overlap.Max.Y; y++ {
					p := comp.Result.CIELabAt(x-origin.X, y-origin.Y)
					pythonf("%.2f, ", p.L)
				}
				pythonln("")
			}
			// srcRegionImage := comp.Result.SubImage(overlap)
			// // Save only the L channel of this image.  I hope I'm making a copy...

			// saveImageLChannel(*srcRegionImage, index, regionIndex, t)
			pythonln("")
		}
	}
	pythonln("])")
	pythonf("plt.hist(l_targ_%d, bins=100)\n", index)
	pythonf(`plt.title("Composited L values that overlap image %d")
`, index)
	pythonln("plt.show()")
}

// func saveImageLChannel(image lib_image.CIELab, index int, regionIndex int, t *testing.T) {
// 	numPixels := 0
// 	numBrightest := 0
// 	for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
// 		for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
// 			numPixels += 1
// 			pix := image.CIELabAt(x, y)
// 			pix.B = 0
// 			if pix.L <= 99.0 {
// 				pix.A = 0
// 			} else {
// 				pix.A = 100.0
// 				numBrightest += 1
// 			}
// 			image.SetCIELab(x, y, pix)
// 		}
// 	}
// 	outDir := troublesomeTilesOutDir
// 	regionFilename := fmt.Sprintf("%s/sub_comp_%d_%d.png", outDir, index, regionIndex)
// 	pythonln("")
// 	pythonln("# Tile", index, "region", regionIndex, "bright fraction:", float64(numBrightest)/float64(numPixels))
// 	savePNG(&image, regionFilename, t)
// }

func printVAM(am *AdjustmentMap, index int) {
	if len(am.L) <= 0 {
		return
	}
	pythonln("")
	pythonln("import matplotlib.pyplot as plt")
	pythonln("import numpy as np")
	pythonf("adj_%d = np.array([", index)
	sortedSrc := make([]float64, 0, len(am.L))
	for src := range am.L {
		sortedSrc = append(sortedSrc, src)
	}
	sort.Float64s(sortedSrc)

	for _, src := range sortedSrc {
		targ := am.L[src]
		pythonf("    [%v, %.4f],\n", src, targ)
	}
	pythonln("])")

	pythonf(`fig = plt.figure(figsize=(6.4, 6.4 * 1.5))
adj_plot, src_hist, targ_hist = fig.subplots(3, 1)
src_vals = adj_%d[:, 0]
targ_vals = adj_%d[:, 1]
adj_plot.plot(src_vals, targ_vals)
adj_plot.set_title("adjustment_%d")
src_hist.hist(src_vals, bins=100)
src_hist.set_title("L values of Tile %d")
targ_hist.hist(targ_vals, bins=100)
targ_hist.set_title("L Values of composited image")
plt.show()


`, index, index, index, index)
}
