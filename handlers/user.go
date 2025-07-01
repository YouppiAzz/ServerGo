package handlers

import (
	"database/sql"
	"strconv"

	"server/models"
	"server/server"
)

type UserHandler struct {
	userRepo *models.UserRepository
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		userRepo: models.NewUserRepository(db),
	}
}

func (h *UserHandler) GetProfile(ctx *server.Context) {
	if ctx.UserID == nil {
		ctx.JSON(401, map[string]string{"error": "User not authenticated"})
		return
	}

	user, err := h.userRepo.GetByID(*ctx.UserID)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Database error"})
		return
	}

	if user == nil {
		ctx.JSON(404, map[string]string{"error": "User not found"})
		return
	}

	ctx.JSON(200, user)
}

type UpdateProfileRequest struct {
	Name string `json:"name"`
}

func (h *UserHandler) UpdateProfile(ctx *server.Context) {
	if ctx.UserID == nil {
		ctx.JSON(401, map[string]string{"error": "User not authenticated"})
		return
	}

	var updateReq UpdateProfileRequest
	if err := ctx.BindJSON(&updateReq); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
		return
	}

	if updateReq.Name == "" {
		ctx.JSON(400, map[string]string{"error": "Name is required"})
		return
	}

	user, err := h.userRepo.GetByID(*ctx.UserID)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Database error"})
		return
	}

	if user == nil {
		ctx.JSON(404, map[string]string{"error": "User not found"})
		return
	}

	user.Name = updateReq.Name
	if err := h.userRepo.Update(user); err != nil {
		ctx.JSON(500, map[string]string{"error": "Failed to update user"})
		return
	}

	ctx.JSON(200, user)
}

func (h *UserHandler) ListUsers(ctx *server.Context) {
	if ctx.UserID == nil {
		ctx.JSON(401, map[string]string{"error": "User not authenticated"})
		return
	}

	limit := 10
	offset := 0

	if l := ctx.QueryParam("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if o := ctx.QueryParam("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	users, err := h.userRepo.List(limit, offset)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Database error"})
		return
	}
	total, err := h.userRepo.Count()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Database error"})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"users":  users,
		"limit":  limit,
		"offset": offset,
		"total":  total,
	})
}
