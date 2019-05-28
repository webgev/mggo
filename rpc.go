package mggo

import (
    "net"
    "net/http"
    "net/rpc"
    "reflect"
    "strings"
    "encoding/json"
)

func runRpc(address string) {
    arith := new(RPCInvoke)
    rpc.Register(arith)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", address)
    if e != nil {
        panic(e)
    }
    go http.Serve(l, nil)
}

// RPCArgs is arguments for rpc invoke
type RPCArgs struct {
    Method string
    Params []byte
}

// RPCInvoke is rpc
type RPCInvoke int

// Invoke is invoke method
func (r *RPCInvoke) Invoke(args *RPCArgs, reply *[]byte) error {
    SQLOpen()
    defer func () {
        SQLClose()
        handlerError(ViewData{}, recover())
    }()
    methods := strings.Split(args.Method, ".")
    contr := getController(methods[0])
    rec := NewRecord()
    json.Unmarshal(args.Params, &rec)
    MapToStruct(rec.ToMap(), &contr)
    contrValue := reflect.ValueOf(contr)

    method := contrValue.MethodByName(methods[1])
    if !method.IsValid() {
        panic(ErrorMethodNotFound{})
    }
    LogInfo("Call rpc method", "->", methods)
    res := method.Call(nil)
    LogInfo("End call rpc method", "->", methods)
    if len(res) > 0 {
        r := NewRecord()
        r.Add("Result", res[0].Interface())
        q, _ := json.Marshal(r)
        *reply = q
    }
    return nil
}

// RPC struct
type RPC struct {
    object string
    service string
}

// Invoke rpc method
func (r *RPC) Invoke(method string, params *Record) (interface{}, error) {
    serverConfig, err := config.GetSection(r.service)
    if err != nil {
        return nil, err
    }
    host, err := serverConfig.GetKey("address")
    if err != nil {
        return nil, err
    }
    client, err := rpc.DialHTTP("tcp", host.String())
    if err != nil {
        return nil, err
    }

    q, _ := json.Marshal(*params)
    args := &RPCArgs{Method: r.object + "." + method, Params: q}
    reply := []byte{}
    err = client.Call("RPCInvoke.Invoke", args, &reply)
    defer client.Close()
    rec := NewRecord()
    json.Unmarshal(reply, &rec)
    if err != nil {
        return nil, err
    }
    return rec.Get("Result"), nil
}

// NewRPC is new RPC
func NewRPC(object, service string) *RPC {
    r := &RPC{
        object: object,
        service: service,
    }
    return r
}
