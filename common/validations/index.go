package validation

import validatorP "gopkg.in/go-playground/validator.v9"

// CustomValidator to validate request
type CustomValidator struct {
	ValidatorX *validatorP.Validate
}

//Validate run validate request
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.ValidatorX.Struct(i)
}
