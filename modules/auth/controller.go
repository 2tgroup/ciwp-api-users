package auth

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
)

//UserTokenHandler is
func UserTokenHandler(c echo.Context) error {
	claims := &AuthJwtClaims{}
	claims.Name = "TPT - Svideo"
	claims.Email = "support@serverapi.host"
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}
	t, e := getJETToken(claims)
	if e != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"token": e,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func UserLoginHandler(c echo.Context) error {

	claims := &AuthJwtClaims{}
	claims.Name = "TPT - Svideo"
	claims.Email = "support@serverapi.host"
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.DataConfig.SecretKey))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"token": err,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func UserRegisterHandler(c echo.Context) error {

	c.JSON(200, map[string]interface{}{"URI": "api user regist"})
	return nil
}

func getJETToken(Authclaims *AuthJwtClaims) (t string, e error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Authclaims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.DataConfig.SecretKey))
	if err != nil {
		return "", err
	}
	return t, err
}
