package handlers

import (
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	for {
		if err != nil {
			return 
		}
		msg := []byte("hi");
		if err := conn.WriteMessage(1, msg); err != nil {
			return
		}
	}
}
