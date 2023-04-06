package errs

type CustomError struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	ErrServer = NewError(500, "Internal server error.")

	ErrParam        = NewError(400, "Invalid parameters.")
	ErrUnAuthorized = NewError(401, "UnAuthorized.")
)

func (e *CustomError) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *CustomError {
	return &CustomError{
		Msg:  msg,
		Code: code,
	}
}

func GetError(e *CustomError, data interface{}) *CustomError {
	return &CustomError{
		Msg:  e.Msg,
		Code: e.Code,
		Data: data,
	}
}
