package mggo

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
)

// BaseContext is call context
type BaseContext struct {
	Response    http.ResponseWriter
	Request     *http.Request
	Path        []string
	Query       url.Values
	CurrentUser User
	uuid        string
}

func newBaseContext(response http.ResponseWriter, request *http.Request, path []string, query url.Values, user User) *BaseContext {
	return &BaseContext{response, request, path, query, user, generateUIIDMethod()}
}

func generateUIIDMethod() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
