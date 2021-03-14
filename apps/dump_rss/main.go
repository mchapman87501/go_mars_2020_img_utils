package main

import (
	"fmt"

	"com.dmoonc/mchapman/go_mars_perseverance_images/lib"
)

func main() {
	cameras := []string{"NAVCAM_LEFT", "NAVCAM_RIGHT"}
	params := lib.GetRequestParams(cameras, 5, 1, -1, -1)
	images, err := lib.GetImageMetadata(params)

	if err != nil {
		fmt.Println("Error retrieving images:", err)
	}

	fmt.Println("Number of images:", len(images))
	if len(images) > 0 {
		for i, v := range images {
			fmt.Printf("Image %d: %v\n\n", i, v)
		}
	}
}
