package common

type OaiError struct {
	Msg    string
	Code   int
	Detail any
}

func (e *OaiError) Error() string {
	return e.Msg
}

func NewOaiError(msg string, code int, detail any) *OaiError {
	return &OaiError{
		Msg:    msg,
		Code:   code,
		Detail: detail,
	}
}

var (
	NoImplementError   = NewOaiError("not implement", 404, nil)
	InternalError      = NewOaiError("internal error", 500, nil)
	CreateRequestError = NewOaiError("create request error", 500, nil)
	StatusCodeError    = NewOaiError("status code error", 500, nil)
)
