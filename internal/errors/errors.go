package errorz

import (
	"errors"
	"gorm.io/gorm"
)

const CommonModule = 1

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")

	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested data is not found")

	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")

	// ErrTimeRegisterInputInvalid will throw if time register request invalid
	ErrTimeRegisterInputInvalid = errors.New("time register invalid")

	// ErrIDInvalid will throw if id params is invalid
	ErrIDInvalid = errors.New("id is invalid")

	// ErrExceptionError will throw if any the Exception Error happen
	ErrExceptionError = errors.New("exception error")
)

var (
	UnknownErrorCode = NewErrorCode(CommonModule, 500)
	errorCodes       = map[error]int{
		// Common errors
		ErrInternalServerError:      NewErrorCode(CommonModule, 0),
		ErrNotFound:                 NewErrorCode(CommonModule, 1),
		gorm.ErrRecordNotFound:      NewErrorCode(CommonModule, 1),
		ErrBadParamInput:            NewErrorCode(CommonModule, 2),
		ErrIDInvalid:                NewErrorCode(CommonModule, 3),
		ErrTimeRegisterInputInvalid: NewErrorCode(CommonModule, 4),
		ErrExceptionError:           NewErrorCode(CommonModule, 5),
	}
)

func NewErrorCode(module, detail int) int {
	code := module*10000 + detail
	return code
}

type ErrResponse struct {
	Description string `json:"error,omitempty"`
	ErrorCode   int    `json:"error_code"`
}

func GetErrorCode(err error) int {
	code, ok := errorCodes[err]
	if !ok {
		if unwrapErr := errors.Unwrap(err); unwrapErr != nil {
			return GetErrorCode(unwrapErr)
		}
		return UnknownErrorCode
	}
	return code
}

func NewErrResponse(err error, messages ...string) *ErrResponse {
	message := err.Error()
	if len(messages) > 0 {
		message = messages[0]
	}
	return &ErrResponse{
		Description: message,
		ErrorCode:   GetErrorCode(err),
	}
}
