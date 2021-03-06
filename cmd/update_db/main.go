package main

import (
	"fmt"
	"log"

	"github.com/mchapman87501/go_mars_2020_img_utils/lib"
)

func main() {
	imageDB, err := lib.NewImageDB()
	if err != nil {
		log.Fatal("Could not instantiate image DB:", err)
	}

	cameras := []string{}
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
