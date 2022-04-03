package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) PingHttp(c *fiber.Ctx) error {
	return util.SendOK(c, nil)
}

func (ctrl *Controller) InitiateClusterHttp(c *fiber.Ctx) error {
	rsconf := &cluster.ClusterRsConf{}

	if err := c.BodyParser(rsconf); err != nil {
		return util.SendError(c, "unable to parse rsconf", err)
	}

	if err := ctrl.c.RsInitiate(rsconf); err != nil {
		return util.SendError(c, "unable to initiate cluster replicaset", err)
	}

	confFuture := ctrl.c.Ra.GetConfiguration()
	if confFuture.Error() != nil {
		return util.SendError(c, "cannot get rsconf from cluster after initiate", confFuture.Error())
	}

	return util.SendOK(c, confFuture.Configuration())
}

func (ctrl *Controller) JoinClusterHttp(c *fiber.Ctx) error {
	node := &cluster.Node{}

	if err := c.BodyParser(node); err != nil {
		return util.SendError(c, "unable to parse node", err)
	}

	if err := ctrl.c.Join(node); err != nil {
		return util.SendError(c, "unable to join cluster replicaset", err)
	}

	conf, err := ctrl.c.GetConf()
	if err != nil {
		return util.SendError(c, "unable to get rsconf from cluster", err)
	}

	return util.SendOK(c, conf)
}

func (ctrl *Controller) LeftClusterHttp(c *fiber.Ctx) error {
	node := &cluster.Node{}

	if err := c.BodyParser(node); err != nil {
		return util.SendError(c, "unable to parse node", err)
	}

	if err := ctrl.c.Left(node); err != nil {
		return util.SendError(c, "unable to left cluster replicaset", err)
	}

	conf, err := ctrl.c.GetConf()
	if err != nil {
		return util.SendError(c, "unable to get rsconf from cluster", err)
	}

	return util.SendOK(c, conf)
}

func (ctrl *Controller) GetClusterRsConfHttp(c *fiber.Ctx) error {
	conf, err := ctrl.c.GetConf()
	if err != nil {
		return util.SendError(c, "unable get cluster's rsconf", err)
	}

	return util.SendOK(c, conf)
}

func (ctrl *Controller) GetClusterStatsHttp(c *fiber.Ctx) error {
	stats, err := ctrl.c.Stats()
	if err != nil {
		return util.SendError(c, "cannot get cluster's stats", err)
	}

	return util.SendOK(c, stats)
}

func (ctrl *Controller) GetClusterLeaderHttp(c *fiber.Ctx) error {
	stats, err := ctrl.c.Leader()
	if err != nil {
		return util.SendError(c, "cannot get leader node", err)
	}

	return util.SendOK(c, stats)
}

func (ctrl *Controller) GetClusterCurrentHttp(c *fiber.Ctx) error {
	stats, err := ctrl.c.Current()
	if err != nil {
		return util.SendError(c, "cannot get current node", err)
	}

	return util.SendOK(c, stats)
}
