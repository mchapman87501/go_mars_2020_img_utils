package lib

import (
	"database/sql"
	"fmt"
)

type CompositeImageInfo struct {
	ImageID string
	Site    int
	Drive   int
	Sclk    float64

	SubframeRect Rect

	Camera string
}

func GetCompositeImageInfoRecords(idb ImageDB, camera string) ([]CompositeImageInfo, error) {
	result := []CompositeImageInfo{}
	rows, err := retrieveImageSets(idb, camera)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	irow := 0
	for rows.Next() {
		record := CompositeImageInfo{}
		origin := Origin{}
		size := Size{}
		err := rows.Scan(
			&record.ImageID, &record.Site, &record.Drive, &record.Sclk,
			&origin.X, &origin.Y, &size.Width, &size.Height, &record.Camera)
		if err != nil {
			err = fmt.Errorf("error extracting row %v: %v", irow, err)
			return result, err
		}
		record.SubframeRect.Origin = origin
		record.SubframeRect.Size = size
		result = append(result, record)
	}

	return result, nil
}

func retrieveImageSets(idb ImageDB, camera string) (*sql.Rows, error) {
	query := `SELECT
			image_id, site, drive, ext_sclk,
			ext_sf_left x, ext_sf_top y,
			ext_sf_width w, ext_sf_height h,
			cam_instrument
		FROM Images
		WHERE cam_instrument = ?
			AND sample_type = 'Full'
			AND ext_scale_factor = 1
			AND x NOT NULL
			AND y NOT NULL
			AND w NOT NULL
			and h NOT NULL
		ORDER BY site, drive, ext_sclk, image_id`

	return idb.DB.Query(query, camera)
}

type CompositeImageSet struct {
	Images []CompositeImageInfo
}

func GetCompositeImageSets(idb ImageDB, camera string) ([]CompositeImageSet, error) {
	result := []CompositeImageSet{}

	records, err := GetCompositeImageInfoRecords(idb, camera)
	if err != nil {
		return result, err
	}

	prevSclk := -1.0
	currImages := []CompositeImageInfo{}
	for _, record := range records {
		if record.Sclk != prevSclk {
			if len(currImages) > 0 {
				result = append(result, CompositeImageSet{currImages})
			}
			currImages = []CompositeImageInfo{}
			prevSclk = record.Sclk
		}
		currImages = append(currImages, record)
	}

	if len(currImages) > 0 {
		result = append(result, CompositeImageSet{currImages})
	}
	return result, nil
}
