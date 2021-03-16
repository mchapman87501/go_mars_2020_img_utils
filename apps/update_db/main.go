package main

import (
	"fmt"
	"log"

	"com.dmoonc/mchapman/go_mars_perseverance_images/lib"
)

func main() {
	imageDB, err := lib.NewImageDB()
	if err != nil {
		log.Fatal("Could not instantiate image DB:", err)
	}

	cameras := []string{} //lib.ValidCameras()
	page := 0
	for {
		fmt.Println("Page", page+1)

		params := lib.GetRequestParams(cameras, 100, page, -1, -1)
		records, err := lib.GetImageMetadata(params)

		if err != nil {
			log.Fatal(err)
		}

		if len(records) <= 0 {
			return
		}

		if err = imageDB.AddOrUpdate(records); err != nil {
			log.Fatal(err)
		}
		page += 1
	}
}
