package lib

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func recreateDB() (ImageDB, error) {
	_, err := os.Stat(dbFilename)
	if err == nil {
		os.Remove(dbFilename)
	}

	return NewImageDB()
}

func TestNewImageDB(t *testing.T) {
	_, err := recreateDB()

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

func TestAddRecords(t *testing.T) {
	idb, err := recreateDB()
	if err != nil {
		t.Errorf("Error creating Image DB: %v", err)
	} else {
		defer os.Remove(dbFilename)

		data, err := ioutil.ReadFile("test_data/sample_rss_response.json")
		if len(data) <= 0 {
			t.Errorf("Failed to read test JSON file: %v", err)
		}

		records, err := ParseImageMetadata(data)
		if err != nil {
			t.Errorf("Error parsing image metadata: %v", err)
		}

		err = idb.AddOrUpdate(records)
		if err != nil {
			t.Errorf("Error adding/updating DB records: %v", err)
		}

		stmt, err := idb.DB().Prepare("SELECT COUNT(*) FROM Images")
		if err != nil {
			t.Errorf("Error preparing count query: %v", err)
		}
		defer stmt.Close()

		countRow := stmt.QueryRow()
		got := -1
		err = countRow.Scan(&got)
		if err != nil {
			t.Errorf("Error retrieving count: %v", err)
		}
		want := len(records)
		if want != got {
			t.Errorf("Expected database to have # rows = %v, got %v", want, got)
		}
	}
}
