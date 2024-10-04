package main

import (
	"fmt"
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/annotalk"
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"github.com/joho/godotenv"
	"log"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	startAnnoTalkChat(err)
}

func startAnnoTalkChat(err error) {
	sio, err := socketIO.NewSocketIOClient(
		"husrv.anotalk.hu",
	)
	if err != nil {
		panic(err)
	}

	chat := annotalk.NewChat(true, sio)
	chat.StartNewChat(annotalk.Persona{
		Name:               "Viktor",
		Age:                25,
		Gender:             annotalk.Man,
		InterestedInGender: annotalk.Whatever,
		Description:        "",
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
