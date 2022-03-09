package controller

import (
	"github.com/gofiber/fiber/v2"
)

type ApiResponse interface {
	Status() int
	Response() *fiber.Map
}

type ApiError struct {
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
		"msg": r.Err.Error(),
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

func SendResponse(c *fiber.Ctx, resp interface{}) error {
	var apiResp ApiResponse

	switch response := resp.(type) {
	case ApiResponse:
	case error:
		return SendApiResponse(c, &ApiError{
			Err:        response,
			StatusCode: fiber.StatusInternalServerError,
		})
	}

	return SendApiResponse(c, apiResp)
}

func SendInternalError(c *fiber.Ctx, err error) error {
	return SendApiResponse(c, &ApiError{
		Err:        err,
		StatusCode: fiber.StatusInternalServerError,
	})
}
