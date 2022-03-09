package controller

import (
	"errors"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/service"
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) GetSessionById(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.srv.GetSessionById(sessionId, true)
	if errors.Is(err, service.ErrNotFound) {
		return SendError(c, fiber.StatusNotFound, "session not found", err)
	} else if err != nil {
		return SendInternalError(c, "unable to get session", err)
	}

	return SendOK(c, session)
}

func (ctrl *Controller) PutSessionById(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session := &schema.Session{}

	if err := c.BodyParser(session); err != nil {
		return SendInternalError(c, "unable to parse put session request payload", err)
	}

	session.Id = sessionId
	session, err := ctrl.srv.PutSessionById(session)
	if errors.Is(err, service.ErrNotFound) {
		return SendError(c, fiber.StatusNotFound, "session not found", err)
	} else if err != nil {
		return SendInternalError(c, "unable to put session", err)
	}

	return SendOK(c, session)
}

func (ctrl *Controller) DeleteSessionById(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.srv.DeleteSessionById(sessionId)
	if errors.Is(err, service.ErrNotFound) {
		return SendError(c, fiber.StatusNotFound, "session not found", err)
	} else if err != nil {
		return SendInternalError(c, "unable to delete session", err)
	}

	return SendOK(c, session)
}

func (ctrl *Controller) StartSession(c *fiber.Ctx) error {
	sessionOpts := schema.NewSessionOption()
	if err := c.BodyParser(sessionOpts); err != nil {
		return SendInternalError(c, "unable to parse start session request payload", err)
	}

	session := schema.NewSession(sessionOpts)
	if _, err := ctrl.srv.StartSession(session); err != nil {
		return SendInternalError(c, "unable to start new session", err)
	}

	return SendOK(c, session)
}

func (ctrl *Controller) JoinSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	partJoinBody := &schema.ParticipantJoinBody{}
	if err := c.BodyParser(partJoinBody); err != nil {
		return SendInternalError(c, "unable to parse join session request payload", err)
	}
	if err := common.GetValidate().Struct(partJoinBody); err != nil {
		return SendError(c, fiber.StatusBadRequest, "invalid join session request payload", err)
	}

	part := schema.NewParticipant()
	part.SessionId = sessionId
	part.ClientId = partJoinBody.ClientId
	part.RequestId = partJoinBody.RequestId

	part, err := ctrl.srv.JoinSession(sessionId, part)
	if errors.Is(err, service.ErrNotFound) {
		return SendError(c, fiber.StatusNotFound, "session not found", err)
	} else if errors.Is(err, service.ErrPreconditionFailed) {
		return SendError(c, fiber.StatusPreconditionFailed, "unable to join session", err)
	} else if err != nil {
		return SendInternalError(c, "unable to join session", err)
	}

	return SendOK(c, part)
}

func (ctrl *Controller) PartialCommit(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	partCommit := &schema.ParticipantCommit{}
	if err := c.BodyParser(partCommit); err != nil {
		return SendInternalError(c, "unable to parse partial commit session request payload", err)
	}
	if err := common.GetValidate().Struct(partCommit); err != nil {
		return SendOK(c, err)
	}

	part, err := ctrl.srv.PartialCommitSession(sessionId, partCommit)
	if errors.Is(err, service.ErrNotFound) {
		return SendError(c, fiber.StatusNotFound, "session not found", err)
	} else if err != nil {
		return SendInternalError(c, "unable to delete session", err)
	}

	return SendOK(c, part)
}

func (ctrl *Controller) CommitSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	part, err := ctrl.srv.CommitSession(sessionId)
	if errors.Is(err, service.ErrPreconditionFailed) {
		return SendError(c, fiber.StatusPreconditionFailed, "unable to commit session", err)
	} else if err != nil {
		return SendInternalError(c, "unable to commit session", err)
	}

	return SendOK(c, part)
}

func (ctrl *Controller) AbortSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	part, err := ctrl.srv.AbortSession(sessionId)
	if errors.Is(err, service.ErrPreconditionFailed) {
		return SendError(c, fiber.StatusPreconditionFailed, "unable to abort session", err)
	} else if err != nil {
		return SendInternalError(c, "unable to abort session", err)
	}

	return SendOK(c, part)
}
