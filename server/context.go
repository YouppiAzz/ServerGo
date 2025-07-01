package server

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string
	Query   map[string]string
	Status  int
	UserID  *int64
}

// JSON sends a JSON response
func (c *Context) JSON(status int, data interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(status)
	json.NewEncoder(c.Writer).Encode(data)
}

// String sends a plain text response
func (c *Context) String(status int, text string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(status)
	c.Writer.Write([]byte(text))
}

// Param gets a path parameter
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// QueryParam gets a query parameter
func (c *Context) QueryParam(key string) string {
	return c.Request.URL.Query().Get(key)
}

// BindJSON parses the request body as JSON into obj
func (c *Context) BindJSON(obj interface{}) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, obj)
}

// SetStatus sets the response status code
func (c *Context) SetStatus(status int) {
	c.Status = status
}

// Header sets a header
func (c *Context) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}
