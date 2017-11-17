package types

type errorResponse struct {
	ErrorCode string `json:"err_code"`
	Message   string `json:"message"`
}

//PayloadResponse is response error message
func PayloadResponse(code, message string) errorResponse {
	return errorResponse{
		ErrorCode: code,
		Message:   message,
	}
}
