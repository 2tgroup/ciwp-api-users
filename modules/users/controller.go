package users

import (
	"fmt"
	"net/http"

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
	err := u.UserUpdate(u.ID.Hex())
	if err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError("action_invaild", fmt.Sprintf("%s", err)))
	}
	return c.JSON(http.StatusOK, types.PayloadResponseOk(echo.Map{}, nil))
}
