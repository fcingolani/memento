package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type IDatastore interface {
	CreateScore(*Score) error
	FindBeatableScore(*Score, int) (*Score, error)
	SaveFile(*File) error
	GetFileById(id int64) (*File, error)
}

type datastore struct {
	*sql.DB
}

func NewDatastore(p string) (IDatastore, error) {
	db, err := sql.Open("sqlite3", p)

	if err != nil {
		return nil, err
	}

	q := `CREATE TABLE IF NOT EXISTS [scores] (
		[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
		[player_name] VARCHAR(20)  NULL,
		[level_number] INTEGER  NOT NULL,
		[level_version] INTEGER  NOT NULL,
		[value] INTEGER  NOT NULL,
		[file_data] BLOB  NULL,
		[file_upload_ticket] VARCHAR(40)  NULL,
		[created_at] TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
		[updated_at] TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	)`

	_, err = db.Exec(q)

	if err != nil {
		return nil, err
	}

	return &datastore{db}, nil
}
