package lib

import (
	"io/ioutil"
	"testing"
)

func TestParseImageMetadata(t *testing.T) {
	// data, err := fs.ReadFile("test_data/sample_rss_response.json")
	data, err := ioutil.ReadFile("test_data/sample_rss_response.json")
	if len(data) <= 0 {
		t.Errorf("Failed to read test JSON file: %v", err)
	}
	got, err := ParseImageMetadata(data)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	want := 100
	if len(got) != want {
		t.Errorf("Expected %v image records, got %v", want, len(got))
	}
}
