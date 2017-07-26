package models

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

type File struct {
	ID           int64     `json:"id,omitempty"`
	UploadTicket uuid.UUID `json:"uploadTicket,omitempty"`
	Data         []byte    `json:"-"`
}

func (ds *datastore) SaveFile(f *File) error {

	q := "UPDATE scores SET file_data = ?, file_upload_ticket = NULL WHERE id = ? AND file_upload_ticket = ? AND file_upload_ticket NOT NULL"
	res, err := ds.Exec(q, f.Data, f.ID, f.UploadTicket)

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

func (ds *datastore) GetFileById(id int64) (*File, error) {

	q := "SELECT file_data FROM scores WHERE id = ?"
	rows, err := ds.Query(q, id)

	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("couldn't find replay data")
	}

	f := &File{ID: id}

	if err := rows.Scan(&f.Data); err != nil {
		return nil, err
	}

	return f, nil

}
