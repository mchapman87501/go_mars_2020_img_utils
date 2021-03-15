package lib

import (
	"errors"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestNewImageDB(t *testing.T) {
	_, err := os.Stat(dbFilename)
	if err == nil {
		os.Remove(dbFilename)
	}

	_, err = NewImageDB()
	if err != nil {
		t.Errorf("Error creating Image DB: %v", err)
	} else {
		_, err := os.Stat(dbFilename)
		if errors.Is(err, os.ErrNotExist) {
			t.Errorf("Expected NewImageDB to create %v, but it did not.", dbFilename)
		} else {
			defer os.Remove(dbFilename)
		}
	}

	// TODO Verify that the database schema was created.
}
