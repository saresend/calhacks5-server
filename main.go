package main

import (
	"net/http"
)

func main() {

	Init()
	go Tick()
	http.HandleFunc("/updates", SocketHandler)
	http.ListenAndServe(":8080", nil)
}
