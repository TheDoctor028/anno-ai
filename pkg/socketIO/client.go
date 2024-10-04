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

	ReceiveMessage chan IncomingMessage
	SendMessage    chan OutgoingMessage
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
		ReceiveMessage: make(chan IncomingMessage),
		SendMessage:    make(chan OutgoingMessage),

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

		go c.handelOutgoing()
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
			return
		}

		if string(msg) == PONG {
			continue
		}
		if strings.HasPrefix(string(msg), MESSAGE) {
			msgWoPrefix, _ := strings.CutPrefix(string(msg), MESSAGE)
			var m []interface{}
			err = json.Unmarshal([]byte(msgWoPrefix), &m)
			if err != nil {
				log.Printf("Error unmarshaling message: %s %v", string(msg), err)
			}

			c.ReceiveMessage <- IncomingMessage{
				Type: m[0].(string),
				Data: m[1].(map[string]interface{}),
			}
		} else {
			log.Printf("Received unknow type message: %s", string(msg))
		}
	}
}

func (c *Client) handelOutgoing() {
	for {
		select {
		case msg := <-c.SendMessage:
			data, err := json.Marshal(msg.Data)
			if err != nil {
				log.Println("Error marshaling message: ", err)
			}
			message := []byte(fmt.Sprintf("%s[\"%s\",%s]", MESSAGE, msg.Type, data))

			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error sending message: ", err)
			}
		case <-c.Done:
			return
		}
	}

}
