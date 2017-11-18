package auth

import (
	validation "bitbucket.org/2tgroup/ciwp-api-users/common/validations"
	"bitbucket.org/2tgroup/ciwp-api-users/config"
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

	auth := e.Group("/auth")

	auth.POST("/register", UserRegisterHandler)

	auth.POST("/login", UserLoginHandler)

	// JWT
	auth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &AuthJwtClaims{},
		SigningKey: []byte(config.DataConfig.SecretKey),
	}))

	auth.GET("/token", UserTokenHandler)

	return e
}
