// websocket_im.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/websocket"
)

var clients []*websocket.Conn

func server() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
		conn, _ := upgrader.Upgrade(w, r, nil)
		clients = append(clients, conn)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			for _, client := range clients {
				client.WriteMessage(websocket.TextMessage, msg)
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func client() {
	conn, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	go func() {
		for {
			_, msg, _ := conn.ReadMessage()
			fmt.Println("Received:", string(msg))
		}
	}()
	var input string
	for {
		fmt.Scanln(&input)
		if input == "exit" {
			break
		}
		conn.WriteMessage(websocket.TextMessage, []byte(input))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run websocket_im.go [server|client]")
		return
	}
	if os.Args[1] == "server" {
		server()
	} else {
		client()
	}
}
