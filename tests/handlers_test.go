package tests

import (
	"net/http/httptest"
	"strings"
	"testing"

	"server/handlers"
	"server/server"
)

func TestAuthHandler_Register(t *testing.T) {
	handler := handlers.NewAuthHandler(nil, "test-secret")

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "Valid registration",
			requestBody:    `{"email":"test@example.com","password":"password123","name":"Test User"}`,
			expectedStatus: 201,
		},
		{
			name:           "Missing email",
			requestBody:    `{"password":"password123","name":"Test User"}`,
			expectedStatus: 400,
		},
		{
			name:           "Missing password",
			requestBody:    `{"email":"test@example.com","name":"Test User"}`,
			expectedStatus: 400,
		},
		{
			name:           "Short password",
			requestBody:    `{"email":"test@example.com","password":"123","name":"Test User"}`,
			expectedStatus: 400,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"email":"test@example.com","password":"password123","name":"Test User"`,
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(tt.requestBody))
			recorder := httptest.NewRecorder()
			ctx := &server.Context{
				Writer:  recorder,
				Request: req,
				Params:  map[string]string{},
				Query:   map[string]string{},
			}
			handler.Register(ctx)
			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	handler := handlers.NewAuthHandler(nil, "test-secret")

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "Valid login",
			requestBody:    `{"email":"test@example.com","password":"password123"}`,
			expectedStatus: 200,
		},
		{
			name:           "Missing email",
			requestBody:    `{"password":"password123"}`,
			expectedStatus: 400,
		},
		{
			name:           "Missing password",
			requestBody:    `{"email":"test@example.com"}`,
			expectedStatus: 400,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"email":"test@example.com","password":"password123"`,
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(tt.requestBody))
			recorder := httptest.NewRecorder()
			ctx := &server.Context{
				Writer:  recorder,
				Request: req,
				Params:  map[string]string{},
				Query:   map[string]string{},
			}
			handler.Login(ctx)
			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	handler := handlers.NewUserHandler(nil)
	userID := int64(1)

	tests := []struct {
		name           string
		userID         *int64
		expectedStatus int
	}{
		{
			name:           "Authenticated user",
			userID:         &userID,
			expectedStatus: 200,
		},
		{
			name:           "Unauthenticated user",
			userID:         nil,
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/auth/me", nil)
			recorder := httptest.NewRecorder()
			ctx := &server.Context{
				Writer:  recorder,
				Request: req,
				Params:  map[string]string{},
				Query:   map[string]string{},
				UserID:  tt.userID,
			}
			handler.GetProfile(ctx)
			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	handler := handlers.NewUserHandler(nil)
	userID := int64(1)

	tests := []struct {
		name           string
		userID         *int64
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "Valid update",
			userID:         &userID,
			requestBody:    `{"name":"Updated Name"}`,
			expectedStatus: 200,
		},
		{
			name:           "Unauthenticated user",
			userID:         nil,
			requestBody:    `{"name":"Updated Name"}`,
			expectedStatus: 401,
		},
		{
			name:           "Missing name",
			userID:         &userID,
			requestBody:    `{"name":""}`,
			expectedStatus: 400,
		},
		{
			name:           "Invalid JSON",
			userID:         &userID,
			requestBody:    `{"name":"Updated Name"`,
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/auth/me", strings.NewReader(tt.requestBody))
			recorder := httptest.NewRecorder()
			ctx := &server.Context{
				Writer:  recorder,
				Request: req,
				Params:  map[string]string{},
				Query:   map[string]string{},
				UserID:  tt.userID,
			}
			handler.UpdateProfile(ctx)
			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	handler := handlers.NewUserHandler(nil)
	userID := int64(1)

	tests := []struct {
		name           string
		userID         *int64
		query          map[string]string
		expectedStatus int
	}{
		{
			name:           "Authenticated user",
			userID:         &userID,
			query:          map[string]string{},
			expectedStatus: 200,
		},
		{
			name:           "Unauthenticated user",
			userID:         nil,
			query:          map[string]string{},
			expectedStatus: 401,
		},
		{
			name:           "With pagination",
			userID:         &userID,
			query:          map[string]string{"limit": "5", "offset": "0"},
			expectedStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/users", nil)
			q := req.URL.Query()
			for k, v := range tt.query {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			recorder := httptest.NewRecorder()
			ctx := &server.Context{
				Writer:  recorder,
				Request: req,
				Params:  map[string]string{},
				Query:   tt.query,
				UserID:  tt.userID,
			}
			handler.ListUsers(ctx)
			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}
