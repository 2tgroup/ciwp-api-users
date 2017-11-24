package auth

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/2tgroup/ciwp-api-users/common/helptests"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/auth"
)

var (
	userJSON = `{"name":"Jon Snow","email":"jon@labstack.com","password":"77"}`
)

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

	fmt.Println("Body Res:", ht.Res.Body.String())

	//assert.Equal(t, userJSON, rec.Body.String())

}
