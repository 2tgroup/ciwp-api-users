package types

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

//PayloadResponseError is response error message
func PayloadResponseError(code, message string) errorResponse {
	return errorResponse{
		Code:    code,
		Message: message,
	}
}

// PayloadResponseOk struct data response success
func PayloadResponseOk(data, meta interface{}) successResponse {
	return successResponse{
		Code:    "Ok",
		Message: "Successfully",
		Data:    data,
		Meta:    meta,
	}
}
