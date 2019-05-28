package mggo

import (
    "github.com/go-ini/ini"
    "net/http"
)

var (
    config   *ini.File
    respose  http.ResponseWriter
    request  http.Request
    initFlag bool
)

type callbackHandler func()

var callbacks []callbackHandler

// startServer is start server
func startServer(w http.ResponseWriter, req *http.Request) {
    respose = w
    request = *req
    SQLOpen()
}

// endServer is stop server and handler error
func endServer(temp ViewData) {
    SQLClose()
    handlerError(temp, recover())
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
func Run(rout Router, cfg *ini.File) {
    if initFlag {
        panic("init")
    }
    initFlag = true
    config = cfg
    
    SQLOpen()
    for _, handler := range callbacks {
        handler()
    }
    SQLClose()
    rout.run()
}
