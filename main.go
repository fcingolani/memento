package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

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

	dp := env("DATABASE_PATH", "./db.sqlite")
	addr := env("LISTEN_ADDRESS", ":3000")
	debug := env("DEBUG", "false")
	corsOrigins := env("CORS_ALLOWED_ORIGINS", "")

	mub, err := strconv.ParseInt(env("MAX_UPLOAD_BYTES", "1048576"), 10, 64)
	if err != nil {
		panic(err)
	}

	ds, err := models.NewDatastore(dp)

	if err != nil {
		panic(err)
	}

	sc := controllers.NewScoreController(ds)
	fc := controllers.NewFileController(ds, mub)

	e := echo.New()

	e.Debug = debug == "true"
	e.HideBanner = true

	e.Validator = &CustomValidator{validator: validator.New()}

	if corsOrigins != "" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: strings.Split(corsOrigins, ","),
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}))
	}

	e.GET("/scores", sc.List)
	e.POST("/scores", sc.Add)
	e.GET("/scores/_beatable", sc.Beatable)

	e.GET("/files/:id/data", fc.GetData)
	e.PUT("/files/:id/data", fc.SaveData)

	e.GET("/check", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(addr))

}
