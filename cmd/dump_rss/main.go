package main

import (
	"fmt"
	"log"

    "com.dmoonc/mchapman87501/mars_2020_img_utils/lib"
)

func main() {
	cameras := []string{} //lib.ValidCameras()
	page := 0
	for {
		fmt.Println("Page", page+1)

		params := lib.GetRequestParams(cameras, 100, page, -1, -1)
		images, err := lib.GetImageMetadata(params)

		if err != nil {
			log.Fatal(err)
		}

		if len(images) <= 0 {
			return
		}

		for _, record := range images {
			fmt.Println(record.ImageID, record.Extended.SubframeRect)
		}
		page += 1
	}
}
