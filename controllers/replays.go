package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fcingolani/memento/models"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type IReplayController interface {
	AddReplay(c echo.Context) error
	GetBeatableReplay(c echo.Context) error
	GetReplayFile(c echo.Context) error
	SaveReplayFile(c echo.Context) error
}

type replayController struct {
	ds             models.IDatastore
	maxUploadBytes int64
}

func NewReplayController(ds models.IDatastore, mub int64) IReplayController {
	return &replayController{ds, mub}
}

func (rc *replayController) AddReplay(c echo.Context) error {

	r := new(models.Replay)

	if err := c.Bind(r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := c.Validate(r); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := rc.ds.CreateReplay(r); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)
}

func (rc *replayController) SaveReplayFile(c echo.Context) error {

	r := http.MaxBytesReader(c.Response(), c.Request().Body, rc.maxUploadBytes)

	b, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}

	fut, err := uuid.FromString(c.Request().Header.Get("x-file-upload-ticket"))

	if err != nil {
		return err
	}

	if err := rc.ds.SaveReplayData(id, fut, b); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)

}

func (rc *replayController) GetBeatableReplay(c echo.Context) error {

	ln, err := strconv.ParseInt(c.Param("level_number"), 10, 64)

	if err != nil {
		return err
	}

	lv, err := strconv.ParseInt(c.Param("level_version"), 10, 64)

	if err != nil {
		return err
	}

	t, err := strconv.ParseInt(c.Param("time"), 10, 64)

	if err != nil {
		return err
	}

	r, err := rc.ds.FindBeatableReplay(ln, lv, t)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)
}

func (rc *replayController) GetReplayFile(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return err
	}

	b, err := rc.ds.GetReplayData(id)

	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/octet-stream", b)
}
