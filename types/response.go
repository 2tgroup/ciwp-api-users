package types

type clientResponse struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

//PayloadResponseError is response error message
func PayloadResponseError(code, message string) clientResponse {
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
