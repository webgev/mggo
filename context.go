package mggo

import (
	"net/http"
)

type BaseContext struct {
	Response http.ResponseWriter
	Request  *http.Request
	Path     []string
}
