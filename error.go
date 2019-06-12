package mggo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
)

// ErrorMethodNotFound - error view not found
type ErrorMethodNotFound struct{}

// ErrorStatusForbidden - error not rights
type ErrorStatusForbidden struct{}

// ErrorViewNotFound - error method not found
type ErrorViewNotFound struct{}

// ErrorAuthenticate - error method not found
type ErrorAuthenticate struct{}

// ErrorInternalServer - error internal server
type ErrorInternalServer struct {
	Message string
}

func (e *ErrorMethodNotFound) Error() string {
	return fmt.Sprintf("method not found")
}
func (e *ErrorInternalServer) Error() string {
	return fmt.Sprintf(e.Message)
}
func (e *ErrorAuthenticate) Error() string {
	return fmt.Sprintf("not authenticate")
}
func (e *ErrorStatusForbidden) Error() string {
	return fmt.Sprintf("not allowed")
}

type errorMethod struct {
	Error interface{}
}

func handlerError(ctx *BaseContext, temp ViewData, r interface{}) {
	if r != nil {
		var printLog bool
		var message = r
		switch e := r.(type) {
		case ErrorViewNotFound:
			ctx.Response.WriteHeader(http.StatusNotFound)
			t, _ := template.ParseFiles(temp.DirView+temp.Template, temp.DirView+"404.html")
			t.Execute(ctx.Response, temp.Data)
			return
		case ErrorMethodNotFound:
			ctx.Response.WriteHeader(http.StatusNotFound)
			message = e.Error()
		case ErrorStatusForbidden:
			ctx.Response.WriteHeader(http.StatusForbidden)
			message = e.Error()
		case ErrorAuthenticate:
			ctx.Response.WriteHeader(http.StatusUnauthorized)
			message = e.Error()
		case ErrorInternalServer:
			ctx.Response.WriteHeader(http.StatusInternalServerError)
			message = e.Error()
			printLog = true
		default:
			ctx.Response.WriteHeader(http.StatusInternalServerError)
			printLog = true
		}

		buf2 := make([]byte, 4096)
		buf2 = buf2[:runtime.Stack(buf2, false)]
		err := errorMethod{Error: message}

		json.NewEncoder(ctx.Response).Encode(err)
		if printLog {
			LogError(ctx, r)
		}
	}
}
