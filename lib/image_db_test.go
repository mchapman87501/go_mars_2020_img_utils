package lib

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func recreateDB() (ImageDB, error) {
	_, err := os.Stat(DefaultDBPathname)
	if err == nil {
		os.Remove(DefaultDBPathname)
	}

	return NewImageDB()
}

func TestNewImageDB(t *testing.T) {
	_, err := recreateDB()

	if err != nil {
		t.Fatal("Error creating Image DB:", err)
	} else {
		_, err := os.Stat(DefaultDBPathname)
		if errors.Is(err, os.ErrNotExist) {
			t.Fatal("Expected NewImageDB to create", DefaultDBPathname, "but it did not.")
		} else {
			defer os.Remove(DefaultDBPathname)
		}
	}
}

func TestAddRecords(t *testing.T) {
	idb, err := recreateDB()
	if err != nil {
		t.Fatal("Error creating Image DB:", err)
	} else {
		defer os.Remove(DefaultDBPathname)

		data, err := ioutil.ReadFile("test_data/sample_rss_response.json")
		if err != nil || len(data) <= 0 {
			t.Fatal("Failed to read test JSON file:", err)
		}

		records, err := ParseImageMetadata(data)
		if err != nil {
			t.Fatal("Error parsing image metadata:", err)
		}

		if err = idb.AddOrUpdate(records); err != nil {
			t.Fatal("Error adding/updating DB records:", err)
		}

		stmt, err := idb.DB.Prepare("SELECT COUNT(*) FROM Images")
		if err != nil {
			t.Fatal("Error preparing count query:", err)
		}
		defer stmt.Close()

		countRow := stmt.QueryRow()
		got := -1
		if err = countRow.Scan(&got); err != nil {
			t.Fatal("Error retrieving count:", err)
		}
		want := len(records)
		if want != got {
			t.Errorf("Expected database to have # rows = %v, got %v", want, got)
		}
	}
}
