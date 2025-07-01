package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"server/auth"
	"server/server"
)

type responseWriterWithStatus struct {
	http.ResponseWriter
	status int
}

func (w *responseWriterWithStatus) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// CORS middleware
func CORS() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if ctx.Request.Method == "OPTIONS" {
				ctx.Writer.WriteHeader(200)
				return
			}
			next(ctx)
		}
	}
}

// Logger middleware
func Logger() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) {
			rw := &responseWriterWithStatus{ResponseWriter: ctx.Writer, status: 200}
			ctx.Writer = rw
			start := time.Now()
			next(ctx)
			duration := time.Since(start)
			log.Printf("%s %s %d %v %s",
				ctx.Request.Method,
				ctx.Request.URL.Path,
				rw.status,
				duration,
				ctx.Request.RemoteAddr)
		}
	}
}

// Security middleware
func Security() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) {
			ctx.Header("X-Content-Type-Options", "nosniff")
			ctx.Header("X-Frame-Options", "DENY")
			ctx.Header("X-XSS-Protection", "1; mode=block")
			ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			next(ctx)
		}
	}
}

// Rate limiter
type rateLimiter struct {
	clients map[string][]time.Time
	mu      sync.RWMutex
	limit   int
}

func RateLimiter(requestsPerMinute int) server.MiddlewareFunc {
	rl := &rateLimiter{
		clients: make(map[string][]time.Time),
		limit:   requestsPerMinute,
	}

	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) {
			clientIP := strings.Split(ctx.Request.RemoteAddr, ":")[0]

			rl.mu.Lock()
			now := time.Now()

			// Clean old entries
			if requests, exists := rl.clients[clientIP]; exists {
				var validRequests []time.Time
				for _, t := range requests {
					if now.Sub(t) < time.Minute {
						validRequests = append(validRequests, t)
					}
				}
				rl.clients[clientIP] = validRequests
			}

			// Check rate limit
			if len(rl.clients[clientIP]) >= rl.limit {
				rl.mu.Unlock()
				ctx.Writer.WriteHeader(429)
				ctx.JSON(429, map[string]string{"error": "Rate limit exceeded"})
				return
			}

			// Add current request
			rl.clients[clientIP] = append(rl.clients[clientIP], now)
			rl.mu.Unlock()

			next(ctx)
		}
	}
}

// Auth middleware
func RequireAuth(jwtSecret string) server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) {
			authHeader := ctx.Request.Header.Get("Authorization")
			if authHeader == "" {
				ctx.Writer.WriteHeader(401)
				ctx.JSON(401, map[string]string{"error": "Authorization header required"})
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				ctx.Writer.WriteHeader(401)
				ctx.JSON(401, map[string]string{"error": "Bearer token required"})
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := auth.ValidateToken(token, jwtSecret)
			if err != nil {
				ctx.Writer.WriteHeader(401)
				ctx.JSON(401, map[string]string{"error": "Invalid token"})
				return
			}

			ctx.UserID = &userID
			next(ctx)
		}
	}
}

// Request ID middleware
func RequestID() server.MiddlewareFunc {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx *server.Context) {
			requestID := ctx.Request.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}
			ctx.Header("X-Request-ID", requestID)
			next(ctx)
		}
	}
}

func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
