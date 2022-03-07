package util

import (
	"github.com/gofiber/fiber/v2"
)

func SendSuccess(c *fiber.Ctx, data interface{}) error {

	return c.JSON(fiber.Map{
		"ok":   true,
		"data": data,
	})
}

func SendError(c *fiber.Ctx, status int, err error) error {

	return c.Status(status).JSON(fiber.Map{
		"ok":  false,
		"msg": err.Error(),
	})
}

func Send400(c *fiber.Ctx, err error) error {
	return SendError(c, fiber.StatusBadRequest, err)
}

func Send404(c *fiber.Ctx, err error) error {
	return SendError(c, fiber.StatusNotFound, err)
}

func Send500(c *fiber.Ctx, err error) error {
	return SendError(c, fiber.StatusInternalServerError, err)
}

func SendErrorWithData(c *fiber.Ctx, status int, err error, data interface{}) error {
	response := fiber.Map{
		"ok":  false,
		"msg": err.Error(),
	}

	if data != nil {
		response["data"] = data
	}

	return c.Status(status).JSON(response)
}

func Send500WithData(c *fiber.Ctx, err error, data interface{}) error {
	return SendErrorWithData(c, fiber.StatusInternalServerError, err, data)
}
