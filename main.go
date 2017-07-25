package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"

	"gopkg.in/go-playground/validator.v9"

	"github.com/fcingolani/memento/controllers"
	"github.com/fcingolani/memento/models"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func env(k, d string) string {
	v := os.Getenv(k)

	if len(v) == 0 {
		return d
	}

	return v
}

func main() {
	e := echo.New()

	dp := env("DATABASE_PATH", "./db.sqlite")
	addr := env("LISTEN_ADDRESS", ":3000")
	debug := env("DEBUG", "false")

	mub, err := strconv.ParseInt(env("MAX_UPLOAD_BYTES", "1048576"), 10, 64)
	if err != nil {
		panic(err)
	}

	e.Debug = debug == "true"
	e.HideBanner = true

	e.Validator = &CustomValidator{validator: validator.New()}

	ds, err := models.NewDatastore(dp)

	if err != nil {
		panic(err)
	}

	rc := controllers.NewReplayController(ds, mub)

	e.POST("/replays", rc.AddReplay)
	e.PUT("/replays/:id/file", rc.SaveReplayFile)
	e.GET("/replays/:id/file", rc.GetReplayFile)
	e.GET("/replays/_tobeat/:level_number/:level_version/:time", rc.GetBeatableReplay)

	e.GET("/check", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(addr))

}
