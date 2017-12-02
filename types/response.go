package types

type clientResponse struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

//PayloadResponseError is response error message
func PayloadResponseError(code, message string) clientResponse {
	return PayloadResponseMgs(code, message)
}

//PayloadResponseMgs is response error message
func PayloadResponseMgs(code, message string) clientResponse {
	return clientResponse{
		Code:    code,
		Message: message,
	}
}

// PayloadResponseOk struct data response success
func PayloadResponseOk(data, meta interface{}) clientResponse {
	return clientResponse{
		Code:    "Ok",
		Message: "Successfully",
		Data:    data,
		Meta:    meta,
	}
}

const (
	//DataInvaild mess wrong struct data
	DataInvaild = "data_invaild"
	// ReqInvaild wrong struct data after bind
	ReqInvaild = "request_invaild"
	// ActionInvaild can not do some action
	ActionInvaild = "action_invaild"
	//NotValidate can not validate
	NotValidate = "not_validated"
	//ActionNotfound is not found somethings
	ActionNotfound = "not_found"
	//DataExist is exists something
	DataExist = "data_exists"
	//DataNotFound is exists something
	DataNotFound = "data_not_found"
	//ActionSuceess is done some action
	ActionSuceess = "action_success"
	//ActionError is action make error
	ActionError = "action_error"
)
