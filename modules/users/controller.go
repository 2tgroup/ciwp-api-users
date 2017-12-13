package users

import (
	"fmt"
	"net/http"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"bitbucket.org/2tgroup/ciwp-api-users/common/findcountry"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

//UserUpdateHandler update user
func UserUpdateHandler(c echo.Context) error {

	u := new(UserBase)

	if err := c.Bind(u); err != nil {
		//log.Errorf("Wrong request %s", err)
		return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.DataInvaild, "Error invaild request, please check your data"))
	}

	user := c.Get("user").(*jwt.Token)

	premission := user.Claims.(*types.AuthJwtClaims)

	if u.Status != 0 {
		if premission.UserType != "admin" {
			u.Status = 0
		}
	}

	if u.Email != "" {
		regEmail := regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w+)+$`)
		if !regEmail.MatchString(u.Email) {
			return c.JSON(http.StatusBadRequest, types.PayloadResponseError(types.DataInvaild, "You should input correct format email!"))
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

	u.UserGetOne(echo.Map{
		"_id": u.ID,
	})

	return c.JSON(http.StatusOK, types.PayloadResponseOk(u.UserResponse(), nil))
}

//UserListCountry return list countries
func UserListCountryByAlpha2(c echo.Context) error {
	return c.JSON(http.StatusOK, types.PayloadResponseOk(findcountry.Country.ListByAlpha2, nil))
}

func UserListCounty(c echo.Context) error {

	dataFound := findcountry.Country.MapByAlpha2(c.Param("alpha2"))

	if dataFound == nil {
		return c.JSON(http.StatusNotFound, types.PayloadResponseError(types.DataNotFound, "Country not found"))
	}

	return c.JSON(http.StatusOK, types.PayloadResponseOk(dataFound, nil))
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
