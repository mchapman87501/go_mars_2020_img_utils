package lib

import (
	"database/sql"
	"fmt"
	"image"
	"math"
)

type CompositeImageInfo struct {
	ImageID string
	Site    int
	Drive   int
	Sclk    float64

	SubframeRect image.Rectangle

	Camera string
}

func valOrNan(fval sql.NullFloat64) float64 {
	if fval.Valid {
		return fval.Float64
	}
	return math.NaN()
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
		x := 0
		y := 0
		width := 0
		height := 0

		// Hm... SQLite3 typically stores NaN as NULL?
		// So hints StackOverflow
		sclkValue := sql.NullFloat64{}

		err := rows.Scan(
			&record.ImageID, &record.Site, &record.Drive, &sclkValue,
			&x, &y, &width, &height, &record.Camera)
		if err != nil {
			err = fmt.Errorf("error extracting row %v: %v", irow, err)
			return result, err
		}
		record.Sclk = valOrNan(sclkValue)
		record.SubframeRect = image.Rect(x, y, x+width, y+height)
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

type CompositeImageSet []CompositeImageInfo

// Implement sort.Interface to order composite images by their subframe rectangles
func (cid CompositeImageSet) Len() int {
	return len(cid)
}

func (cid CompositeImageSet) Swap(i, j int) {
	cid[i], cid[j] = cid[j], cid[i]
}

func (cid CompositeImageSet) Less(i, j int) bool {
	imin := cid[i].SubframeRect.Min
	jmin := cid[j].SubframeRect.Min
	if imin.X < jmin.X {
		return true
	}
	if imin.X == jmin.X {
		return imin.Y < jmin.Y
	}
	return false
}

func (imageSet CompositeImageSet) Name() string {
	if len(imageSet) <= 0 {
		return "empty_image_set"
	}
	firstImage := imageSet[0]
	return fmt.Sprintf("image_set_%v_%v_%v_%v", firstImage.Camera, firstImage.Sclk, firstImage.Site, firstImage.Drive)
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
				result = append(result, currImages)
			}
			currImages = []CompositeImageInfo{}
			prevSclk = record.Sclk
		}
		currImages = append(currImages, record)
	}

	if len(currImages) > 0 {
		result = append(result, currImages)
	}
	return result, nil
}
