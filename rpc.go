package mggo

import (
	"encoding/json"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"reflect"
	"strings"
)

func runRPC(address string) {
	cal := new(RPCInvoke)
	server := rpc.NewServer()
	server.Register(cal)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", address)
	if e != nil {
		panic(e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			panic(err)
		} else {
			LogInfo(nil, "new connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
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
	defer func() {
		//handlerError(ViewData{}, recover())
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
	res := method.Call(nil)
	LogInfo(nil, "End call rpc method", "->", methods)
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
	object  string
	service string
	ctx     *BaseContext
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
	client, err := net.Dial("tcp", host.String())
	defer client.Close()
	if err != nil {
		return nil, err
	}
	c := jsonrpc.NewClient(client)
	q, _ := json.Marshal(*params)
	args := &RPCArgs{Method: r.object + "." + method, Params: q}
	reply := []byte{}
	err = c.Call("RPCInvoke.Invoke", args, &reply)
	rec := NewRecord()
	json.Unmarshal(reply, &rec)
	if err != nil {
		return nil, err
	}
	return rec.Get("Result"), nil
}

// NewRPC is new RPC
func NewRPC(ctx *BaseContext, object, service string) *RPC {
	r := &RPC{
		object:  object,
		service: service,
		ctx:     ctx,
	}
	return r
}
