package exception

import (
	"errors"
	"fmt"

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
	err        error
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
	}

	if r.err != nil {
		resp["err"] = r.err.Error()
	}

	if r.Detail != nil {
		resp["detail"] = r.Detail
	}

	return &resp
}

func (r *ApiError) Err() error {
	return r.err
}

func (r *ApiError) Error() string {
	if r.err != nil {
		return r.err.Error()
	}

	return r.Msg
}

func NewApiError(err error, msg string, status int, detail interface{}) *ApiError {
	return &ApiError{
		err:        err,
		Msg:        msg,
		StatusCode: status,
		Detail:     detail,
	}
}

func ApiErrorf(status int, msg string, format string, a ...interface{}) *ApiError {
	return NewApiError(fmt.Errorf(format, a...), msg, status, nil)
}

func ApiUnauthorized(err error) *ApiError {
	return NewApiError(err, "un authorized", fiber.StatusUnauthorized, nil)
}

func ApiNotFound(err error) *ApiError {
	return NewApiError(err, "not found", fiber.StatusNotFound, nil)
}

func ApiNotFoundf(format string, a ...interface{}) *ApiError {
	return ApiNotFound(fmt.Errorf(format, a...))
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

func ApiPreconditionFailed(err error) *ApiError {
	return NewApiError(err, "precondition failed", fiber.StatusPreconditionFailed, nil)
}

func ApiPreconditionFailedf(format string, a ...interface{}) *ApiError {
	return ApiPreconditionFailed(fmt.Errorf(format, a...))
}

func ApiUnprocessableEntity(err error) *ApiError {
	return NewApiError(err, "unprosessable entity", fiber.StatusUnprocessableEntity, nil)
}

func ApiUnprocessableEntityf(format string, a ...interface{}) *ApiError {
	return ApiUnprocessableEntity(fmt.Errorf(format, a...))
}
