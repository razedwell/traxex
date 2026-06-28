package traxerr

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
)

type Traxerr struct {
	Code    Code
	Message string
	Cause   error
	Details map[string]any
}

func (e *Traxerr) Error() string {
	cs, ds := "", ""
	if e.Cause != nil {
		cs = fmt.Sprintf("->%v", e.Cause)
	}
	if e.Details != nil {
		ds = fmt.Sprintf(" %+v", e.Details)
	}
	return fmt.Sprintf("%s: %s%s%s", e.Code, e.Message, cs, ds)
}

func (e *Traxerr) Unwrap() error {
	return e.Cause
}

func New(code Code, msg string) *Traxerr {
	return &Traxerr{
		Code:    code,
		Message: msg,
	}
}

func Wrap(code Code, msg string, cause error) *Traxerr {
	return &Traxerr{
		Code:    code,
		Message: msg,
		Cause:   cause,
	}
}

func IsCode(err error, code Code) bool {
	var traxerr *Traxerr
	if errors.As(err, &traxerr) {
		return traxerr.Code == code
	}
	return false
}

func ToGRPCCode(err error) codes.Code {
	var traxerr *Traxerr
	if errors.As(err, &traxerr) {
		return GRPCCode(traxerr.Code)
	}
	return codes.Unknown
}
