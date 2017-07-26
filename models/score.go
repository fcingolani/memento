package models

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

type Score struct {
	ID int64 `json:"id"`

	PlayerName   string `json:"playerName" form:"player_name" query:"player_name" validate:"omitempty,alphanum"`
	LevelNumber  int    `json:"levelNumber" form:"level_number" query:"level_number" validate:"required,gt=0"`
	LevelVersion int    `json:"levelVersion" form:"level_version" query:"level_version" validate:"required,gt=0"`
	Value        int    `json:"value" form:"value" query:"value" validate:"required,gt=0"`

	File *File `json:"file,omitempty"`
}

const (
	HigherScore = 1
	LowerScore  = -1
)

func (ds *datastore) FindScores() (*[]Score, error) {

	q := "SELECT id, player_name, level_number, level_version, value FROM scores"

	res, err := ds.Query(q)

	if err != nil {
		return nil, err
	}

	ss := []Score{}

	for res.Next() {
		var s Score
		err = res.Scan(&s.ID, &s.PlayerName, &s.LevelNumber, &s.LevelVersion, &s.Value)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}

	return &ss, nil

}

func (ds *datastore) CreateScore(s *Score) error {

	s.File = &File{UploadTicket: uuid.NewV4()}

	q := "INSERT INTO scores (player_name, level_number, level_version, value, file_upload_ticket) VALUES (?,?,?,?,?)"
	res, err := ds.Exec(q, s.PlayerName, s.LevelNumber, s.LevelVersion, s.Value, s.File.UploadTicket)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	s.ID = id
	s.File.ID = id

	return nil

}

func (ds *datastore) FindBeatableScore(s *Score, t int, f bool) (*Score, error) {

	q := "SELECT id, player_name, level_number, level_version, value FROM scores WHERE level_number = ? AND level_version = ? AND value < ?"

	if f == true {
		q += " AND file_data NOT NULL"
	}

	if t == LowerScore {
		q += " ORDER BY value DESC LIMIT 1"
	} else {
		q += " ORDER BY value ASC LIMIT 1"
	}

	rows, err := ds.Query(q, s.LevelNumber, s.LevelVersion, s.Value)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("couldn't find a beatable score")
	}

	b := &Score{}

	if err = rows.Scan(&b.ID, &b.PlayerName, &b.LevelNumber, &b.LevelVersion, &b.Value); err != nil {
		return nil, err
	}

	return b, nil
}
