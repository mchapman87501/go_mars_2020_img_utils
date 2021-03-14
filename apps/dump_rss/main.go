package main

import (
	"fmt"

	"com.dmoonc/mchapman/go_mars_perseverance_images/lib"
)

func main() {
	cameras := []string{"NAVCAM_LEFT", "NAVCAM_RIGHT"}
	params := lib.GetRequestParams(cameras, 5, 1, -1, -1)
	responseText := lib.GetImageMetadata(params)
	fmt.Println("Response:", responseText)
}
