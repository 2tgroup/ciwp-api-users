package users

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"

	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

func init() {

}

func UserUpdateHandler(c echo.Context) error {
	u := new(UserBase)
	if err := c.Bind(u); err != nil {
		log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("data_invaild", "Có lỗi xảy ra, vui lòng thử lại"))
	}

	if u.Status != 0 {
		user := c.Get("user").(*jwt.Token)
		premission := user.Claims.(*AuthJwtClaims)
		if premission.UserType != "admin" {
			u.Status = 0
		}
	}

	err := u.UserUpdate(u.ID.Hex())
	if err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("action_invaild", fmt.Sprintf("%s", err)))
	}
	return c.JSON(http.StatusOK, types.PayloadResponseOk(nil, nil))
}
