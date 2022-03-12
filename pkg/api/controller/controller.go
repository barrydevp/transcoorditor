package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	srv *service.Service
	l   *logrus.Entry
}

func NewController(srv *service.Service) *Controller {

	return &Controller{
		srv: srv,
		l: common.Logger().WithFields(logrus.Fields{
			"pkg": "ctrl",
		}),
	}
}

func (ctrl *Controller) PublicRoutes(a *fiber.App) {
	// create route
	route := a.Group("/api/v1")

	// register routes
	route.Get("/sessions", ctrl.ListSession)
	route.Get("/sessions/:sessionId", ctrl.GetSessionById)
	route.Put("/sessions/:sessionId", ctrl.PutSessionById)
	route.Delete("/sessions/:sessionId", ctrl.DeleteSessionById)
	route.Post("/sessions", ctrl.StartSession)
	route.Post("/sessions/:sessionId/join", ctrl.JoinSession)
	route.Post("/sessions/:sessionId/partial-commit", ctrl.PartialCommit)
	route.Post("/sessions/:sessionId/commit", ctrl.CommitSession)
	route.Post("/sessions/:sessionId/abort", ctrl.AbortSession)

}
