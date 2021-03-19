package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"sort"

	"com.dmoonc/mchapman87501/mars_2020_img_utils/lib"
)

func savePNG(image image.Image, filename string) {
	outf, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating %v: %v\n", filename, err)
	}
	defer outf.Close()

	err = png.Encode(outf, image)
	if err != nil {
		fmt.Printf("Error saving %v: %v\n", filename, err)
	}
}

func assembleImageSet(cache lib.ImageCache, imageSet lib.CompositeImageSet) {
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
		image, err := cache.FullSize(record.ImageID)
		if err != nil {
			fmt.Println("Error retrieving full size image", record.ImageID, "- skipping")
		} else {
			compositor.AddImage(image, record.SubframeRect)
		}
	}

	savePNG(compositor.Result, imageSet.Name()+".png")
}

func assembleImageSets(imageDB lib.ImageDB, imageSets []lib.CompositeImageSet) {
	cache, err := lib.NewImageCache(imageDB)
	if err != nil {
		log.Fatal("Could not create image cache:", err)
	}
	for _, imageSet := range imageSets {
		assembleImageSet(cache, imageSet)
	}
}

func main() {
	imageDB, err := lib.NewImageDB()
	if err != nil {
		log.Fatal("Could not instantiate image DB:", err)
	}

	cameras := imageDB.Cameras()
	fmt.Println("Cameras:", cameras)
	for _, camera := range cameras {
		fmt.Println("Finding composite image sets from", camera)
		imageSets, err := lib.GetCompositeImageSets(imageDB, camera)
		if err != nil {
			fmt.Println("Error retrieving image sets for", camera, "-", err)
		} else {
			assembleImageSets(imageDB, imageSets)
		}
	}
}