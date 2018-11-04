package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/saresend/calhacks5-server/handlers"
	"github.com/saresend/calhacks5-server/state"
)

func main() {

	state.Init()
	go Tick()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/updates", handlers.SocketHandler)
	r.HandleFunc("/getPrompts", handlers.GetPrompts).Methods("GET")
	r.HandleFunc("/setPrompts", handlers.SetPrompts).Methods("POST")
	r.HandleFunc("/vote", handlers.Vote).Methods("POST")
	http.ListenAndServe(":"+port, r)
}

func Tick() {
	c := time.Tick(time.Second)
	for {
		<-c
		state.UpdateState()
		handlers.Broadcast()
	}
}
