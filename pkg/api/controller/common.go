package controller

import (
	"github.com/gofiber/fiber/v2"
)

type ApiResponse interface {
	Status() int
	Response() *fiber.Map
}

type ApiError struct {
	Msg        string
	Err        error
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

func NewApiError(err error, status int) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: status,
	}
}

type ApiSuccess struct {
	Data       interface{}
	StatusCode int
}

func (r *ApiSuccess) Status() int {
	if r.StatusCode < 200 || r.StatusCode >= 400 {
		return fiber.StatusInternalServerError
	}

	return r.StatusCode
}

func (r *ApiSuccess) Response() *fiber.Map {
	resp := fiber.Map{
		"ok": true,
	}

	if r.Data != nil {
		resp["data"] = r.Data
	}

	return &resp
}

func SendApiResponse(c *fiber.Ctx, apiResp ApiResponse) error {
	return c.Status(apiResp.Status()).JSON(apiResp.Response())
}

func SendSuccess(c *fiber.Ctx, status int, data interface{}) error {
	return SendApiResponse(c, &ApiSuccess{
		Data:       data,
		StatusCode: status,
	})
}

func SendOK(c *fiber.Ctx, data interface{}) error {
	return SendSuccess(c, fiber.StatusOK, data)
}

func SendError(c *fiber.Ctx, status int, msg string, err error) error {
	return SendApiResponse(c, &ApiError{
		Msg:        msg,
		Err:        err,
		StatusCode: status,
	})
}

func SendInternalError(c *fiber.Ctx, msg string, err error) error {
	return SendError(c, fiber.StatusInternalServerError, msg, err)
}
