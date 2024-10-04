package main

import (
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"log"
	"time"
)

func main() {
	sio, err := socketIO.NewSocketIOClient(
		"husrv.anotalk.hu",
	)
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case msg := <-sio.ReceiveMessage:
			log.Printf("Received message: %s", string(msg))
		case <-ticker.C:
			sio.SendMessage <- []byte("Hello")
		case <-sio.Done:
			log.Println("Connection closed")
			return

		}
	}

}
