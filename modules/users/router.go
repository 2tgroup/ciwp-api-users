package users

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

func init() {
	log.Info("Loaded module USERS")
}

func RoutersUser() *echo.Echo {
	// Echo instance
	e := echo.New()
	// Customization
	if config.DataConfig.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("User")
	// CSRF
	/* e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:" + echo.HeaderXCSRFToken,
	})) */

	// Gzip
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	//VALIDATION
	//e.Validator = &validation.CustomValidator{ValidatorX: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routers
	routerUser := e.Group("/users")
	routerUser.GET("/country", UserListCountryByAlpha2)
	routerUser.GET("/country/:alpha2", UserListCounty)
	// JWT
	routerUser.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &types.AuthJwtClaims{},
		SigningKey: []byte(config.DataConfig.SecretKey),
	}))
	routerUser.POST("/update", UserUpdateHandler)

	return e
}
