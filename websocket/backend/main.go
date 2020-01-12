package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"time"
	"math/rand"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main() {
	fmt.Println("App has started")
	http.HandleFunc("/ws", wsEntrypoint)
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Hello there")
}

func wsEntrypoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func (r *http.Request) bool {return true}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Connections successfull")

	go reader(ws)
	go sender(ws)
}

func sender(ws *websocket.Conn) {
	for {
		time.Sleep(500 * time.Millisecond)
		ws.WriteMessage(websocket.TextMessage,[]byte(fmt.Sprintf("%d", rand.Float32())))
	}
}

func reader(ws *websocket.Conn){
	for {
		mType,msg, err :=  ws.ReadMessage()
		if err != nil {
			log.Println("Failed to read message")
			continue
		}

		if err := ws.WriteMessage(mType, msg); err != nil {
			log.Println("Failed to send message")
			log.Println(err)
		}
	}
}