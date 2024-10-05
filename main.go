package main

import (
	"fmt"
	"github.com/TheDoctor028/anno-ai/pkg/annotalk"
	"github.com/TheDoctor028/anno-ai/pkg/socketIO"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	personas := []annotalk.Persona{
		{
			Name:               "Alice",
			Age:                25,
			Gender:             annotalk.Man,
			InterestedInGender: annotalk.Woman,
			Description:        "",
		},
	}

	var chats []*annotalk.Chat

	for _, persona := range personas {
		go func() {
			chat := startAnnoTalkChat(persona)
			chats = append(chats, chat)
		}()
	}

	for {
		log.Println("Type 'exit' to end all chats")
		for {
			msg := ""
			_, err := fmt.Scanln(&msg)
			if err != nil {
				log.Println(err)
			}
			if msg == "exit" {
				for _, chat := range chats {
					chat.EndChat()
				}
				return
			}
		}
	}
}

func startAnnoTalkChat(persona annotalk.Persona) *annotalk.Chat {
	sio, err := socketIO.NewSocketIOClient(
		"husrv.anotalk.hu",
	)
	if err != nil {
		panic(err)
	}

	chat := annotalk.NewChat(true, sio, true)
	chat.StartNewChat(persona)
	return chat
}
