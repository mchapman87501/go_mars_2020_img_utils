package main

import (
	"encoding/json"
	"fmt"

	"com.dmoonc/mchapman/go_mars_perseverance_images/lib"
)

func main() {
	cameras := []string{"NAVCAM_LEFT", "NAVCAM_RIGHT"}
	params := lib.GetRequestParams(cameras, 5, 1, -1, -1)
	responseBytes := lib.GetImageMetadata(params)

	var fields interface{}
	err := json.Unmarshal(responseBytes, &fields)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
	} else {
		m := fields.(map[string]interface{})
		images := m["images"].([]interface{})
		if len(images) <= 0 {
			fmt.Println("Reply had no images")
		} else {
			firstImage := images[0]
			fmt.Println("First image:", firstImage)
		}
	}
}
