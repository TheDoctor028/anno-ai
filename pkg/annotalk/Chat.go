package annotalk

import (
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"log"
)
import "github.com/TheDoctor028/annotalk-chatgpt/pkg/utils"

const (
	Partner = "partner"
	User    = "user"
	Bot     = "bot"
)

type Message struct {
	Entity string `json:"entity"`
	Msg    string `json:"message"`
}

type Chat struct {
	filterStats    bool
	client         *socketIO.Client
	alreadyHadChat bool
	messages       []Message

	MessageEventsChannels *MessageEvents
}

func NewChat(filterStats bool, client *socketIO.Client) *Chat {
	c := &Chat{
		MessageEventsChannels: NewMessageEvents(),

		alreadyHadChat: false,
		filterStats:    filterStats,
		client:         client,

		messages: []Message{},
	}

	go c.MessageHandler()
	return c
}

func (c *Chat) StartNewChat(self Person) {
	log.Printf("Starting new chat as %s(%d %s) to talk with %s", self.Name, self.Age, self.Gender, self.InterestedInGender)
	c.client.SendMessage <- socketIO.OutgoingMessage{
		Type: string(InitChat),
		Data: InitChatData{
			Gender:        self.Gender,
			PartnerGender: self.InterestedInGender,
			CaptchaID:     utils.RandStringRunes(20),
		},
	}
	<-c.MessageEventsChannels.SearchingPartner
	<-c.MessageEventsChannels.ChatStart
	c.alreadyHadChat = true
}

func (c *Chat) MessageHandler() {
	for {
		select {
		case msg := <-c.client.ReceiveMessage:
			switch msg.Type {
			case string(OnStatistics):
				if !c.filterStats {
					log.Printf("Statistics: %v", msg.Data)
				}
				go func() { c.MessageEventsChannels.Stats <- NewOnStatisticsData(msg.Data) }()
			case string(OnChatStart):
				log.Printf("Chat started with a %s", NewOnChatStartData(msg.Data).PartnerGender)
				go func() { c.MessageEventsChannels.ChatStart <- NewOnChatStartData(msg.Data) }()
			case string(OnMessage):
				if NewOnMessageData(msg.Data).IsYou == 0 {
					msgTxt := NewOnMessageData(msg.Data).Message
					log.Printf("Your partner: %v", msgTxt)
					c.messages = append(c.messages, Message{Entity: Partner, Msg: msgTxt})
					go func() { c.MessageEventsChannels.Message <- NewOnMessageData(msg.Data) }()
				}
			case string(OnChatEnd):
				log.Println("Chat ended")
				go func() { c.MessageEventsChannels.ChatEnd <- struct{}{} }()
			case string(OnSearchingPartner):
				log.Println("Searching for partner")
				go func() { c.MessageEventsChannels.SearchingPartner <- struct{}{} }()
			}
		}
	}
}
