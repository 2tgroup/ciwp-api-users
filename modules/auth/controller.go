package auth

import (
	"fmt"
	"net/http"

	"bitbucket.org/2tgroup/ciwp-api-users/types"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/users"
)

//UserTokenHandler is
func UserTokenHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*users.AuthJwtClaims)
	t, err := getJWToken(claims)
	if err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("token_invaild", fmt.Sprintf("%s", err)))
	}
	return c.JSON(http.StatusOK, types.PayloadResponseOk(echo.Map{
		"token": t,
	}, nil))
}

func UserLoginHandler(c echo.Context) error {

	u := new(users.UserBase)

	if err := c.Bind(u); err != nil {
		log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("request_invaild", "Có lỗi xảy ra, vui lòng thử lại"))
	}
	if err := c.Validate(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("validate", fmt.Sprintf("%s", err)))
	}

	err := u.UserGetOne(echo.Map{
		"email": u.Email,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, types.PayloadResponseError("not_found", fmt.Sprintf("%s", err)))
	}

	if checked := u.UserCheckPass(); checked != true {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("not_found", "Wrong password or email"))
	}

	t, _ := u.AuthSignupToken()

	return c.JSON(http.StatusOK, types.PayloadResponseOk(echo.Map{
		"token": t,
		"user": echo.Map{
			"_id":       u.ID.Hex(),
			"email":     u.Email,
			"name":      u.Name,
			"user_type": u.UserType,
			"status":    u.Status,
			"info":      u.UserInfo,
		},
	}, nil))
}

func UserRegisterHandler(c echo.Context) error {

	u := new(users.UserBase)

	if err := c.Bind(u); err != nil {
		log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("request_invaild", "Có lỗi xảy ra, vui lòng thử lại"))
	}
	if err := c.Validate(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("validate", fmt.Sprintf("%s", err)))
	}

	/*Checking User exist or not*/

	u.UserGetOne(echo.Map{
		"email": u.Email,
	})

	if u.ID.Hex() != "" {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("user_exist", "Email exist"))
	}

	if errAdd := u.UserAdd(); errAdd != nil {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("add_user", fmt.Sprintf("%s", errAdd)))
	}

	t, err := u.AuthSignupToken()

	if err != nil {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("add_user", fmt.Sprintf("%s", err)))
	}

	return c.JSON(http.StatusOK, types.PayloadResponseOk(echo.Map{
		"token": t,
		"user": echo.Map{
			"_id":       u.ID.Hex(),
			"email":     u.Email,
			"name":      u.Name,
			"user_type": u.UserType,
			"status":    u.Status,
			"info":      u.UserInfo,
		},
	}, nil))
}

func getJWToken(Authclaims *users.AuthJwtClaims) (t string, e error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Authclaims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.DataConfig.SecretKey))
	if err != nil {
		return "", err
	}
	return t, err
}
