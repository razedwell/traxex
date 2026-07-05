package traxerr

import "net/http"

func HTTPStatus(code Code) int {
	switch code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeInvalidInput:
		return http.StatusBadRequest
	case CodeConflict:
		return http.StatusConflict
	case CodeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
