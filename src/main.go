package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("Client connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("Notice:", msg)
		s.Emit("reply", "Received: "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "Received: " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last, _ := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, err error) {
		fmt.Println("Error:", err)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Client disconnected:", s.ID(), "Reason:", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("SocketIO server error: %v", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))

	log.Println("Server running at http://localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
