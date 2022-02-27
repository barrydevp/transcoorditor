package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/action"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	ac *action.Action
}

func NewController(ac *action.Action) *Controller {

	return &Controller{
		ac: ac,
	}
}

func (ctrl *Controller) PublicRoutes(a *fiber.App) {
	// create route
	route := a.Group("/api/v1")

	// register routes
	route.Get("/session/:sessionId", ctrl.GetSessionById)
	route.Post("/session", ctrl.StartSession)
	// route.Post("/session/:sesionId/join")
	// route.Post("/session/:sesionId/update")
	// route.Post("/session/:sesionId/end")

}