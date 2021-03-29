package lib

import "os"

func FileExists(pathname string) bool {
	_, err := os.Stat(pathname)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// Something else could be amiss here.
	return false
}
