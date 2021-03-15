package lib

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

type optFloat float64

func isUnknown(data []byte) bool {
	return string(data) == `"UNK"`
}

func (v *optFloat) UnmarshalJSON(data []byte) error {
	sVal := string(data)
	if isUnknown(data) {
		*v = optFloat(math.NaN())
		return nil
	}
	fmtStr := "%g"
	if sVal[0] == '"' {
		fmtStr = `"%g"`
	}
	_, err := fmt.Sscanf(sVal, fmtStr, v)
	return err
}

type optInt int

func (v *optInt) UnmarshalJSON(data []byte) error {
	sVal := string(data)
	if isUnknown(data) {
		*v = 0
		return nil
	}
	fmtStr := "%d"
	if sVal[0] == '"' {
		fmtStr = `"%d"`
	}
	_, err := fmt.Sscanf(sVal, fmtStr, v)

	return err
}

type FloatTuple []float64

func (v *FloatTuple) UnmarshalJSON(data []byte) error {
	if isUnknown(data) {
		return nil // Empty
	}
	// avoid recursion.
	var vBase []float64
	err := json.Unmarshal(tupleAsArray(data), &vBase)
	if err != nil {
		return err
	}
	*v = vBase
	return nil
}

type CameraInfo struct {
	CameraModelComponentList interface{} `json:"camera_model_component_list"`
	CameraModelType          string      `json:"camera_model_type"`
	CameraPosition           FloatTuple  `json:"camera_position"`
	CameraVector             FloatTuple  `json:"camera_vector"`
	FilterName               string      `json:"filter_name"`
	Instrument               string      `json:"instrument"`
}

type ImageFileUrls struct {
	Medium  string `json:"medium"`
	Small   string `json:"small"`
	FullRes string `json:"full_res"`
	Large   string `json:"large"`
}

type Origin struct {
	X int
	Y int
}
type Size struct {
	Width  int
	Height int
}

func tupleAsArray(data []byte) []byte {
	return []byte(strings.Replace(
		strings.Replace(
			strings.Replace(string(data), ")", "]", 1),
			"(", "[", 1),
		`"`, "", -1))
}

func (size *Size) UnmarshalJSON(data []byte) error {
	if isUnknown(data) {
		return nil
	}
	var coords []int
	err := json.Unmarshal(tupleAsArray(data), &coords)
	if err != nil {
		return err
	}
	size.Width = coords[0]
	size.Height = coords[1]
	return nil
}

type Rect struct {
	Origin Origin
	Size   Size
}

func (rect *Rect) UnmarshalJSON(data []byte) error {
	if isUnknown(data) {
		return nil
	}
	arrayBytes := tupleAsArray(data)

	var coords []int
	err := json.Unmarshal(arrayBytes, &coords)
	if err != nil {
		return err
	}
	rect.Origin.X = coords[0]
	rect.Origin.Y = coords[1]
	rect.Size.Width = coords[2]
	rect.Size.Height = coords[3]
	return nil
}

type ExtendedInfo struct {
	// These are all either a float or "UNK"
	MastAzimuth   optFloat `json:"mastAz"`
	MastElevation optFloat `json:"mastEl"`
	Sclk          optFloat `json:"sclk"`
	// "UNK" or a positive integer, typically 1 or 4.
	ScaleFactor optInt `json:"optInt"`

	XYZ FloatTuple `json:"xyz"`
	// How to convince JSON to parse an array of 4 floats as a struct?
	SubframeRect Rect `json:"subframeRect"`
	Dimension    Size `json:"dimension"`
}

type ImageInfo struct {
	ImageId string `json:"imageid"`
	Credit  string `json:"credit"`
	Caption string `json:"caption"`
	Title   string `json:"title"`

	Attitude FloatTuple `json:"attitude"`
	Sol      optInt     `json:"sol"`

	// How to parse date strings in Go?
	DateTakenMars string `json:"date_taken_mars"`
	DateTakenUtc  string `json:"date_taken_utc"`
	DateReceived  string `json:"date_received"`

	Drive optInt `json:"drive"`
	Site  optInt `json:"site"`

	SampleType string `json:"sample_type"`

	Camera     CameraInfo    `json:"camera"`
	ImageFiles ImageFileUrls `json:"image_files"`

	Extended ExtendedInfo `json:"extended"`
}
