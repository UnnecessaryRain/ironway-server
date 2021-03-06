package server

import (
	"net/http"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/client"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveSocket(server *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := client.NewClient(conn, server.receivedChan, server.registerChan, server.unregisterChan)
	go client.StartWriter()
	go client.StartReader()
}
