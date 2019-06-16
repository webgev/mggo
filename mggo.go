package mggo

import (
	"math/rand"
	"net/http"
	"os"

	"github.com/go-ini/ini"
)

var (
	config   *ini.File
	initFlag bool
)

type callbackHandler func()

var callbacks []callbackHandler

// startServer is start server
func startServer(w http.ResponseWriter, req *http.Request) {
	SQLOpen()
}

// endServer is stop server and handler error
func endServer(ctx *BaseContext, temp ViewData) {
	handlerError(ctx, temp, recover())
}

// InitCallback registers a callback function that will be called when the configuration is initialized.
func InitCallback(handler callbackHandler) {
	if callbacks == nil {
		callbacks = []callbackHandler{handler}
	} else {
		callbacks = append(callbacks, handler)
	}
}

//Run http service
func Run(rout Router, pathConfing string) {
	if initFlag {
		panic("init")
	}

	cfg, err := ini.Load(pathConfing)
	if err != nil {
		os.Exit(1)
	}

	initFlag = true
	config = cfg

	SQLOpen()
	for _, handler := range callbacks {
		handler()
	}
	rout.run()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
