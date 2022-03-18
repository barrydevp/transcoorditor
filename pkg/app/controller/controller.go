package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/controlplane"
	"github.com/barrydevp/transcoorditor/pkg/controlplane/reconciler"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var logger = common.Logger().WithFields(logrus.Fields{
	"pkg": "app/ctrl",
})

type Controller struct {
	srv  *service.Service
	recl *reconciler.ScheduleReconciler
	l    *logrus.Entry
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
	route.Get("/sessions", ctrl.ListSessionHttp)
	route.Get("/sessions/:sessionId", ctrl.GetSessionByIdHttp)
	route.Put("/sessions/:sessionId", ctrl.PutSessionByIdHttp)
	route.Delete("/sessions/:sessionId", ctrl.DeleteSessionByIdHttp)
	route.Post("/sessions", ctrl.StartSessionHttp)
	route.Post("/sessions/:sessionId/join", ctrl.JoinSessionHttp)
	route.Post("/sessions/:sessionId/partial-commit", ctrl.PartialCommitHttp)
	route.Post("/sessions/:sessionId/commit", ctrl.CommitSessionHttp)
	route.Post("/sessions/:sessionId/abort", ctrl.AbortSessionHttp)

}

func (ctrl *Controller) RegisterReconciler(c *controlplane.ControlPlane) {
	// init reconciler
	recl := reconciler.NewScheduleReconciler(ctrl.InitTimeoutSessionQueueRecl, ctrl.HandleTimeoutSessionRecl)
	ctrl.recl = recl
	c.RegisterRecl(recl)
}
