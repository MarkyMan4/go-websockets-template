package main

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/websocket"
)

type Server struct {
	rooms map[string]map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{rooms: make(map[string]map[*websocket.Conn]bool)}
}

func (s *Server) handleWebSocket(ws *websocket.Conn) {
	// based on URL, can close socket if not valid
	// client can check ws.readyState to see if they are actually connected
	urlParts := strings.Split(ws.Request().URL.Path, "/")

	if len(urlParts) != 3 || len(strings.Trim(urlParts[2], " ")) == 0 {
		fmt.Println(len(urlParts))
		fmt.Println("invalid URL:", ws.Request().URL.Path)
		ws.Close()
		return
	}

	roomId := urlParts[2]
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	fmt.Println("room ID:", roomId)

	if _, ok := s.rooms[roomId]; !ok {
		s.rooms[roomId] = make(map[*websocket.Conn]bool)
		fmt.Println("new room created:", roomId)
	}

	// add custom logic here around client info, checking if room is at capacity, etc.

	s.rooms[roomId][ws] = true
	s.readLoop(ws, roomId)
}

func (s *Server) readLoop(ws *websocket.Conn, roomId string) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				// remove connection once closed by client
				delete(s.rooms[roomId], ws)
				break
			}

			// if no one left in the room, delete the room
			if len(s.rooms[roomId]) == 0 {
				delete(s.rooms, roomId)
				fmt.Printf("room %s deleted", roomId)
			}

			fmt.Println("read error:", err)
			continue
		}

		msg := buf[:n]

		// instead of broadcasting, here is where you can parse the event
		// and decide what to do with it
		s.broadcast(msg, roomId)
	}
}

func (s *Server) broadcast(b []byte, roomId string) {
	for ws := range s.rooms[roomId] {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}
