package lib

import (
	"database/sql"
	"os"
)

type ImageDB struct {
	db *sql.DB
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
		cam_position TEXT NOT NULL,

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
	result.db, err = sql.Open("sqlite", "./mars_perseverance_image_info.db")
	if err != nil {
		return result, err
	}
	return result, result.initSchema()
}
