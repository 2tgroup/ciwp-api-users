package users

import (
	validation "bitbucket.org/2tgroup/ciwp-api-users/common/validations"
	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/auth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

func Routers() *echo.Echo {
	// Echo instance
	e := echo.New()
	// Customization
	if config.DataConfig.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("Auth")
	// CSRF
	/* e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:" + echo.HeaderXCSRFToken,
	})) */

	// Gzip
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	//VALIDATION
	e.Validator = &validation.CustomValidator{ValidatorX: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routers

	router := e.Group("/users")
	// JWT
	router.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.AuthJwtClaims{},
		SigningKey: []byte(config.DataConfig.SecretKey),
	}))
	router.POST("/update", UserUpdateHandler)

	return e
}
