package mggo

import (
	"github.com/gorilla/websocket"
)

var sockets = map[int][]*websocket.Conn{}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SocketData struct {
	EventName string       `json:"event_name"`
	Msg       MapStringAny `json:"msg"`
}

// подключаемся к сокету. нужен пользователь
func socketConnect(ctx *BaseContext, userID int) {
	if userID == 0 {
		return
	}

	conn, _ := upgrader.Upgrade(ctx.Response, ctx.Request, nil)
	conns := sockets[userID]
	conns = append(conns, conn)
	sockets[userID] = conns
	go func() {
		s := &SocketData{
			EventName: "Socket.Connect",
			Msg: MapStringAny{
				"addr": conn.RemoteAddr().String(),
			},
		}
		SendSocket(s, conn)
	}()
}

// SendSocketUser is will send the user a message through the socket
func SendSocketUser(s *SocketData, userID int) {
	if conns, ok := sockets[userID]; ok {
		for _, conn := range conns {
			SendSocket(s, conn)
		}
	}
}

func SendSocket(s *SocketData, conn *websocket.Conn) {
	conn.WriteJSON(s)
}

func sendSockets(eventName string, users []int, msg MapStringAny) {
	s := &SocketData{
		EventName: eventName,
		Msg:       msg,
	}
	if len(users) == 0 {
		for _, conns := range sockets {
			for _, conn := range conns {
				SendSocket(s, conn)
			}
		}
	} else {
		for _, user := range users {
			SendSocketUser(s, user)
		}
	}
}
