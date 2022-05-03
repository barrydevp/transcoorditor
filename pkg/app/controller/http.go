package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/exception"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/gofiber/fiber/v2"
)

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
	var err error
	if err = c.BodyParser(sessionOpts); err != nil {
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

func (ctrl *Controller) ForgetSessionHttp(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.srv.ForgetSession(sessionId)
	if err != nil {
		return util.SendError(c, "unable to forget session", err)
	}

	return util.SendOK(c, session)
}
