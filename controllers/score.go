package controllers

import (
	"net/http"

	"github.com/fcingolani/memento/models"
	"github.com/labstack/echo"
)

type IScoreController interface {
	List(c echo.Context) error
	Add(c echo.Context) error
	Beatable(c echo.Context) error
}

type scoreController struct {
	ds models.IDatastore
}

func NewScoreController(ds models.IDatastore) IScoreController {
	return &scoreController{ds}
}

func (sc *scoreController) List(c echo.Context) error {

	ss, err := sc.ds.FindScores()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ss)

}

func (sc *scoreController) Add(c echo.Context) error {

	s := &models.Score{}

	if err := c.Bind(s); err != nil {
		return err
	}

	if err := c.Validate(s); err != nil {
		return err
	}

	if err := sc.ds.CreateScore(s); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func (sc *scoreController) Beatable(c echo.Context) error {

	s := &models.Score{}

	if err := c.Bind(s); err != nil {
		return err
	}

	if err := c.Validate(s); err != nil {
		return err
	}

	var t int

	if c.QueryParam("type") == "lower" {
		t = models.LowerScore
	} else {
		t = models.HigherScore
	}

	b, err := sc.ds.FindBeatableScore(s, t)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}
