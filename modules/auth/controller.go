package auth

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/helpers/crypt"
	"bitbucket.org/2tgroup/ciwp-api-users/libaries/redisCache"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/users"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

//UserTokenHandler is
func UserTokenHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*types.AuthJwtClaims)
	t, err := getJWToken(claims)
	if err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ReqInvaild, fmt.Sprintf("%s", err)))
	}
	return c.JSON(http.StatusOK, types.PayloadResponseOk(echo.Map{
		"token": t,
		"user":  user.Claims,
	}, nil))
}

func UserLoginHandler(c echo.Context) error {

	u := new(users.UserBase)

	if err := c.Bind(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusUnauthorized, types.PayloadResponseError(types.ReqInvaild, "error invaild request, please check your data"))
	}
	if err := c.Validate(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusUnauthorized, types.PayloadResponseError(types.NotValidate, "Wrong password or email"))
	}

	err := u.UserGetOne(echo.Map{
		"email": u.Email,
	})

	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.PayloadResponseError(types.ActionNotfound, "Wrong password or email"))
	}

	if checked := u.UserCheckPass(); checked != true {
		return c.JSON(http.StatusUnauthorized, types.PayloadResponseError(types.ActionNotfound, "Wrong password or email"))
	}

	t, err := u.AuthSignupToken()

	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.PayloadResponseError(types.ActionNotfound, "Wrong password or email"))
	}

	uRes := new(AuthResponse)

	uRes.Token = t

	uRes.AuthSetResponse(*u)

	return c.JSON(http.StatusOK, types.PayloadResponseOk(uRes, nil))
}

func UserRegisterHandler(c echo.Context) error {

	u := new(users.UserBase)

	if err := c.Bind(u); err != nil {
		log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ReqInvaild, "error invaild request, please check your data"))
		//return c.JSON(http.StatusBadRequest, types.PayloadResponseError("request_invaild", "error invaild request, please check your data"))
	}
	if err := c.Validate(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.NotValidate, fmt.Sprintf("%s", err)))
	}

	/*Checking User exist or not*/

	u.UserGetOne(echo.Map{
		"email": u.Email,
	})

	if u.ID.Hex() != "" {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.DataExist, "Email exist"))
	}

	if errAdd := u.UserAdd(); errAdd != nil {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ActionError, fmt.Sprintf("%s", errAdd)))
	}

	t, err := u.AuthSignupToken()

	if err != nil {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ActionError, fmt.Sprintf("%s", err)))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ActionNotfound, fmt.Sprintf("%s", err)))
	}

	uRes := new(AuthResponse)

	uRes.Token = t

	uRes.AuthSetResponse(*u)

	return c.JSON(http.StatusOK, types.PayloadResponseOk(uRes, nil))

}

func UserForgotHandler(c echo.Context) error {

	u := new(users.UserBase)

	if err := c.Bind(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ReqInvaild, "error invaild request, please check your data"))
	}

	err := u.UserGetOne(echo.Map{
		"email": u.Email,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, types.PayloadResponseError(types.ActionNotfound, "Email not exist"))
	}

	u.UserGeneratePass()

	if err != nil {
		return c.JSON(http.StatusUnauthorized, types.PayloadResponseError(types.ActionNotfound, "Wrong password or email"))
	}

	return c.JSON(http.StatusOK, types.PayloadResponseMgs(types.ActionSuceess, "The new password has been sent to your email "+u.Email))
}

func UserLogoutHandler(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)

	blackJWT := helperCtypt.GenerateCrypt(c.Request().Header.Get("Authorization"))

	claims := user.Claims.(*types.AuthJwtClaims)

	/* claims.ExpiresAt = time.Now().Add(time.Hour * -72).Unix()

	claims.IssuedAt = claims.ExpiresAt

	t, err := getJWToken(claims) */

	if err := redisCache.Set("TokenKey", claims.ID.Hex(), []byte(blackJWT), 5); err != true {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ReqInvaild, ""))
	}

	return c.JSON(http.StatusOK, types.PayloadResponseOk(nil, nil))
}

func getJWToken(Authclaims *types.AuthJwtClaims) (t string, e error) {
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Authclaims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.DataConfig.SecretKey))
	if err != nil {
		return "", err
	}
	return t, err
}
