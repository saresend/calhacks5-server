package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var connections = make(map[string]*websocket.Conn)

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	user := new(User)
	conn.ReadJSON(user)

	connections[user.GetUsername()] = conn

}

func Broadcast(newState PublicState) {
	println("hey")
	for id, conn := range connections {
		if err := conn.WriteJSON(newState); err != nil {
			delete(connections, id)
			RemoveUser(id)
		}
	}
}
