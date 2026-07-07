package traxerr

type Code string

const (
	CodeNotFound        Code = "NOT_FOUND"
	CodeSessionNotFound Code = "SESSION_NOT_FOUND"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeInvalidInput    Code = "INVALID_INPUT"
	CodeConflict        Code = "CONFLICT"
	CodeInternal        Code = "INTERNAL"
	CodeUnavailable     Code = "UNAVAILABLE"
)
