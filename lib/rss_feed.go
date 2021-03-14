package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

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

func GetImageMetadata(params url.Values) []byte {
	apiUrl := "https://mars.nasa.gov/rss/api/"
	fullUrl := apiUrl + "?" + params.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("Failed getting:", fullUrl, "reason", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed reading response body:", err)
	}
	return body
}
