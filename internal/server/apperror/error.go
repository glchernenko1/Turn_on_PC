package apperror

import "encoding/json"

var (
	ErrNotFound  = NewAppError(nil, "not found", "", "404")
	BadRequest   = NewAppError(nil, "bad request", "", "400")
	Unauthorized = NewAppError(nil, "unauthorized", "", "401")
	NameTaken    = NewAppError(nil, "This name is already taken", "", "400")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func SystemError(err error) *AppError {
	return &AppError{
		Err:              err,
		Message:          "system error",
		DeveloperMessage: err.Error(),
		Code:             "US-000",
	}
}
