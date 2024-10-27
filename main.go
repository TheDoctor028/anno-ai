package main

import (
	"fmt"
	"github.com/TheDoctor028/anno-ai/pkg/annotalk"
	"github.com/TheDoctor028/anno-ai/pkg/socketIO"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	personas := []annotalk.Persona{ // TODO read this from a file
		{
			Name:               "Lali",
			Age:                25,
			Gender:             annotalk.Man,
			InterestedInGender: annotalk.Whatever,
			Description:        "",
		},
		//{
		//	Name:               "Judit",
		//	Age:                23,
		//	Gender:             annotalk.Woman,
		//	InterestedInGender: annotalk.Whatever,
		//	Description:        "Kicsit felvállalós, kicsit szégyenlős",
		//},
		{
			Name:               "Emese",
			Age:                18,
			Gender:             annotalk.Woman,
			InterestedInGender: annotalk.Woman,
			Description:        "Szeretek olvasni, és a természetben lenni",
		},
		//{
		//	Name:               "Jimmy",
		//	Age:                19,
		//	Gender:             annotalk.Man,
		//	InterestedInGender: annotalk.Man,
		//	Description:        "Gamer, kedvenc játékom a Minecraft",
		//},
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
					chat.SaveChat(fmt.Sprintf("%d.bak", chat.GetTS().Unix()))
					if chat.Client.IsConnected() {
						chat.EndChat()
					}
				}
				return
			}
		}
	}
}

func startAnnoTalkChat(persona annotalk.Persona) *annotalk.Chat {
	sio, err := socketIO.NewSocketIOClient(
		os.Getenv("ANO_TALK_WEBSOCKET_HOST"),
	)
	if err != nil {
		panic(err)
	}

	chat := annotalk.NewChat(true, sio, true)
	chat.StartNewChat(persona)
	return chat
}
