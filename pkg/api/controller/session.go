package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) GetSessionById(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	session, err := ctrl.ac.GetSessionById(sessionId, true)

	if err != nil {
		return util.Send500(c, err)
	}

	if session == nil {
		return util.Send404(c, util.Errorf("session not found"))
	}

	return util.SendSuccess(c, session)
}

func (ctrl *Controller) StartSession(c *fiber.Ctx) error {
	sessionOpts := schema.NewSessionOption()

	if err := c.BodyParser(sessionOpts); err != nil {
		return util.Send500(c, err)
	}

	session := schema.NewSession(sessionOpts)

	if _, err := ctrl.ac.StartSession(session); err != nil {
		return util.Send500(c, err)
	}

	return util.SendSuccess(c, session)
}

func (ctrl *Controller) JoinSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")
	partJoinBody := &schema.ParticipantJoinBody{}

	if err := c.BodyParser(partJoinBody); err != nil {
		return util.Send500(c, err)
	}

	if err := common.GetValidate().Struct(partJoinBody); err != nil {
		return util.Send400(c, err)
	}

	part := schema.NewParticipant()
	part.SessionId = sessionId
	part.ClientId = partJoinBody.ClientId
	part.RequestId = partJoinBody.RequestId

	part, err := ctrl.ac.JoinSession(sessionId, part)
	if err != nil {
		return util.Send500(c, err)
	}

	return util.SendSuccess(c, part)
}

func (ctrl *Controller) PartialCommit(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")
	partCommit := &schema.ParticipantCommit{}

	if err := c.BodyParser(partCommit); err != nil {
		return util.Send500(c, err)
	}

	if err := common.GetValidate().Struct(partCommit); err != nil {
		return util.Send400(c, err)
	}

	part, err := ctrl.ac.PartialCommitSession(sessionId, partCommit)
	if err != nil {
		return util.Send500(c, err)
	}

	return util.SendSuccess(c, part)
}

func (ctrl *Controller) CommitSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	part, err := ctrl.ac.CommitSession(sessionId)
	if err != nil {
		return util.Send500WithData(c, util.Errorf("cannot commit session, %w", err), part)
	}

	return util.SendSuccess(c, part)
}

func (ctrl *Controller) AbortSession(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

	part, err := ctrl.ac.AbortSession(sessionId)
	if err != nil {
		return util.Send500WithData(c, util.Errorf("cannot abort session, %w", err), part)
	}

	return util.SendSuccess(c, part)
}

func (ctrl *Controller) PutSessionById(c *fiber.Ctx) error {
	sessionId := c.Params("sessionId")

    session := &schema.Session{}

	if err := c.BodyParser(session); err != nil {
		return util.Send500(c, err)
	}

    session.Id = sessionId
	session, err := ctrl.ac.PutSessionById(session)

	if err != nil {
		return util.Send500(c, err)
	}

	if session == nil {
		return util.Send404(c, util.Errorf("session not found"))
	}

	return util.SendSuccess(c, session)
}

