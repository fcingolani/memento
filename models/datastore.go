package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

type IDatastore interface {
	CreateReplay(*Replay) error
	FindBeatableReplay(int64, int64, int64) (*Replay, error)
	SaveReplayData(id int64, fut uuid.UUID, fud []byte) error
	GetReplayData(id int64) ([]byte, error)
}

type datastore struct {
	*sql.DB
}

func NewDatastore(p string) (IDatastore, error) {
	db, err := sql.Open("sqlite3", p)

	if err != nil {
		return nil, err
	}

	q := `CREATE TABLE IF NOT EXISTS [replays] (
		[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
		[player_name] VARCHAR(20)  NULL,
		[level_number] INTEGER  NOT NULL,
		[level_version] INTEGER  NOT NULL,
		[time] INTEGER  NOT NULL,
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
