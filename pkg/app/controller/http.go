package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/cluster"
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/gofiber/fiber/v2"
)

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

func (ctrl *Controller) GetSessionByIdHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.srv.GetSessionById(sessionId, true)
	if err != nil {
		return util.SendError(c, "unable to get session", err)
	}

	return util.SendOK(c, session)
}

func (ctrl *Controller) ListSessionHttp(c *fiber.Ctx) error {
	session, err := ctrl.srv.ListSession()

	if err != nil {
		return util.SendError(c, "unable to list session", err)
	}

	return util.SendOK(c, session)
}

func (ctrl *Controller) PutSessionByIdHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session := &schema.Session{}

	if err := c.BodyParser(session); err != nil {
		return util.SendError(c, "unable to parse put session request payload", err)
	}

	session.Id = sessionId
	session, err := ctrl.srv.PutSessionById(session)
	if err != nil {
		return util.SendError(c, "unable to put session", err)
	}

	return util.SendOK(c, session)
}

func (ctrl *Controller) DeleteSessionByIdHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.srv.DeleteSessionById(sessionId)
	if err != nil {
		return util.SendError(c, "unable to delete session", err)
	}

	return util.SendOK(c, session)
}

func (ctrl *Controller) StartSessionHttp(c *fiber.Ctx) error {
	sessionOpts := schema.NewSessionOption()
	if err := c.BodyParser(sessionOpts); err != nil {
		return util.SendError(c, "unable to parse start session request payload", err)
	}

	session := schema.NewSession(sessionOpts)
	if _, err := ctrl.srv.StartSession(session); err != nil {
		return util.SendError(c, "unable to start new session", err)
	}

	// schedule the cleanup of session when timeout
	ctrl.recl.Schedule(&TimeoutSessionEntry{
		TimedoutAt: session.TimedoutAt(),
		SessionId:  session.Id,
	})

	return util.SendOK(c, session)
}

func (ctrl *Controller) JoinSessionHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	partJoinBody := &schema.ParticipantJoinBody{}
	if err := c.BodyParser(partJoinBody); err != nil {
		return util.SendError(c, "unable to parse join session request payload", err)
	}
	if err := common.GetValidate().Struct(partJoinBody); err != nil {
		return util.SendError(c, "invalid join session request payload", exception.AppBadRequest(err))
	}

	part := schema.NewParticipant()
	part.SessionId = sessionId
	part.ClientId = partJoinBody.ClientId
	part.RequestId = partJoinBody.RequestId

	part, err := ctrl.srv.JoinSession(sessionId, part)
	if err != nil {
		return util.SendError(c, "unable to join session", err)
	}

	return util.SendOK(c, part)
}

func (ctrl *Controller) PartialCommitHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	partCommit := &schema.ParticipantCommit{}
	if err := c.BodyParser(partCommit); err != nil {
		return util.SendError(c, "unable to parse partial commit session request payload", err)
	}
	if err := common.GetValidate().Struct(partCommit); err != nil {
		return util.SendError(c, "invalid partial commit session request payload", err)
	}

	part, err := ctrl.srv.PartialCommitSession(sessionId, partCommit)
	if err != nil {
		return util.SendError(c, "unable to partial commit session", err)
	}

	return util.SendOK(c, part)
}

func (ctrl *Controller) CommitSessionHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.srv.CommitSession(sessionId)
	if err != nil {
		return util.SendError(c, "unable to commit session", err)
	}

	return util.SendOK(c, session)
}

func (ctrl *Controller) AbortSessionHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.srv.AbortSession(sessionId)
	if err != nil {
		return util.SendError(c, "unable to abort session", err)
	}

	return util.SendOK(c, session)
}
