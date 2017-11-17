package auth

import (
	"fmt"
	"net/http"
	"time"

	"g.ghn.vn/go-training/tientp/types"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/users"
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

	u := new(users.UserBase)

	if err := c.Bind(u); err != nil {
		log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponse("request_invaild", "Có lỗi xảy ra, vui lòng thử lại"))
	}
	if err := c.Validate(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponse("validate", fmt.Sprintf("%s", err)))
	}
	if errAdd := u.UserAdd(); errAdd != nil {
		return c.JSON(http.StatusBadRequest, types.PayloadResponse("add_user", fmt.Sprintf("%s", errAdd)))
	}
	return c.JSON(http.StatusOK, u)
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
