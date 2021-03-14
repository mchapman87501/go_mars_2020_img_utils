package lib

import (
	"encoding/json"
	"strings"
)

type CameraInfo struct {
	CameraModelComponentList interface{} `json:"camera_model_component_list"`
	CameraModelType          string      `json:"camera_model_type"`
	CameraPosition           []float64   `json:"camera_position"`
	CameraVector             []float64   `json:"camera_vector"`
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
	MastAzimuth   float64 `json:"mastAz"`
	MastElevation float64 `json:"mastEl"`

	Sclk float64 `json:"sclk"`

	ScaleFactor int       `json:"scaleFactor"`
	XYZ         []float64 `json:"xyz"`
	// How to convince JSON to parse an array of 4 floats as a struct?
	SubframeRect Rect `json:"subframeRect"`
	Dimension    Size `json:"dimension"`
}

type ImageInfo struct {
	ImageId string `json:"imageid"`
	Credit  string `json:"credit"`
	Caption string `json:"caption"`
	Title   string `json:"title"`

	Attitude []float64 `json:"attitude"`
	Sol      int       `json:"sol"`

	// How to parse date strings in Go?
	DateTakenMars string `json:"date_taken_mars"`
	DateTakenUtc  string `json:"date_taken_utc"`
	DateReceived  string `json:"date_received"`

	Drive int `json:"drive"`
	Site  int `json:"site"`

	SampleType string `json:"sample_type"`

	Camera     CameraInfo    `json:"camera"`
	ImageFiles ImageFileUrls `json:"image_files"`

	Extended ExtendedInfo `json:"extended"`
}
