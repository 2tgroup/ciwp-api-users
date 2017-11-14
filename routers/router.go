package router

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/auth"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

/*Hold all endpoint*/
type (
	Host struct {
		Echo *echo.Echo
	}
)

var HostNames = make(map[string]*Host)

func init() {
	HostNames[config.DataConfig.Server["api_user"]] = &Host{auth.Routers()}
	//HostNames[Conf.Server.DomainWeb] = &Host{web.Routers()}
}

func GetInfoEndpoint() map[string]*Host {
	return HostNames
}

/*Init Run Router*/
func InitRouter() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.SecureWithConfig(middleware.DefaultSecureConfig))
	middleware.MethodOverride()
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://" + config.DataConfig.Server["api_user"]},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, echo.HeaderAuthorization},
	}))

	// SEND REQUEST TO ENDPOINT
	e.Any("/*", func(c echo.Context) (err error) {
		req, res := c.Request(), c.Response()
		u, _err := url.Parse(c.Scheme() + "://" + req.Host)
		if _err != nil {
			e.Logger.Errorf("Request URL parse error:%v", _err)
		}
		host := HostNames[u.Hostname()]
		if host == nil {
			e.Logger.Info("Host not found")
			err = echo.ErrNotFound
		} else {

			host.Echo.ServeHTTP(res, req)
		}
		return
	})

	if config.DataConfig.Server["graceful"] != "true" {
		e.Logger.Fatal(e.Start(config.DataConfig.Host))
	} else {
		// Start server
		go func() {
			if err := e.Start(config.DataConfig.Host); err != nil {
				e.Logger.Errorf("Shutting down the server with error:%v", err)
			}
		}()
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}

}
