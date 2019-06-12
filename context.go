package mggo

import (
	"net/http"
)

// BaseContext is call context
type BaseContext struct {
	Response    http.ResponseWriter
	Request     *http.Request
	Path        []string
	CurrentUser User
}

func newBaseContext(response http.ResponseWriter, request *http.Request, path []string, user User) *BaseContext {
	return &BaseContext{response, request, path, user}
}
