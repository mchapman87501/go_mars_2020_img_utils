package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
)

func ValidCameras() []string {
	return []string{
		"FRONT_HAZCAM_LEFT_A",
		"FRONT_HAZCAM_LEFT_B",
		"FRONT_HAZCAM_RIGHT_A",
		"FRONT_HAZCAM_RIGHT_B",

		"REAR_HAZCAM_LEFT",
		"REAR_HAZCAM_RIGHT",

		"NAVCAM_LEFT",
		"NAVCAM_RIGHT",

		"MCZ_LEFT",
		"MCZ_RIGHT",

		"EDL_DDCAM",
		"EDL_PUCAM1",
		"EDL_PUCAM2",
		"EDL_RUCAM",
	}
}

func getCamerasParam(cameras []string) string {
	result := ""
	sep := ""
	for _, v := range cameras {
		result = result + sep + v
		sep = "|"
	}
	return result
}

func GetRequestParams(cameras []string, numPerPage int, page int, minSol int, maxSol int) url.Values {
	result := url.Values{}
	result.Set("feed", "raw_images")
	result.Set("category", "mars2020")
	result.Set("feedtype", "json")
	result.Set("num", fmt.Sprint(numPerPage))
	result.Set("page", fmt.Sprint(page))
	result.Set("order", "sol desc")

	if len(cameras) > 0 {
		result.Set("search", getCamerasParam(cameras))
	}
	if minSol >= 0 {
		result.Set("condition_2", fmt.Sprintf("%d:sol:gte", minSol))
	}
	if (maxSol >= minSol) && (maxSol > 0) {
		result.Set("condition_3", fmt.Sprintf("%d:sol:lte", maxSol))
	}
	return result
}

func parseImages(body []byte) ([]ImageInfo, error) {
	result := []ImageInfo{}
	// https://mariadesouza.com/2017/09/07/custom-unmarshal-json-in-golang/
	var raw map[string]*json.RawMessage
	err := json.Unmarshal(body, &raw)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(*raw["images"], &result)
	if err != nil {
		debug.PrintStack()
	}
	return result, err
}

func GetImageMetadata(params url.Values) ([]ImageInfo, error) {
	apiUrl := "https://mars.nasa.gov/rss/api/"
	fullUrl := apiUrl + "?" + params.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Failed getting:", fullUrl, "reason", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug.PrintStack()
		return []ImageInfo{}, err
	}
	return parseImages(body)
}
