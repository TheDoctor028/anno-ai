package main

import (
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/annotalk"
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"log"
)

func main() {
	sio, err := socketIO.NewSocketIOClient(
		"husrv.anotalk.hu",
	)
	if err != nil {
		panic(err)
	}

	chat := annotalk.NewChat(true, sio)
	chat.StartNewChat(annotalk.Person{
		Name:               "Viktor",
		Age:                27,
		Gender:             annotalk.Man,
		InterestedInGender: annotalk.Whatever,
		Description:        "I'm a bot",
	})

	for {
		select {
		case msg := <-sio.ReceiveMessage:
			log.Printf("Received message: %s", msg.Type)
		case <-sio.Done:
			log.Println("Connection closed")
			return

		}
	}

}
