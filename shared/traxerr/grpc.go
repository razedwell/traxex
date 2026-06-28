package traxerr

import (
	"google.golang.org/grpc/codes"
)

func GRPCCode(code Code) codes.Code {
	switch code {
	case CodeNotFound:
		return codes.NotFound
	case CodeUnauthorized:
		return codes.Unauthenticated
	case CodeInvalidInput:
		return codes.InvalidArgument
	case CodeConflict:
		return codes.AlreadyExists
	case CodeInternal:
		return codes.Internal
	case CodeUnavailable:
		return codes.Unavailable
	}
	return codes.Unknown
}
