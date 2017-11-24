package helperTesting

import (
	"net/http"
	"net/http/httptest"
	"strings"

	validation "bitbucket.org/2tgroup/ciwp-api-users/common/validations"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

func init() {

}

//HelperTesting make testing easy
type HelperTesting struct {
	Res *httptest.ResponseRecorder
	Req *http.Request
	Eco *echo.Echo
}

//HelperTestMakeRequest make request simple
func (ht *HelperTesting) HelperTestMakeRequest(method, target string, body string) {
	ht.Eco = echo.New()
	ht.Req = httptest.NewRequest(echo.POST, target, strings.NewReader(body))
	ht.Res = httptest.NewRecorder()
}

//HelperTestSetDataType set content type send to server
func (ht *HelperTesting) HelperTestSetDataType(contentType string) {
	ht.Req.Header.Set("Content-Type", contentType)
}

//HelperTestSetHeader set custom header
func (ht *HelperTesting) HelperTestSetHeader(cusHeader map[string]string) {
	for field, val := range cusHeader {
		ht.Req.Header.Set(field, val)
	}
}

//HelperTestFakeContext fake context for echo
func (ht *HelperTesting) HelperTestFakeContext() echo.Context {
	return ht.Eco.NewContext(ht.Req, ht.Res)
}

//HelperTestAddValidator add my validator struct
func (ht *HelperTesting) HelperTestAddValidator() {
	ht.Eco.Validator = &validation.CustomValidator{ValidatorX: validator.New()}
}
