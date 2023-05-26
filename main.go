package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()
	http.Handle("/ws/", websocket.Handler(server.handleWebSocket))

	fmt.Println("listening on port 8000...")
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
