package handlers

import (
	"database/sql"
	"strings"

	"server/auth"
	"server/models"
	"server/server"
)

type AuthHandler struct {
	userRepo  *models.UserRepository
	jwtSecret string
}

func NewAuthHandler(db *sql.DB, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  models.NewUserRepository(db),
		jwtSecret: jwtSecret,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Register(ctx *server.Context) {
	var registerReq RegisterRequest
	if err := ctx.BindJSON(&registerReq); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
		return
	}

	if registerReq.Email == "" || registerReq.Password == "" || registerReq.Name == "" {
		ctx.JSON(400, map[string]string{"error": "Email, password, and name are required"})
		return
	}

	if len(registerReq.Password) < 6 {
		ctx.JSON(400, map[string]string{"error": "Password must be at least 6 characters"})
		return
	}

	existingUser, err := h.userRepo.GetByEmail(registerReq.Email)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Database error"})
		return
	}

	if existingUser != nil {
		ctx.JSON(409, map[string]string{"error": "User already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(registerReq.Password)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Failed to hash password"})
		return
	}

	user := &models.User{
		Email:        strings.ToLower(registerReq.Email),
		PasswordHash: hashedPassword,
		Name:         registerReq.Name,
	}

	if err := h.userRepo.Create(user); err != nil {
		ctx.JSON(500, map[string]string{"error": "Failed to create user"})
		return
	}

	token, err := auth.GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Failed to generate token"})
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	ctx.JSON(201, response)
}

func (h *AuthHandler) Login(ctx *server.Context) {
	var loginReq LoginRequest
	if err := ctx.BindJSON(&loginReq); err != nil {
		ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
		return
	}

	user, err := h.userRepo.GetByEmail(strings.ToLower(loginReq.Email))
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Database error"})
		return
	}

	if user == nil || !auth.CheckPasswordHash(loginReq.Password, user.PasswordHash) {
		ctx.JSON(401, map[string]string{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Failed to generate token"})
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	ctx.JSON(200, response)
}
