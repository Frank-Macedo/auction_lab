package rest_err

import "lab_fullcyle-auction_go/cmd/internal/internal_error"

type RestErr struct {
	Message string   `json:"message"`
	Err     string   `json:"error"`
	Code    int      `json:"code"`
	Causes  []Causes `json:"causes,omitempty"`
}

type Causes struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *RestErr) Error() string {
	return e.Message
}

func NewBadRequestError(message string, causes ...Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "Bad request",
		Code:    404,
		Causes:  causes,
	}
}

func NewInternalServerError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    500,
		Causes:  causes,
	}
}

func NewnotfoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not_found",
		Code:    404,
		Causes:  nil,
	}
}

func ConvertError(internalError *internal_error.InternalError) *RestErr {
	switch internalError.Err {
	case "bad request":
		return NewBadRequestError(internalError.Message)
	case "not found":
		return NewnotfoundError(internalError.Message)
	default:
		return NewInternalServerError("internal server error", nil)
	}
}
