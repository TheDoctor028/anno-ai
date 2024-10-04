package main

import (
	"fmt"
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

	go func() {
		log.Println("Type your message to send")
		for {
			msg := ""
			_, err := fmt.Scanln(&msg)
			if err != nil {
				log.Println(err)
			}

			chat.SendMessage(msg)
		}
	}()

	for {
		select {
		case <-sio.Done:
			log.Println("Connection closed")
			return

		}
	}

}
