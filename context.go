package mggo

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

// BaseContext is call context
type BaseContext struct {
	Response    http.ResponseWriter
	Request     *http.Request
	Path        []string
	CurrentUser User
	uuid        string
}

func newBaseContext(response http.ResponseWriter, request *http.Request, path []string, user User) *BaseContext {
	return &BaseContext{response, request, path, user, generateUIIDMethod()}
}

func generateUIIDMethod() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
