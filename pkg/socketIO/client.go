package socketIO

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type SocketIOClient struct {
	ws *websocket.Conn

	ReceiveMessage chan []byte
	SendMessage    chan []byte
	done           chan struct{}
	dialer         websocket.Dialer
}

func NewSocketIOClient(host string, queryParams string, healthCheckTimed time.Duration) *SocketIOClient {

	u, err := url.Parse(fmt.Sprintf("wss://%s/socket.io/%s", host, queryParams))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connecting to %s", u.String())

	d := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  30 * time.Second,
		Subprotocols:      []string{"soap"},
		EnableCompression: true,
	}

	h := http.Header{}
	h.Add("Accept", "*/*")

	ws, res, err := d.Dial(u.String(), h)
	if err != nil {
		if res != nil {
			defer res.Body.Close()
			bs, _ := io.ReadAll(res.Body)
			log.Printf("WS Dial respons: %s", string(bs))
		}
		log.Fatal(err)
	}

	c := &SocketIOClient{
		ReceiveMessage: make(chan []byte),
		SendMessage:    make(chan []byte),

		ws:     ws,
		dialer: d,
		done:   make(chan struct{}),
	}

	log.Printf("Connected to %s", u.String())

	log.Println("Sending CHELLO")
	if err = ws.WriteMessage(websocket.TextMessage, []byte(CHELLO)); err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting for SHELLO")

	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Println("Error reading message: ", err)
	}

	if string(msg) == SHELLO {
		log.Println("Received SHELLO")
		if err = ws.WriteMessage(websocket.TextMessage, []byte(ACKHELLO)); err != nil {
			log.Println("Error sending ACKHELLO: ", err)
		}
		c.startPing(healthCheckTimed, c.done)
	}

	return c
}

func (c *SocketIOClient) startPing(healthCheckTime time.Duration, done chan struct{}) {
	ticker := time.NewTicker(healthCheckTime)
	defer c.ws.Close()

	for {
		select {
		case <-ticker.C:
			if err := c.ws.WriteMessage(websocket.TextMessage, []byte(PING)); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
		case <-done:
			return
		}
	}

}
