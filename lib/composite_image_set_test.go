package lib

import (
	"fmt"
	"io/ioutil"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TODO extract these to a test helper file.
func loadDB(idb ImageDB, jsonPath string, t *testing.T) {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil || len(data) <= 0 {
		t.Fatal("Failed to read test JSON file:", err)
	}

	records, err := ParseImageMetadata(data)
	if err != nil {
		t.Fatal("Error parsing image metadata:", err)
	}

	fmt.Println("Record count:", len(records))
	if err = idb.AddOrUpdate(records); err != nil {
		t.Fatal("Error adding/updating DB records:", err)
	}
}

func loadFakeData(idb ImageDB, t *testing.T) {
	loadDB(idb, "test_data/with_fake_pano_set.json", t)
}

func loadSampleData(idb ImageDB, t *testing.T) {
	loadDB(idb, "test_data/sample_rss_response.json", t)
}

func recreateInMemDB(t *testing.T) ImageDB {
	idb, err := NewImageDBAtPath(":memory:")
	if err != nil {
		t.Fatal("Could not create in-memory database:", err)
	}

	loadFakeData(idb, t)
	return idb
}

func getCIRecords(camera string, t *testing.T) []CompositeImageInfo {
	idb := recreateInMemDB(t)
	records, err := GetCompositeImageInfoRecords(idb, camera)
	if err != nil {
		t.Fatal(err)
	}
	return records
}

func TestGetCompositeImageInfoRecords(t *testing.T) {
	records := getCIRecords("invalid camera", t)
	if len(records) > 0 {
		t.Errorf("Expected no image info records, got %v", records)
	}

	records = getCIRecords("NAVCAM_LEFT", t)
	if len(records) <= 0 {
		t.Errorf("Expected some image info records, got none.")
	}
}

func getCISets(camera string, t *testing.T) []CompositeImageSet {
	idb := recreateInMemDB(t)
	imageSets, err := GetCompositeImageSets(idb, camera)
	if err != nil {
		t.Fatal(err)
	}
	return imageSets
}

func TestGetCompositeImageSets(t *testing.T) {
	records := getCISets("invalid camera", t)
	if len(records) > 0 {
		t.Errorf("Expected no image sets, got %v", records)
	}

	cam := "NAVCAM_LEFT"
	records = getCISets(cam, t)
	got := len(records)
	want := 1
	if got != want {
		t.Errorf("Expected # image sets: %v, got %v.", want, got)
	}

	for _, imageRecords := range records {
		if len(imageRecords) <= 1 {
			t.Errorf("Image sets must contain multiple images.")
		}
		for _, imageRecord := range imageRecords {
			if imageRecord.Camera != cam {
				t.Fatalf("Image set record contains unexpected camera %v", imageRecord.Camera)
			}
			fmt.Println(imageRecord.ImageID, imageRecord.SubframeRect)
		}
	}
}

// TODO add sorting test.

// TODO verify that all sets returned by GetCompositeImageSets have multiple images.
