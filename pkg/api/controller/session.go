package controller

import (
	"github.com/barrydevp/transcoorditor/pkg/action"
	// "github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) GetSessionById(c *fiber.Ctx) error {
	sesisonId := c.Params("sessionId")

	session, err := ctrl.ac.GetSessionById(sesisonId, false)

	if err != nil {
		return util.Send500(c, err)
	}

	if session == nil {
		return util.Send404(c, util.NewError("session not found"))
	}

	return util.SendSuccess(c, session)
}

func (ctrl *Controller) StartSession(c *fiber.Ctx) error {
	sessionOpts := action.NewSessionOption()

	if err := c.BodyParser(sessionOpts); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, err)
	}

	session := action.NewSession(sessionOpts)

	if err := ctrl.ac.StartSession(session); err != nil {
		return util.Send500(c, err)
	}

	return util.SendSuccess(c, session)
}

// func (ctrl *Controller) JoinSession(c *fiber.Ctx) error {
// 	sessionId := c.Params("sessionId")
//
// 	if err := c.BodyParser(sessionOpts); err != nil {
// 		return util.SendError(c, fiber.StatusBadRequest, err)
// 	}
//
// 	session := action.NewSession(sessionOpts)
//
// 	if err := ctrl.ac.StartSession(session); err != nil {
// 		return util.Send500(c, err)
// 	}
//
// 	return util.SendSuccess(c, fiber.Map{
// 		"session": session,
// 	})
// }
