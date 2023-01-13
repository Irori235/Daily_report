package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	c := e.NewContext(nil, nil)

	apiToken := os.Getenv("TOGGL_API_TOKEN")
	baseURL := os.Getenv("TOGGL_BASE_URL")

	Toggl := NewToggl(apiToken, baseURL)

	log, err := Toggl.GetTimeLogs(c)
}
