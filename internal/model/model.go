package model

type GenericResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GenericError(code int, message string) GenericResponse {
	return GenericResponse{
		Code:    code,
		Message: message,
	}
}
