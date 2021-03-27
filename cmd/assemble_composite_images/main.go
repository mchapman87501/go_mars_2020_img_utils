package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"

	"com.dmoonc/mchapman87501/mars_2020_img_utils/lib"
)

const outDir = "composite_images/"

func savePNG(image image.Image, filename string) {
	outf, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating %v: %v\n", filename, err)
		return
	}
	defer outf.Close()

	err = png.Encode(outf, image)
	if err != nil {
		fmt.Printf("Error saving %v: %v\n", filename, err)
	}
}

func demosaiced(cache lib.ImageCache, record lib.CompositeImageInfo) (image.Image, error) {
	image, err := cache.FullSize(record.ImageID)
	if err != nil {
		return image, err
	}
	if record.ColorType == "E" {
		return lib.DemosaicRGBGray(image)
	}
	return image, nil
}

func saveMetadata(records lib.CompositeImageSet, filename string) {
	b, err := json.MarshalIndent(records, "", "  ")
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

func assembleImageSet(cache lib.ImageCache, imageSet lib.CompositeImageSet) {
	fmt.Println("Processing", imageSet.Name())
	filename := outDir + imageSet.Name() + ".png"
	metadataFilename := outDir + imageSet.Name() + "_metadata.json"

	if len(imageSet) < 2 {
		fmt.Println("Image set does not contain multiple images.")
		return
	}
	sorted := make(lib.CompositeImageSet, len(imageSet))
	copy(sorted, imageSet)
	sort.Sort(sorted)

	compositeRect := image.Rectangle{}
	for _, record := range sorted {
		compositeRect = compositeRect.Union(record.SubframeRect)
	}

	compositor := lib.NewCompositor(compositeRect)
	if compositor.Bounds.Empty() {
		fmt.Println("composite image set has no extent.")
		return
	}

	for _, record := range sorted {
		image, err := demosaiced(cache, record)
		if err != nil {
			fmt.Println("Error retrieving full size image", record.ImageID, "- skipping")
		} else {
			compositor.AddImage(image, record.SubframeRect)
		}
	}

	compositor.AutoEnhance()
	savePNG(compositor.Result, filename)
	saveMetadata(sorted, metadataFilename)
}

func enqueueCameraImageSets(
	imageDB lib.ImageDB, camera string, jobs chan lib.CompositeImageSet,
) {
	imageSets, err := lib.GetCompositeImageSets(imageDB, camera)
	if err != nil {
		fmt.Println("Error retrieving image sets for", camera, "-", err)
	} else {
		for _, imageSet := range imageSets {
			jobs <- imageSet
		}
	}

}

func processJobs(jobs chan lib.CompositeImageSet, cache lib.ImageCache, wg *sync.WaitGroup) {
	for {
		imageSet, ok := <-jobs
		if !ok {
			wg.Done()
			return
		}
		assembleImageSet(cache, imageSet)
	}
}

func main() {
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		log.Fatal("Could not create output directory", outDir, ":", err)
	}

	imageDB, err := lib.NewImageDB()
	if err != nil {
		log.Fatal("Could not instantiate image DB:", err)
	}

	// Lots of thread-safety issues here -- need to mutex access to
	// cache operations.
	cache, err := lib.NewImageCache(imageDB)
	if err != nil {
		log.Fatal("Could not instantiate image cache:", err)
	}

	cameras := imageDB.Cameras()
	concurrency := runtime.NumCPU() * 3

	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	jobs := make(chan lib.CompositeImageSet, concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			processJobs(jobs, cache, &wg)
		}()
	}
	for _, camera := range cameras {
		enqueueCameraImageSets(imageDB, camera, jobs)
	}
	close(jobs)
	wg.Wait()
}
