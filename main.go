package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func ping(out chan string) {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			out <- strconv.Itoa(rand.Intn(9999))
		}
	}
}

func serveWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		upgrader := websocket.Upgrader{}
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}
		defer ws.Close()

		data := make(chan string)
		go ping(data)

		for {
			err = ws.WriteMessage(1, []byte(<-data))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}

func main() {
	router := gin.Default()

	router.StaticFile("/", "client.html")
	router.GET("/ws", serveWs())

	router.Run(":5000")
}
