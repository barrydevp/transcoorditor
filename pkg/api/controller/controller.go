package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	srv *service.Service
}

func NewController(srv *service.Service) *Controller {

	return &Controller{
		srv: srv,
	}
}

func (ctrl *Controller) PublicRoutes(a *fiber.App) {
	// create route
	route := a.Group("/api/v1")

	// register routes
	route.Get("/session/:sessionId", ctrl.GetSessionById)
	route.Put("/session/:sessionId", ctrl.PutSessionById)
	route.Delete("/session/:sessionId", ctrl.DeleteSessionById)
	route.Post("/session", ctrl.StartSession)
	route.Post("/session/:sessionId/join", ctrl.JoinSession)
	route.Post("/session/:sessionId/partial-commit", ctrl.PartialCommit)
	route.Post("/session/:sessionId/commit", ctrl.CommitSession)
	route.Post("/session/:sessionId/abort", ctrl.AbortSession)

}