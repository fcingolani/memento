package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fcingolani/memento/models"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type IFileController interface {
	GetData(c echo.Context) error
	SaveData(c echo.Context) error
}

type fileController struct {
	ds             models.IDatastore
	maxUploadBytes int64
}

func NewFileController(ds models.IDatastore, mub int64) IFileController {
	return &fileController{ds, mub}
}

func (fc *fileController) SaveData(c echo.Context) error {

	r := http.MaxBytesReader(c.Response(), c.Request().Body, fc.maxUploadBytes)

	b, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}

	t, err := uuid.FromString(c.Request().Header.Get("x-file-upload-ticket"))

	if err != nil {
		return err
	}

	f := &models.File{Data: b, ID: id, UploadTicket: t}

	if err := fc.ds.SaveFile(f); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)

}

func (fc *fileController) GetData(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}

	f, err := fc.ds.GetFileById(id)

	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/octet-stream", f.Data)
}
