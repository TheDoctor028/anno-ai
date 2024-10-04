package socketIO

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	ws *websocket.Conn

	ReceiveMessage chan []byte
	SendMessage    chan []byte
	Done           chan struct{}
	pong           chan struct{}
	dialer         *websocket.Dialer
	rest           *resty.Client
}

func NewSocketIOClient(host string) (*Client, error) {
	socketIOConfig, err := getSID(host)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf("wss://%s/socket.io/?EIO=3&transport=websocket&sid=%s", host, socketIOConfig.Sid))
	if err != nil {
		return nil, err
	}

	log.Printf("Connecting to %s", u.String())

	d := &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  5 * time.Second,
		Subprotocols:      []string{"soap"},
		EnableCompression: false,
	}

	ws, res, err := d.Dial(u.String(), nil)
	if err != nil {
		if res != nil {
			defer res.Body.Close()
			bs, _ := io.ReadAll(res.Body)
			log.Printf("WS Dial respons: %s", string(bs))
		}
		log.Fatal(err)
	}

	c := &Client{
		ReceiveMessage: make(chan []byte),
		SendMessage:    make(chan []byte),

		ws:     ws,
		dialer: d,
		Done:   make(chan struct{}),
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
		log.Println("Sending ACKHELLO")
		if err = ws.WriteMessage(websocket.TextMessage, []byte(ACKHELLO)); err != nil {
			log.Println("Error sending ACKHELLO: ", err)
		}
		log.Println("Sent ACKHELLO")

		go c.handleIncoming()
		go c.startPing(time.Duration(socketIOConfig.PingInterval)*time.Millisecond, c.Done)
		ws.SetCloseHandler(func(code int, text string) error {
			c.Done <- struct{}{}
			return nil
		})
	}

	return c, nil
}

func getSID(host string) (*struct {
	Sid          string `json:"sid,omitempty"`
	PingInterval int    `json:"pingInterval,omitempty"`
	PingTimeout  int    `json:"pingTimeout,omitempty"`
}, error) {
	r := resty.New()

	initRes, err := r.R().Get(fmt.Sprintf("https://%s/socket.io/?EIO=3&transport=polling&t=P9JsPIX", host))
	if err != nil {
		return nil, err
	}

	if initRes == nil || initRes.StatusCode() != http.StatusOK {
		return nil, err
	}

	woPrefix, _ := strings.CutPrefix(initRes.String(), "96:0")
	jsonString, _ := strings.CutSuffix(woPrefix, "2:40")

	var socketIOConfig struct {
		Sid          string `json:"sid,omitempty"`
		PingInterval int    `json:"pingInterval,omitempty"`
		PingTimeout  int    `json:"pingTimeout,omitempty"`
	}
	err = json.Unmarshal([]byte(jsonString), &socketIOConfig)
	if err != nil {
		return nil, err
	}
	return &socketIOConfig, err
}

func (c *Client) startPing(healthCheckTime time.Duration, done chan struct{}) {
	ticker := time.NewTicker(healthCheckTime)
	defer c.ws.Close()

	for {
		select {
		case <-ticker.C:
			if err := c.ws.WriteMessage(websocket.TextMessage, []byte(PING)); err != nil {
				log.Printf("Error sending ping: %v", err)
				return
			}
			log.Println("Ping ---->")
		case <-done:
			return
		}
	}

}

func (c *Client) handleIncoming() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
		}

		if string(msg) == PONG {
			log.Println("<---- Pong")
			continue
		}
		if strings.HasPrefix(string(msg), MESSAGE_SERVER) {
			c.ReceiveMessage <- []byte(strings.TrimPrefix(string(msg), MESSAGE_SERVER))
		} else {
			log.Printf("Received unknow type message: %s", string(msg))
		}
	}
}
