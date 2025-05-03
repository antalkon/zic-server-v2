package handlers

import (
	"backend/internal/transport/rest/req"
	"backend/internal/transport/rest/res"
	"backend/internal/transport/service"
	"backend/internal/utils"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ActionsHandler struct {
	actions  *service.ActionsService
	validate *validator.Validate
}

func NewActionsHandler(actions *service.ActionsService) *ActionsHandler {
	return &ActionsHandler{
		actions:  actions,
		validate: validator.New(),
	}
}

func (h *ActionsHandler) SendReboot(c echo.Context) error {
	var r req.SendRebootReq

	if err := c.Bind(&r); err != nil {
		code, resp := utils.BadRequestError()
		return c.JSON(code, resp)
	}

	if err := h.validate.Struct(r); err != nil {
		code, resp := utils.ValidationError()
		resp.Message = err.Error()
		return c.JSON(code, resp)
	}

	if err := h.actions.SendReboot(r.ComputerID, r.Delay); err != nil {
		code, resp := utils.InternalServerError("Failed to send reboot command: " + err.Error())
		return c.JSON(code, resp)
	}

	return c.JSON(200, res.MessageRes{
		Message: "Reboot command sent successfully",
	})
}

func (h *ActionsHandler) SendShutdown(c echo.Context) error {
	var r req.SendShutdownReq
	if err := c.Bind(&r); err != nil {
		code, resp := utils.BadRequestError()
		return c.JSON(code, resp)
	}
	if err := h.validate.Struct(r); err != nil {
		code, resp := utils.ValidationError()
		resp.Message = err.Error()
		return c.JSON(code, resp)
	}
	if err := h.actions.SendShutdown(r.ComputerID, r.Delay); err != nil {
		code, resp := utils.InternalServerError("Failed to send shutdown command: " + err.Error())
		return c.JSON(code, resp)
	}
	return c.JSON(200, res.MessageRes{
		Message: "Shutdown command sent successfully",
	})
}

func (h *ActionsHandler) SendBlock(c echo.Context) error {
	var r req.SendBlockReq
	if err := c.Bind(&r); err != nil {
		code, resp := utils.BadRequestError()
		return c.JSON(code, resp)
	}
	if err := h.validate.Struct(r); err != nil {
		code, resp := utils.ValidationError()
		resp.Message = err.Error()
		return c.JSON(code, resp)
	}
	if err := h.actions.SendBlock(r.ComputerID); err != nil {
		code, resp := utils.InternalServerError("Failed to send block command: " + err.Error())
		return c.JSON(code, resp)
	}
	return c.JSON(200, res.MessageRes{
		Message: "block command sent successfully",
	})
}

func (h *ActionsHandler) SendUnblock(c echo.Context) error {
	var r req.SendUnblockReq
	if err := c.Bind(&r); err != nil {
		code, resp := utils.BadRequestError()
		return c.JSON(code, resp)
	}
	if err := h.validate.Struct(r); err != nil {
		code, resp := utils.ValidationError()
		resp.Message = err.Error()
		return c.JSON(code, resp)
	}
	if err := h.actions.SendUnblock(r.ComputerID); err != nil {
		code, resp := utils.InternalServerError("Failed to send unblock command: " + err.Error())
		return c.JSON(code, resp)
	}
	return c.JSON(200, res.MessageRes{
		Message: "unblock command sent successfully",
	})
}

func (h *ActionsHandler) SendLockScreen(c echo.Context) error {
	var r req.SendLockScreenReq
	if err := c.Bind(&r); err != nil {
		code, resp := utils.BadRequestError()
		return c.JSON(code, resp)
	}
	if err := h.validate.Struct(r); err != nil {
		code, resp := utils.ValidationError()
		resp.Message = err.Error()
		return c.JSON(code, resp)
	}
	if err := h.actions.SendLockScreen(r.ComputerID); err != nil {
		code, resp := utils.InternalServerError("Failed to send lock screen command: " + err.Error())
		return c.JSON(code, resp)
	}
	return c.JSON(200, res.MessageRes{
		Message: "lock screen command sent successfully",
	})
}

func (h *ActionsHandler) SendUrl(c echo.Context) error {
	var r req.SendUrlReq
	if err := c.Bind(&r); err != nil {
		code, resp := utils.BadRequestError()
		return c.JSON(code, resp)
	}
	if err := h.validate.Struct(r); err != nil {
		code, resp := utils.ValidationError()
		resp.Message = err.Error()
		return c.JSON(code, resp)
	}
	if err := h.actions.SendUrl(r.ComputerID, r.Url); err != nil {
		code, resp := utils.InternalServerError("Failed to send URL command: " + err.Error())
		return c.JSON(code, resp)
	}
	return c.JSON(200, res.MessageRes{
		Message: "URL command sent successfully",
	})
}

func (h *ActionsHandler) SendMessage(c echo.Context) error {
	var r req.SendMessageReq
	if err := c.Bind(&r); err != nil {
		code, resp := utils.BadRequestError()
		return c.JSON(code, resp)
	}
	if err := h.validate.Struct(r); err != nil {
		code, resp := utils.ValidationError()
		resp.Message = err.Error()
		return c.JSON(code, resp)
	}
	if err := h.actions.SendMessage(r.ComputerID, r.Message, r.Type); err != nil {
		code, resp := utils.InternalServerError("Failed to send message command: " + err.Error())
		return c.JSON(code, resp)
	}
	return c.JSON(200, res.MessageRes{
		Message: "message command sent successfully",
	})
}
