package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type HTTPRequest struct {
	Method     string
	Path       string
	Version    string
	Headers    map[string]string
	Body       string
	Context    context.Context
	UserID     *int64
	RemoteAddr string
	Query      map[string]string
}

type HTTPResponse struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       string
}

type HandlerFunc func(ctx *Context)
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

type Server struct {
	port       string
	middleware []MiddlewareFunc
	router     *mux.Router
	server     *http.Server
	shutdown   chan struct{}
	wg         sync.WaitGroup
	isShutdown bool
	mu         sync.RWMutex
	logger     *logrus.Logger
}

func NewServer(port string) *Server {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return &Server{
		port:     port,
		shutdown: make(chan struct{}),
		router:   mux.NewRouter(),
		logger:   logger,
	}
}

func (s *Server) Use(middleware MiddlewareFunc) {
	s.middleware = append(s.middleware, middleware)
}

func (s *Server) AddRoute(method, path string, handler HandlerFunc) {
	// Apply middleware
	finalHandler := handler
	for i := len(s.middleware) - 1; i >= 0; i-- {
		finalHandler = s.middleware[i](finalHandler)
	}

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{
			Writer:  w,
			Request: r,
			Params:  mux.Vars(r),
			Query:   map[string]string{},
		}
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				ctx.Query[key] = values[0]
			}
		}
		finalHandler(ctx)
	}

	s.router.HandleFunc(path, httpHandler).Methods(method)
}

func (s *Server) GET(path string, handler HandlerFunc) {
	s.AddRoute("GET", path, handler)
}

func (s *Server) POST(path string, handler HandlerFunc) {
	s.AddRoute("POST", path, handler)
}

func (s *Server) PUT(path string, handler HandlerFunc) {
	s.AddRoute("PUT", path, handler)
}

func (s *Server) DELETE(path string, handler HandlerFunc) {
	s.AddRoute("DELETE", path, handler)
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Infof("Starting server on port %s", s.port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop() {
	s.mu.Lock()
	s.isShutdown = true
	s.mu.Unlock()

	s.logger.Info("Shutting down server...")

	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.Errorf("Error during server shutdown: %v", err)
		}
	}

	close(s.shutdown)
	s.wg.Wait()
	s.logger.Info("Server stopped")
}

// Helper function to create JSON responses
func JSONResponse(statusCode int, data interface{}) *HTTPResponse {
	body, _ := json.Marshal(data)
	return &HTTPResponse{
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(body),
	}
}

// Helper function to create error responses
func ErrorResponse(statusCode int, message string) *HTTPResponse {
	return JSONResponse(statusCode, map[string]string{"error": message})
}
