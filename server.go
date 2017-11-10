package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/2tgroup/ciwp-api-users/config"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		configEnv := config.GetConfig()

		fmt.Println("WE GO?")

		return c.JSON(http.StatusOK, configEnv)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
