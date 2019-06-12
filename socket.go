package mggo

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var sockets map[int][]*websocket.Conn = map[int][]*websocket.Conn{}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type socketData struct {
	EventName string      `json:"event_name"`
	Msg       interface{} `json:"msg"`
}

func socketConnect(userID int, w http.ResponseWriter, r *http.Request) {
	if userID == 0 {
		return
	}

	conn, _ := upgrader.Upgrade(w, r, nil)
	conns := sockets[userID]
	conns = append(conns, conn)
	sockets[userID] = conns
}

// sendSocketUser is will send the user a message through the socket
func sendSocketUser(s *socketData, userID int, msg interface{}) {
	if conns, ok := sockets[userID]; ok {
		for _, conn := range conns {
			sendSocket(s, conn, msg)
		}
	}
}

func sendSocket(s *socketData, conn *websocket.Conn, msg interface{}) {
	conn.WriteJSON(s)
}

func sendSockets(eventName string, users []int, msg interface{}) {
	s := &socketData{
		EventName: eventName,
		Msg:       msg,
	}
	if len(users) == 0 {
		for _, conns := range sockets {
			for _, conn := range conns {
				sendSocket(s, conn, msg)
			}
		}
	} else {
		for _, user := range users {
			sendSocketUser(s, user, msg)
		}
	}
}
