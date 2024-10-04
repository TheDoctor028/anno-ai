package main

import (
	"fmt"
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/annotalk"
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"log"
	"strings"
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
		Gender:             annotalk.Woman,
		InterestedInGender: annotalk.Man,
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
			if strings.Compare(msg, "exit") == 0 {
				chat.EndChat()
				continue
			}

			if strings.Compare(msg, "start") == 0 {
				chat.FindNewPartner()
				continue
			}

			chat.SendMessage(msg, annotalk.User)
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
