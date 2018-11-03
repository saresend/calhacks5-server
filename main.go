package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/saresend/calhacks5-server/handlers"
	"github.com/saresend/calhacks5-server/state"
)

func main() {

	state.Init()
	go Tick()

	r := mux.NewRouter()
	r.HandleFunc("/updates", handlers.SocketHandler)
	r.HandleFunc("/getPrompts", handlers.GetPrompts).Methods("GET")
	r.HandleFunc("/setPrompts", handlers.SetPrompts).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func Tick() {
	c := time.Tick(time.Second)
	for {
		<-c
		state.UpdateState()
		handlers.Broadcast()
	}
}
