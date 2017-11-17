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
	//e.GET("/login", UserLoginHandler)
	//e.GET("/register", UserRegisterHandler)

	// github.com/hb-go/json
	//e.GET("/json/encode", handler(JsonEncodeHandler))

	// JWT
	auth := e.Group("/auth")

	auth.POST("/register", UserRegisterHandler)

	auth.POST("/login", UserLoginHandler)

	auth.GET("/token", UserTokenHandler)

	/* auth.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &AuthJwtClaims{},
		SigningKey: []byte(config.DataConfig.SecretKey),
	})) */

	//r.GET("/", handler(ApiHandler))

	// curl http://echo.api.localhost:8080/restricted/user -H "Authorization: Bearer XXX"
	//r.GET("/user", UserHandler)
	return e
}

type (
	HandlerFunc func(*echo.Context) error
)

/**
 * handler Request
 */
func handler(h HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(echo.Context)
		return h(&ctx)
	}
}
