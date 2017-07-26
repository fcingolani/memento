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

func (ds *datastore) FindBeatableScore(s *Score, t int) (*Score, error) {

	println(t)

	var q string

	if t == LowerScore {
		q = "SELECT id, player_name, level_number, level_version, value FROM scores WHERE level_number = ? AND level_version = ? AND value < ? ORDER BY value DESC LIMIT 1"
	} else {
		q = "SELECT id, player_name, level_number, level_version, value FROM scores WHERE level_number = ? AND level_version = ? AND value > ? ORDER BY value ASC LIMIT 1"
	}

	rows, err := ds.Query(q, s.LevelNumber, s.LevelVersion, s.Value)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("couldn't find a beatable score")
	}

	b := new(Score)

	if err = rows.Scan(&b.ID, &b.PlayerName, &b.LevelNumber, &b.LevelVersion, &b.Value); err != nil {
		return nil, err
	}

	return b, nil
}

/*
func (ds *datastore) CreateReplay(r *Replay) error {

	r.FileUploadTicket = uuid.NewV4()

	q := "INSERT INTO replays (player_name, level_number, level_version, time, file_upload_ticket) VALUES (?,?,?,?,?)"
	res, err := ds.Exec(q, r.PlayerName, r.LevelNumber, r.LevelVersion, r.Time, r.FileUploadTicket)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	r.ID = id

	return nil

}

func (ds *datastore) FindBeatableReplay(ln, lv, t int64) (*Replay, error) {

	q := "SELECT id, player_name, level_number, level_version, time FROM replays WHERE level_number = ? AND level_version = ? AND time < ? ORDER BY time ASC LIMIT 1"
	rows, err := ds.Query(q, ln, lv, t)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("couldn't find a beatable replay")
	}

	replay := new(Replay)

	if err = rows.Scan(&replay.ID, &replay.PlayerName, &replay.LevelNumber, &replay.LevelVersion, &replay.Time); err != nil {
		return nil, err
	}

	return replay, nil
}

func (ds *datastore) SaveReplayData(id int64, fut uuid.UUID, fd []byte) error {

	q := "UPDATE replays SET file_data = ?, file_upload_ticket = NULL WHERE id = ? AND file_upload_ticket = ? AND file_upload_ticket NOT NULL"
	res, err := ds.Exec(q, fd, id, fut)

	if err != nil {
		return err
	}

	r, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if r != 1 {
		return errors.New("error while saving data")
	}

	return nil

}

func (ds *datastore) GetReplayData(id int64) ([]byte, error) {

	q := "SELECT file_data FROM replays WHERE id = ?"
	rows, err := ds.Query(q, id)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("couldn't find replay data")
	}

	var b []byte

	if err := rows.Scan(&b); err != nil {
		return nil, err
	}

	return b, nil

}
*/
