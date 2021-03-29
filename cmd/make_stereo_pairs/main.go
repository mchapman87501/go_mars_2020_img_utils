package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"com.dmoonc/mchapman87501/mars_2020_img_utils/lib"
)

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
	suffixes := []string{"_LEFT", "_RIGHT"}

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
	return (pf1 == pf2) && (lr1 != "") && (lr2 != "")
}

func findStereoPairs(imageDB lib.ImageDB) []StereoPair {
	result := []StereoPair{}

	query := `SELECT image_id, cam_instrument, ext_sclk
FROM Images
WHERE color_type = 'F' and sample_type = 'Full'
ORDER BY ext_sclk, cam_instrument`

	rows, err := imageDB.DB.Query(query)
	if err == nil {
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
				if (dt <= 2.0) && idsMatch(prevID, imageID) && instrumentsMatch(instrument, prevInstrument) {
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

func main() {
	imageDB, err := lib.NewImageDB()
	if err != nil {
		log.Fatalf("Can't create database.")
	}
	for _, pair := range findStereoPairs(imageDB) {
		fmt.Println("L:", pair.Left, "R:", pair.Right)
	}
}

