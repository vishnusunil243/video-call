package main

import (
	"log"
	"net/http"

	"main.go/server"
)

func main() {
	server.AllRooms.Init()
	http.Handle("/create", http.HandlerFunc(server.CreateRoomRequestHandler))
	http.Handle("/join", http.HandlerFunc(server.JoinRoomRequestHandler))
	go server.Broadcaster()
	log.Println("starting servier on port 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
