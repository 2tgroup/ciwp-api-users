package users

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"bitbucket.org/2tgroup/ciwp-api-users/common/findcountry"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

func init() {
	UserLoadCountry()
}

//UserUpdateHandler update user
func UserUpdateHandler(c echo.Context) error {

	u := new(UserBase)

	if err := c.Bind(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.DataInvaild, "error invaild request, please check your data"))
	}

	user := c.Get("user").(*jwt.Token)

	premission := user.Claims.(*types.AuthJwtClaims)

	if u.Status != 0 {
		if premission.UserType != "admin" {
			u.Status = 0
		}
	}

	if premission.ID.Hex() != u.ID.Hex() && premission.UserType != "admin" {
		return c.JSON(http.StatusForbidden, types.PayloadResponseError(types.ActionInvaild, "Not owner profile"))
	}

	if u.UserCheckEmailExits(u.ID.Hex()) == true {
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.DataExist, "Email exist"))
	}

	if err := u.UserUpdate(u.ID.Hex()); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.ActionInvaild, fmt.Sprintf("%s", err)))
	}

	return c.JSON(http.StatusOK, types.PayloadResponseOk(nil, nil))
}

func UserLoadCountry() {

	data := findcountry.Country.MapByName("South Korea")
	fmt.Println(data.Name)           // Will Print: South Korea
	fmt.Println(data.Alpha2)         // Will Print: KR
	fmt.Println(data.Alpha3)         // Will Print: KOR
	fmt.Println(data.Currency[0])    // Will Print: KRW
	fmt.Println(data.CallingCode[0]) // Will Print: 82
	fmt.Println(data.Region)         // Will Print: Asia
	fmt.Println(data.Subregion)      // Will Print: Eastern Asia

}
