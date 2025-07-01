package handlers

import (
	"time"

	"server/server"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

func (h *HealthHandler) Health(ctx *server.Context) {
	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Database:  "disconnected", // For now, always disconnected
		Version:   "2.0.0",
		Uptime:    time.Since(startTime).String(),
	}
	ctx.JSON(200, health)
}

var startTime = time.Now()
