package lib

type CameraInfo struct {
	CameraModelComponentList interface{} `json:"camera_model_component_list"`
	CameraModelType          string      `json:"camera_model_type"`
	CameraPosition           []float64   `json:"camera_position"`
	CameraVector             []float64   `json:"camera_vector"`
	FilterName               string      `json:"filter_name"`
	Instrument               string      `json:"instrument"`
}

type ImageInfo struct {
	ImageId string `json:"image_id"`
	Credit  string `json:"credit"`
	Caption string `json:"caption"`
	Title   string `json:"title"`

	Attitude []float64 `json:"attitude"`

	Camera CameraInfo `json:"camera"`
}
