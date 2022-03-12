package exception

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound           = errors.New("NOT_FOUND")
	ErrInvalidArgument    = errors.New("INVALID_ARGUMENT")
	ErrPreconditionFailed = errors.New("PRECONDITION_FAILED")
	ErrAborted            = errors.New("ABORTED")
	ErrUnauthorized       = errors.New("UNAUTHORIZED")
	ErrForbidden          = errors.New("FORBIDDEN")
	ErrGone               = errors.New("GONE")
)

type ApiError struct {
	Err        error
	Msg        string
	StatusCode int
	Detail     interface{}
}

func (r *ApiError) Status() int {
	if r.StatusCode < 400 {
		return fiber.StatusInternalServerError
	}

	return r.StatusCode
}

func (r *ApiError) Response() *fiber.Map {
	resp := fiber.Map{
		"ok":  false,
		"msg": r.Msg,
		"err": r.Err.Error(),
	}

	if r.Detail != nil {
		resp["detail"] = r.Detail
	}

	return &resp
}

func (r *ApiError) Error() string {
	return r.Err.Error()
}

func NewApiError(err error, msg string, status int, detail interface{}) *ApiError {
	return &ApiError{
		Err:        err,
		Msg:        msg,
		StatusCode: status,
		Detail:     detail,
	}
}

func ApiUnauthorized(err error) *ApiError {
	return NewApiError(err, "un authorized", fiber.StatusUnauthorized, nil)
}

func ApiNotFound(err error) *ApiError {
	return NewApiError(err, "not found", fiber.StatusNotFound, nil)
}

func ApiForbidden(err error) *ApiError {
	return NewApiError(err, "forbidden", fiber.StatusForbidden, nil)
}

func ApiConflict(err error) *ApiError {
	return NewApiError(err, "conflict", fiber.StatusConflict, nil)
}

func ApiGone(err error) *ApiError {
	return NewApiError(err, "gone", fiber.StatusGone, nil)
}

func ApiBadRequest(err error) *ApiError {
	return NewApiError(err, "bad request", fiber.StatusBadRequest, nil)
}

func ApiServiceUnavailable(err error) *ApiError {
	return NewApiError(err, "bad request", fiber.StatusServiceUnavailable, nil)
}

func ApiInternalError(err error) *ApiError {
	return NewApiError(err, "internal error", fiber.StatusInternalServerError, nil)
}

func ApiTimeout(err error) *ApiError {
	return NewApiError(err, "timeout", fiber.StatusGatewayTimeout, nil)
}
