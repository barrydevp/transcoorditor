package util

import (
	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/gofiber/fiber/v2"
)

type ApiResponse interface {
	Status() int
	Response() *fiber.Map
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

func SendError(c *fiber.Ctx, msg string, err error) error {
	var apiErr *exception.ApiError

	switch v := err.(type) {
	case *exception.ApiError:
		apiErr = v
	default:
		apiErr = exception.ApiInternalError(err)
	}

	if msg != "" {
		apiErr.Msg = msg
	}

	return SendApiResponse(c, apiErr)
}
