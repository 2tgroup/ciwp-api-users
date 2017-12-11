package auth

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	validator "gopkg.in/go-playground/validator.v9"

	validation "bitbucket.org/2tgroup/ciwp-api-users/common/validations"
	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

func init() {
	log.SetPrefix("modules/auth")
	log.Info("Loaded modules/auth")
}
func RoutersAuth() *echo.Echo {
	// Echo instance
	e := echo.New()
	// Customization
	if config.DataConfig.ReleaseMode {
		e.Debug = false
	}
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

	routerAuth.POST("/forgot", UserForgotHandler)

	// JWT
	routerAuth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &types.AuthJwtClaims{},
		SigningKey: []byte(config.DataConfig.SecretKey),
	}))

	// need way to store old token to backlist
	routerAuth.GET("/logout", UserLogoutHandler)

	routerAuth.GET("/token", UserTokenHandler)

	return e
}
