package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", httpHandler)

	panic(http.ListenAndServe(":8080", nil))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go doEvery(2*time.Second, pushCurrentTime, conn)
}

func doEvery(d time.Duration, f func(time.Time, *websocket.Conn), conn *websocket.Conn) {
	for x := range time.Tick(d) {
		f(x, conn)
	}
}

func pushCurrentTime(t time.Time, conn *websocket.Conn) {
	err := conn.WriteJSON(t.Format("3:4:5 PM"))

	if err != nil {
		fmt.Println(err)
	}
}
