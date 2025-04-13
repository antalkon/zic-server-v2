package handlers

import (
	"backend/internal/models"
	"backend/internal/transport/rest/req"
	"backend/internal/transport/service"
	"backend/internal/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	room     *service.RoomService
	validate *validator.Validate
}

func NewRoomHandler(room *service.RoomService) *RoomHandler {
	return &RoomHandler{
		room:     room,
		validate: validator.New(),
	}
}

func (h *RoomHandler) CreateRoom(c echo.Context) error {
	var room models.Room
	if err := c.Bind(&room); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "failed to bind request: "+err.Error())
		return c.JSON(code, msg)
	}

	if err := h.validate.Struct(room); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "validation failed: "+err.Error())
		return c.JSON(code, msg)
	}
	userRole := c.Get("user_role").(string)

	roomID, err := h.room.CreateRoom(&room, userRole)
	if err != nil {
		code, msg := utils.InternalServerError("failed to create room: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusCreated, roomID)
}

func (h *RoomHandler) GetAllRooms(c echo.Context) error {
	rooms, err := h.room.GetAllRooms()
	if err != nil {
		code, msg := utils.InternalServerError("failed to get rooms: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, rooms)
}

func (h *RoomHandler) GetRoomByID(c echo.Context) error {
	roomID := c.Param("id")
	if roomID == "" {
		code, msg := utils.NewError(http.StatusBadRequest, "room ID is required")
		return c.JSON(code, msg)
	}

	room, err := h.room.GetRoomByID(roomID)
	if err != nil {
		code, msg := utils.InternalServerError("failed to get room: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.JSON(http.StatusOK, room)
}

func (h *RoomHandler) UpdateRoom(c echo.Context) error {
	roomID := c.Param("id")
	if roomID == "" {
		code, msg := utils.NewError(http.StatusBadRequest, "room ID is required")
		return c.JSON(code, msg)
	}

	userRole := c.Get("user_role").(string)

	var roomReq req.UpdateRoomRequest
	if err := c.Bind(&roomReq); err != nil {
		code, msg := utils.NewError(http.StatusBadRequest, "failed to bind request: "+err.Error())
		return c.JSON(code, msg)
	}

	err := h.room.UpdateRoom(roomID, &roomReq, userRole)
	if err != nil {
		code, msg := utils.InternalServerError("failed to update room: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RoomHandler) DeleteRoom(c echo.Context) error {
	roomID := c.Param("id")
	if roomID == "" {
		code, msg := utils.NewError(http.StatusBadRequest, "room ID is required")
		return c.JSON(code, msg)
	}

	userRole := c.Get("user_role").(string)

	err := h.room.DeleteRoom(roomID, userRole)
	if err != nil {
		code, msg := utils.InternalServerError("failed to delete room: " + err.Error())
		return c.JSON(code, msg)
	}

	return c.NoContent(http.StatusNoContent)
}
