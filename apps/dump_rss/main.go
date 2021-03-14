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

	var m interface{}
	err := json.Unmarshal(responseBytes, &m)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
	} else {
		fmt.Println("Response:", m)
	}
}
