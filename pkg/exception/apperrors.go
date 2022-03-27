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

type AppError struct {
	err        error
	Msg        string
	StatusCode int
	Detail     interface{}
}

func (r *AppError) Status() int {
	if r.StatusCode < 400 {
		return fiber.StatusInternalServerError
	}

	return r.StatusCode
}

func (r *AppError) Response() *fiber.Map {
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

func (r *AppError) Err() error {
	return r.err
}

func (r *AppError) Error() string {
	if r.err != nil {
		return r.err.Error()
	}

	return r.Msg
}

func NewAppError(err error, msg string, status int, detail interface{}) *AppError {
	return &AppError{
		err:        err,
		Msg:        msg,
		StatusCode: status,
		Detail:     detail,
	}
}

func AppErrorf(status int, msg string, format string, a ...interface{}) *AppError {
	return NewAppError(fmt.Errorf(format, a...), msg, status, nil)
}

func AppUnauthorized(err error) *AppError {
	return NewAppError(err, "un authorized", fiber.StatusUnauthorized, nil)
}

func AppNotFound(err error) *AppError {
	return NewAppError(err, "not found", fiber.StatusNotFound, nil)
}

func AppNotFoundf(format string, a ...interface{}) *AppError {
	return AppNotFound(fmt.Errorf(format, a...))
}

func AppForbidden(err error) *AppError {
	return NewAppError(err, "forbidden", fiber.StatusForbidden, nil)
}

func AppConflict(err error) *AppError {
	return NewAppError(err, "conflict", fiber.StatusConflict, nil)
}

func AppGone(err error) *AppError {
	return NewAppError(err, "gone", fiber.StatusGone, nil)
}

func AppGonef(format string, a ...interface{}) *AppError {
	return AppGone(fmt.Errorf(format, a...))
}

func AppBadRequest(err error) *AppError {
	return NewAppError(err, "bad request", fiber.StatusBadRequest, nil)
}

func AppServiceUnavailable(err error) *AppError {
	return NewAppError(err, "bad request", fiber.StatusServiceUnavailable, nil)
}

func AppInternalError(err error) *AppError {
	return NewAppError(err, "internal error", fiber.StatusInternalServerError, nil)
}

func AppTimeout(err error) *AppError {
	return NewAppError(err, "timeout", fiber.StatusGatewayTimeout, nil)
}

func AppPreconditionFailed(err error) *AppError {
	return NewAppError(err, "precondition failed", fiber.StatusPreconditionFailed, nil)
}

func AppPreconditionFailedf(format string, a ...interface{}) *AppError {
	return AppPreconditionFailed(fmt.Errorf(format, a...))
}

func AppUnprocessableEntity(err error) *AppError {
	return NewAppError(err, "unprosessable entity", fiber.StatusUnprocessableEntity, nil)
}

func AppUnprocessableEntityf(format string, a ...interface{}) *AppError {
	return AppUnprocessableEntity(fmt.Errorf(format, a...))
}
