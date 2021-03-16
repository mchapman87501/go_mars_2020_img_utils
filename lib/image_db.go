package lib

import (
	"database/sql"
	"fmt"
	"math"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type ImageDB struct {
	db *sql.DB
	// TODO add placeholders for prepared statements, for common ops such as
	// insertion.
}

const dbFilename = "./mars_perseverance_image_info.db"

const schema string = `CREATE TABLE IF NOT EXISTS Images (
		image_id TEXT NOT NULL PRIMARY KEY,

		credit TEXT NOT NULL,
		caption TEXT NOT NULL,
		title TEXT NOT NULL,

		cam_instrument TEXT NOT NULL,
		cam_filter TEXT NOT NULL,
		cam_model_component_list TEXT NOT NULL,
		cam_model_type TEXT NOT NULL,
		
		cam_pos_x REAL,
		cam_pos_y REAL,
		cam_pos_z REAL,

		sample_type TEXT NOT NULL,

		small_url TEXT,
		full_res_url TEXT,
		json_url TEXT,

		date_taken_utc TIMESTAMP NOT NULL,
		-- date_taken_mars TIMESTAMP NOT NULL,
		-- date_received TIMESTAMP NOT NULL,
		-- sol INTEGER NOT NULL,

		-- misc
		attitude TEXT NOT NULL, -- 3-tuple of floats, I think
		drive INTEGER,
		site INTEGER,

		-- extended properties:
		ext_mast_azimuth REAL,
		ext_mast_elevation REAL,
		ext_sclk REAL,
		ext_scale_factor REAL,

		-- position?  What coordinates?
		ext_x REAL,
		ext_y REAL,
		ext_z REAL,

		-- subframe rect:
		ext_sf_left REAL,
		ext_sf_top REAL,
		ext_sf_width REAL,
		ext_sf_height REAL,

		-- dimension: (width, height), appears to be image size in pixels
		ext_width REAL,
		ext_height REAL
	);`

func (idb *ImageDB) initSchema() error {
	_, err := idb.db.Exec(schema)
	return err
}

// Creates a new image database in the current working directory.
func NewImageDB() (ImageDB, error) {
	result := ImageDB{}
	file, err := os.Create(dbFilename)
	file.Close()

	if err != nil {
		return result, err
	}
	result.db, err = sql.Open("sqlite3", "./mars_perseverance_image_info.db")
	if err != nil {
		return result, err
	}
	return result, result.initSchema()
}

// Add or update Images from provided records.
func (idb *ImageDB) AddOrUpdate(records []ImageInfo) error {
	statement, err := idb.prepareUpdateOne()

	if err != nil {
		return err
	}
	defer statement.Close()

	for _, record := range records {
		err := addOrUpdateOne(statement, record)
		if err != nil {
			return err
		}
	}
	return nil
}

func (idb *ImageDB) prepareUpdateOne() (*sql.Stmt, error) {
	// SQLite3 supports named query parameters.  Go's sql.DB support
	// for named parameters looks a bit verbose to me.
	// https://golang.org/pkg/database/sql/#Named
	query := `INSERT OR REPLACE INTO Images
	(
		image_id, credit, caption, title,
		cam_instrument, cam_filter, cam_model_component_list,
		cam_model_type,
		cam_pos_x, cam_pos_y, cam_pos_z,
		sample_type,
		full_res_url, json_url,
		date_taken_utc,
		attitude, drive, site,
		ext_mast_azimuth, ext_mast_elevation,
		ext_sclk,
		ext_scale_factor,
		ext_x, ext_y, ext_z,
		ext_sf_left, ext_sf_top, ext_sf_width, ext_sf_height,
		ext_width, ext_height
	) VALUES (
		?, ?, ?, ?,
		?, ?, ?,
		?,
		?, ?, ?,
		?,
		?, ?,
		?,
		?, ?, ?,
		?, ?,
		?,
		?,
		?, ?, ?,
		?, ?, ?, ?,
		?, ?
	)`
	return idb.db.Prepare(query)
}

func addOrUpdateOne(statement *sql.Stmt, record ImageInfo) error {
	nan := math.NaN()
	// Pad out arrays/slices so they have at least the expected length.
	camPos := append(record.Camera.CameraPosition, nan, nan, nan)

	// TODO learn how Go sql drivers might convert to/from []float64, etc.
	attitudeStr := fmt.Sprint(record.Attitude)

	extXYZ := append(record.Extended.XYZ, nan, nan, nan)

	_, err := statement.Exec(
		record.ImageId,
		record.Credit,
		record.Caption,
		record.Title,
		record.Camera.Instrument,
		record.Camera.FilterName,
		record.Camera.CameraModelComponentList,
		record.Camera.CameraModelType,
		camPos[0], camPos[1], camPos[2],
		record.SampleType, record.ImageFiles.FullRes, record.JsonLink,
		record.DateTakenUtc, attitudeStr, record.Drive, record.Site,
		record.Extended.MastAzimuth, record.Extended.MastElevation,
		record.Extended.Sclk, record.Extended.ScaleFactor,
		extXYZ[0], extXYZ[1], extXYZ[2],
		record.Extended.SubframeRect.Origin.X, record.Extended.SubframeRect.Origin.Y,
		record.Extended.SubframeRect.Size.Width, record.Extended.SubframeRect.Size.Height,
		record.Extended.Dimension.Width, record.Extended.Dimension.Height,
	)
	return err
}

func (idb *ImageDB) DB() *sql.DB {
	return idb.db
}