package lib

import (
	"database/sql"
	"fmt"
	"image"
	"math"
)

// Note about ColorType:
// The 3rd letter of the image ID can tell what type of color data the image contains:
// R -- image is (usually) PNG RGB, but it represents the red band of
//      a full image
// G -- as above, but green
// B -- as above, but blue
// E -- PNG RGB containing only grayscale data; this is the full
//      readout of all of the pixels underneath the Bayer RGGB color filter
//      array, and it needs to be de-mosaiced to create a full color image
// F -- PNG RGB, all color channels, already de-mosaiced

type CompositeImageInfo struct {
	ImageID string
	Site    int
	Drive   int
	Sclk    float64

	SubframeRect image.Rectangle

	Camera    string
	ColorType string
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
		colorTypeStr := "U"

		err := rows.Scan(
			&record.ImageID, &record.Site, &record.Drive, &sclkValue,
			&colorTypeStr, &x, &y, &width, &height, &record.Camera)
		if err != nil {
			err = fmt.Errorf("error extracting row %v: %v", irow, err)
			return result, err
		}
		record.Sclk = valOrNan(sclkValue)
		record.SubframeRect = image.Rect(x, y, x+width, y+height)
		record.ColorType = colorTypeStr
		result = append(result, record)
	}

	return result, nil
}

func retrieveImageSets(idb ImageDB, camera string) (*sql.Rows, error) {
	query := `SELECT
			image_id, site, drive, ext_sclk, color_type,
			ext_sf_left x, ext_sf_top y,
			ext_sf_width w, ext_sf_height h,
			cam_instrument
		FROM Images
		WHERE cam_instrument = ?
			AND sample_type = 'Full'
			AND color_type = 'E'
			AND ext_scale_factor = 1
			AND x NOT NULL
			AND y NOT NULL
			AND w NOT NULL
			and h NOT NULL
		ORDER BY site, drive, ext_sclk, color_type, image_id`

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
	return fmt.Sprintf("%v_%v_%.0f_%v_%v", firstImage.Camera, firstImage.ColorType, firstImage.Sclk, firstImage.Site, firstImage.Drive)
}

func GetCompositeImageSets(idb ImageDB, camera string) ([]CompositeImageSet, error) {
	result := []CompositeImageSet{}

	records, err := GetCompositeImageInfoRecords(idb, camera)
	if err != nil {
		return result, err
	}

	prevSclk := -1.0
	prevColorType := "U"
	currImages := []CompositeImageInfo{}
	for _, record := range records {
		// Images constituting a single composite have the same sclk and the
		// same color type.
		if (record.Sclk != prevSclk) || (record.ColorType != prevColorType) {
			// Composite image sets must have more than one constituent image.
			if len(currImages) > 1 {
				result = append(result, currImages)
			}
			currImages = []CompositeImageInfo{}
			prevSclk = record.Sclk
			prevColorType = record.ColorType
		}
		currImages = append(currImages, record)
	}

	if len(currImages) > 1 {
		result = append(result, currImages)
	}
	return result, nil
}
