package auth

import (
	validation "bitbucket.org/2tgroup/ciwp-api-users/common/validations"
	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/users"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	validator "gopkg.in/go-playground/validator.v9"
)

func RoutersAuth() *echo.Echo {
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

	routerAuth := e.Group("/auth")

	routerAuth.POST("/register", UserRegisterHandler)

	routerAuth.POST("/login", UserLoginHandler)

	// JWT
	routerAuth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &users.AuthJwtClaims{},
		SigningKey: []byte(config.DataConfig.SecretKey),
	}))

	routerAuth.GET("/token", UserTokenHandler)

	return e
}
