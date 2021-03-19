package lib

import (
	"fmt"
	"os"
	"testing"
)

func TestGetAThumbnail(t *testing.T) {
	idb, err := recreateDB()
	if err != nil {
		t.Fatal("Error recreating db.", err)
	}
	defer os.Remove(idb.DBName)

	loadSampleData(idb, t)

	cache, err := NewImageCache(idb)
	if err != nil {
		t.Fatal("Error creating image cache:", err)
	}
	defer os.RemoveAll(cache.rootdir)

	imageID := "SI0_0024_0669080939_077ECM_N0030792SRLC07015_0000LUJ"

	image, err := cache.ThumbNail(imageID)
	if err != nil {
		t.Fatal("Error loading thumbnail image:", err)
	}

	bounds := image.Bounds()
	fmt.Println("Got thumbnail with bounds", bounds)
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Fatal("Implausible thumbnail extents ", bounds.Dx(), "x", bounds.Dy())
	}
}

func TestGetAFullRes(t *testing.T) {
	idb, err := recreateDB()
	if err != nil {
		t.Fatal("Error recreating db.", err)
	}
	defer os.Remove(idb.DBName)

	loadSampleData(idb, t)

	cache, err := NewImageCache(idb)
	if err != nil {
		t.Fatal("Error creating image cache:", err)
	}
	defer os.RemoveAll(cache.rootdir)

	imageID := "NLF_0024_0669080250_161ECM_N0030792NCAM00194_01_290J"

	image, err := cache.FullSize(imageID)
	if err != nil {
		t.Fatal("Error loading full size image:", err)
	}

	bounds := image.Bounds()
	fmt.Println("Got full-res image with bounds", bounds)
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Fatal("Implausible image extents ", bounds.Dx(), "x", bounds.Dy())
	}
	// TODO verify that subsequent retrievals use data from cache.
}
