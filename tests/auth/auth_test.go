package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/2tgroup/ciwp-api-users/common/helptests"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/auth"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

var (
	userJSON  = `{"name":"Jon Snow","email":"jon@labstack.com","password":"123456"}`
	userLogin = `{"email":"jon@labstack.com","password":"123456"}`
	UserToken = ""
)

type getResponse struct {
	Code    string            `json:"code,omitempty"`
	Message string            `json:"message,omitempty"`
	Data    auth.AuthResponse `json:"data,omitempty"`
	Meta    interface{}       `json:"meta,omitempty"`
}

func TestResgisterUser(t *testing.T) {

	ht := helperTesting.HelperTesting{}

	// Make request
	ht.HelperTestMakeRequest(echo.POST, "/auth/register", userJSON)

	// Because have vaidator in module
	ht.HelperTestAddValidator()

	// set content type server allow
	ht.HelperTestSetDataType(echo.MIMEApplicationJSON)

	//Make context in echo
	c := ht.HelperTestFakeContext()

	auth.UserRegisterHandler(c)

	// Echo
	assert.Equal(t, ht.Eco, c.Echo())

	// Request
	assert.NotNil(t, c.Request())

	// Response
	assert.NotNil(t, c.Response())

	assert.Equal(t, http.StatusOK, ht.Res.Code)

	fmt.Println("Body Register:", ht.Res.Body.String())

}

func TestLoginUser(t *testing.T) {

	ht := helperTesting.HelperTesting{}

	// Make request
	ht.HelperTestMakeRequest(echo.POST, "/auth/login", userLogin)

	// Because have vaidator in module
	ht.HelperTestAddValidator()

	// set content type server allow
	ht.HelperTestSetDataType(echo.MIMEApplicationJSON)

	//Make context in echo
	c := ht.HelperTestFakeContext()

	auth.UserLoginHandler(c)

	// Echo
	assert.Equal(t, ht.Eco, c.Echo())

	// Request
	assert.NotNil(t, c.Request())

	// Response
	assert.NotNil(t, c.Response())

	assert.Equal(t, http.StatusOK, ht.Res.Code)

	jsonMap := new(getResponse)

	json.NewDecoder(ht.Res.Body).Decode(&jsonMap)

	UserToken = jsonMap.Data.Token

}

func TestGeTokenUser(t *testing.T) {

	// Parse the token
	token, err := jwt.ParseWithClaims(UserToken, &types.AuthJwtClaims{}, func(token *jwt.Token) (verifyKey interface{}, err error) {
		return verifyKey, err
	})

	fmt.Println(err)

	claims := token.Claims.(*types.AuthJwtClaims)

	assert.Equal(t, claims.Email, "jon@labstack.com", "Email should be equal")

}
