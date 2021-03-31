package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"

	"dmoonc.com/mchapman87501/mars_2020_img_utils/lib"
)

const outDir = "stereo_images/"

func init() {
	os.MkdirAll(outDir, 0755)
}

type StereoPair struct {
	Left, Right string // image_ids
}

func idsMatch(imageID1, imageID2 string) bool {
	// If TWO IDs differ only in the 2nd rune - one being
	// 'L' and the other being 'R' - then the IDs are
	// probably for the same snapshot.
	run1 := []rune(imageID1)
	run2 := []rune(imageID2)
	if (len(run1) > 3) && (len(run2) > 3) {
		prefix1 := string(run1[:1])
		prefix2 := string(run2[:1])
		suffix1 := string(run1[2:])
		suffix2 := string(run2[2:])
		if (prefix1 == prefix2) && (suffix1 == suffix2) {
			return true
		}
	}
	return false
}

func splitLRSuffix(s string) (prefix, suffix string) {
	// Front hazcams use, e.g., "LEFT_A"
	suffixes := []string{"_LEFT", "_RIGHT", "LEFT_A", "RIGHT_A"}

	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return strings.Replace(s, suffix, "", 1), suffix
		}
	}
	return s, ""
}

func instrumentsMatch(inst1, inst2 string) bool {
	pf1, lr1 := splitLRSuffix(inst1)
	pf2, lr2 := splitLRSuffix(inst2)
	return (pf1 == pf2) && (lr1 != "") && (lr2 != "") && (lr1 != lr2)
}

func findStereoPairs(imageDB lib.ImageDB) []StereoPair {
	result := []StereoPair{}

	query := `SELECT image_id, cam_instrument, ext_sclk
FROM Images
WHERE color_type = 'F' and sample_type = 'Full'
ORDER BY ext_sclk, cam_instrument`

	rows, err := imageDB.DB.Query(query)
	if err == nil {
		defer rows.Close()
		prevID := ""
		prevInstrument := ""
		prevSclk := -1.0

		for rows.Next() {
			var imageID string
			var instrument string
			var sclk float64
			err := rows.Scan(&imageID, &instrument, &sclk)
			if err == nil {
				dt := math.Abs(sclk - prevSclk)
				if (dt <= 1.0) && idsMatch(prevID, imageID) && instrumentsMatch(instrument, prevInstrument) {
					pair := StereoPair{imageID, prevID}
					leftOrRightPrev := []rune(prevID)[1]
					if leftOrRightPrev == 'L' {
						pair = StereoPair{prevID, imageID}
					}
					result = append(result, pair)
				}
			}

			prevID = imageID
			prevInstrument = instrument
			prevSclk = sclk
		}
	}
	return result
}

func savePNG(image image.Image, filename string) {
	if err := lib.SavePNG(image, filename); err != nil {
		fmt.Printf("Error saving %v: %v\n", filename, err)
	}
}

func saveMetadata(sp StereoPair, filename string) {
	b, err := json.MarshalIndent(sp, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling composite image set to JSON:", err)
	} else {
		outf, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Error creating %v: %v\n", filename, err)
			return
		}
		defer outf.Close()
		bytesWritten, err := outf.Write(b)
		if err != nil {
			fmt.Printf("Error writing %v: %v\n", filename, err)
		}
		if bytesWritten < len(b) {
			fmt.Printf("Did not write all JSON data to %v.\n", filename)
		}
	}
}

func makeImage(imageDB lib.ImageDB, sp StereoPair) (image.Image, error) {
	cache, err := lib.NewImageCache(imageDB)
	if err != nil {
		return nil, fmt.Errorf("can't create image cache: %v", err)
	}

	leftImage, err := cache.FullSize(sp.Left)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve image %v: %v", sp.Left, err)
	}

	rightImage, err := cache.FullSize(sp.Right)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve image %v: %v", sp.Right, err)
	}

	rightReExposed := lib.MatchExposure(rightImage, leftImage)

	leftBounds := leftImage.Bounds()
	rightBounds := rightReExposed.Bounds()
	if leftBounds.Dx() != rightBounds.Dx() {
		return nil, fmt.Errorf("images have different widths: %v=%v, %v=%v", sp.Left, leftBounds.Dx(), sp.Right, rightBounds.Dx())
	}
	if leftBounds.Dy() != rightBounds.Dy() {
		return nil, fmt.Errorf("images have different heights: %v=%v, %v=%v", sp.Left, leftBounds.Dy(), sp.Right, rightBounds.Dy())
	}
	width := leftBounds.Dx() + rightBounds.Dx()
	height := leftBounds.Dy()
	bounds := image.Rect(0, 0, width, height)
	result := image.NewRGBA(bounds)

	leftRect := image.Rect(0, 0, leftBounds.Dx(), leftBounds.Dy())
	draw.Src.Draw(result, leftRect, leftImage, leftBounds.Min)
	rightRect := image.Rect(leftBounds.Dx(), 0, leftBounds.Dx()+rightBounds.Dx(), rightBounds.Dy())
	draw.Src.Draw(result, rightRect, rightReExposed, rightBounds.Min)

	// TODO adjust dynamic range.
	return result, nil
}

type Job struct {
	Index int
	Pair  StereoPair
}

func processJobs(
	workerID int, jobs chan Job, imageDB lib.ImageDB,
	wg *sync.WaitGroup,
) {
	for {
		job, ok := <-jobs
		if !ok {
			wg.Done()
			return
		}

		i := job.Index
		pair := job.Pair
		name := fmt.Sprintf(
			"stereo_%04d_%v", i, strings.Replace(pair.Left, "L", "", 1))
		pngName := outDir + name + ".png"
		jsonName := outDir + name + ".json"

		if !lib.FileExists(pngName) {
			fmt.Println("L:", pair.Left, "R:", pair.Right)
			image, err := makeImage(imageDB, pair)
			if err != nil {
				fmt.Println("Error creating stereo pair:", err)
			} else {
				savePNG(image, pngName)
				saveMetadata(pair, jsonName)
			}
		}
	}
}

func processConcurrently(imageDB lib.ImageDB) {
	concurrency := runtime.NumCPU()

	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	jobs := make(chan Job, concurrency)
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			processJobs(workerID, jobs, imageDB, &wg)
		}(i)
	}

	for i, pair := range findStereoPairs(imageDB) {
		jobs <- Job{i, pair}
	}

	close(jobs)
	wg.Wait()
}

func main() {
	imageDB, err := lib.NewImageDB()
	if err != nil {
		log.Fatal("Could not instantiate image DB:", err)
	}

	processConcurrently(imageDB)
}
