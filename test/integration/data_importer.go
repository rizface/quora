package integration

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ImportSQL(db *sql.DB, path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	q, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(q))
	if err != nil {
		log.Fatal(err)
	}
}
