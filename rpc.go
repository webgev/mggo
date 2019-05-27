package mggo

import (
    "net"
    "net/http"
    "net/rpc"
    "reflect"
    "strings"
)

func rpcServe(address string) {
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
    Params Record
}

// RPCInvoke is rpc
type RPCInvoke int

// Invoke is invoke method
func (r *RPCInvoke) Invoke(args *RPCArgs, reply *Record) error {
    LogInfo(args)
    methods := strings.Split(args.Method, ".")
    contr := getController(methods[0])
    MapToStruct(args.Params.ToMap(), &contr)
    contrValue := reflect.ValueOf(contr)

    method := contrValue.MethodByName(methods[1])
    if !method.IsValid() {
        panic(ErrorMethodNotFound{})
    }
    res := method.Call(nil)

    if len(res) > 0 {
        *reply = *StructToRecord(res[0].Interface())
    }
    return nil
}

// RPC struct
type RPC struct {
    client *rpc.Client
    object string
}

// Invoke rpc method
func (r *RPC) Invoke(method string, params *Record) (*Record, error) {
    args := &RPCArgs{Method: r.object + "." + method, Params: *params}
    reply := NewRecord()
    err := r.client.Call("RPCInvoke.Invoke", args, reply)
    if err != nil {
        LogError(err)
    }
    return reply, err
}

// NewRPC is new RPC
func NewRPC(object, service string) *RPC {
    serverConfig, err := config.GetSection("user")
    if err != nil {
        panic(err)
    }
    host, err := serverConfig.GetKey("address")
    if err != nil {
        panic(err)
    }
    client, _ := rpc.DialHTTP("tcp", host.String())
    r := &RPC{
        client: client,
        object: object,
    }
    return r
}
