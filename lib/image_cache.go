package lib

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ImageCache struct {
	idb     ImageDB
	rootdir string
}

const DefaultCachePathname = "./image_cache"

const thumbDir = "thumbnail"
const fullDir = "full_res"

func NewImageCache(idb ImageDB) ImageCache {
	return ImageCache{idb, DefaultCachePathname}
}

func (cache *ImageCache) ThumbDir() string {
	return filepath.Join(cache.rootdir, thumbDir)
}

func (cache *ImageCache) ThumbPath(imageID string) string {
	return filepath.Join(cache.rootdir, thumbDir, imageID+".png")
}

func (cache *ImageCache) FullSizeDir() string {
	return filepath.Join(cache.rootdir, fullDir)
}

func (cache *ImageCache) FullSizePath(imageID string) string {
	return filepath.Join(cache.rootdir, fullDir, imageID+".png")
}

func ensureDirExists(pathname string) error {
	abspath, err := filepath.Abs(pathname)
	if err != nil {
		return err
	}
	dirpath := filepath.Dir(abspath)
	return os.MkdirAll(dirpath, 0775)
}

func imageData(pathname string) (image.Image, error) {
	var result image.Image = image.NewRGBA(image.Rect(0, 0, 0, 0))
	reader, err := os.Open(pathname)
	if err != nil {
		return result, err
	}
	defer reader.Close()

	result, _, err = image.Decode(reader)
	return result, err
}

func downloadImage(url string, destpath string) (image.Image, error) {
	var result image.Image = image.NewRGBA(image.Rect(0, 0, 0, 0))
	response, err := http.Get(url)
	if err != nil {
		return result, err
	}

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	byteReader := bytes.NewReader(imageData)
	result, formatName, err := image.Decode(byteReader)
	if err != nil {
		return result, err
	}
	fmt.Println("Decoded image data with format", formatName)

	ensureDirExists(destpath)
	// Save all images as PNG.
	fileWriter, err := os.Create(destpath)
	if err != nil {
		return result, err
	}
	defer fileWriter.Close()

	err = png.Encode(fileWriter, result)
	return result, err
}

func (cache *ImageCache) ThumbNail(imageID string) (image.Image, error) {
	imagePathname := cache.ThumbPath(imageID)
	result, err := imageData(imagePathname)
	if err == nil {
		return result, err
	}
	url, err := cache.idb.ThumbnailURL(imageID)
	if err != nil {
		return result, err
	}

	return downloadImage(url, imagePathname)
}

func (cache *ImageCache) FullSize(imageID string) (image.Image, error) {
	imagePathname := cache.FullSizePath(imageID)
	result, err := imageData(imagePathname)
	if err == nil {
		return result, err
	}
	url, err := cache.idb.FullSizeURL(imageID)
	if err != nil {
		return result, err
	}

	return downloadImage(url, imagePathname)
}
