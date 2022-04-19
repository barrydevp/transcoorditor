package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/cluster"
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
	c    *cluster.Cluster
	srv  *service.Service
	recl *reconciler.ScheduleReconciler
	l    *logrus.Entry
}

func NewController(c *cluster.Cluster, srv *service.Service) *Controller {
	return &Controller{
		c:   c,
		srv: srv,
		l: common.Logger().WithFields(logrus.Fields{
			"pkg": "ctrl",
		}),
	}
}

func (ctrl *Controller) SystemRoutes(a *fiber.App) {
	// create route
	route := a.Group("/api/sys")

	// heartbeat
	route.Get("/ping", ctrl.PingHttp)

	// cluster routes
	route.Post("/initiate", ctrl.InitiateClusterHttp)
	route.Post("/join", ctrl.JoinClusterHttp)
	route.Post("/left", ctrl.LeftClusterHttp)
	route.Get("/rsconf", ctrl.GetClusterRsConfHttp)
	route.Get("/stats", ctrl.GetClusterStatsHttp)
	route.Get("/leader", ctrl.GetClusterLeaderHttp)
	route.Get("/nconf", ctrl.GetClusterCurrentHttp)
}

func (ctrl *Controller) PublicRoutes(a *fiber.App) {
	// create route
	route := a.Group("/api/v1")

	// txn routes
	route.Get("/sessions", ctrl.ListSessionHttp)
	route.Get("/sessions/:sessionId", ctrl.GetSessionByIdHttp)
	route.Put("/sessions/:sessionId", ctrl.PutSessionByIdHttp)
	route.Post("/sessions", ctrl.StartSessionHttp)
	route.Post("/sessions/:sessionId/join", ctrl.JoinSessionHttp)
	route.Post("/sessions/:sessionId/partial-commit", ctrl.PartialCommitHttp)
	route.Post("/sessions/:sessionId/commit", ctrl.CommitSessionHttp)
	route.Post("/sessions/:sessionId/abort", ctrl.AbortSessionHttp)
	route.Post("/sessions/:sessionId/forget", ctrl.ForgetSessionHttp)

    // internal testing
	route.Delete("/sessions/:sessionId", ctrl.DeleteSessionByIdHttp)

}

func (ctrl *Controller) RegisterReconciler(c *controlplane.ControlPlane) {
	// init reconciler
	recl := reconciler.NewScheduleReconciler(ctrl.InitTimeoutSessionQueueRecl, ctrl.HandleTimeoutSessionRecl)
	ctrl.recl = recl
	c.RegisterRecl(recl)
}
