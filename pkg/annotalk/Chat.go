package annotalk

import (
	"encoding/json"
	"fmt"
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"html"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)
import "github.com/TheDoctor028/annotalk-chatgpt/pkg/utils"

type Entity string

const (
	Partner Entity = "partner"
	User    Entity = "user"
	Bot     Entity = "bot"
)

type Message struct {
	Entity Entity `json:"entity"`
	Msg    string `json:"message"`
}

type Chat struct {
	filterStats    bool
	client         *socketIO.Client
	ai             *AI
	inChat         bool
	alreadyHadChat bool

	conversationsID *string
	partnerGender   *PersonGender
	messages        []Message
	stats           OnStatisticsData
	person          *Persona
	typing          sync.Mutex
	aiResponseTimer *time.Timer

	MessageEventsChannels *MessageEvents
}

func NewChat(filterStats bool, client *socketIO.Client) *Chat {
	c := &Chat{
		MessageEventsChannels: NewMessageEvents(),

		inChat:         false,
		alreadyHadChat: false,
		filterStats:    filterStats,
		client:         client,

		messages:        []Message{},
		person:          nil,
		aiResponseTimer: time.NewTimer(time.Duration((rand.Int()%15)+5) * time.Second),
	}

	go c.messageHandler()
	return c
}

func (c *Chat) StartNewChat(self Persona) {
	if c.inChat {
		log.Println("You are already in a chat")
		return
	}

	c.person = &self
	c.messages = []Message{}
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

	var err error
	c.ai, err = NewAI(self, self.InterestedInGender)
	if err != nil {
		log.Println(err)
		c.EndChat()
	}
}

func (c *Chat) FindNewPartner() {
	if !c.inChat && c.person != nil {
		c.messages = []Message{}
		log.Println("Finding new partner")
		c.client.SendMessage <- socketIO.OutgoingMessage{
			Type: string(LookForPartner),
		}
		<-c.MessageEventsChannels.ChatStart
	}
}

func (c *Chat) SendMessage(msg string, entity Entity) {
	if c.inChat {
		c.client.SendMessage <- socketIO.OutgoingMessage{
			Type: string(SendMessage),
			Data: SendMessageData{
				Message: msg,
			},
		}
		c.messages = append(c.messages, Message{Entity: entity, Msg: msg})
		if entity == Bot {
			log.Println("Bot: ", msg)
		}
	} else {
		log.Println("You are not in a chat")
	}
}

func (c *Chat) EndChat() {
	if c.inChat {
		c.client.SendMessage <- socketIO.OutgoingMessage{
			Type: string(LeaveChat),
			Data: map[string]interface{}{},
		}
		<-c.MessageEventsChannels.ChatEnd
		log.Println("Chat ended")
	}
}

func (c *Chat) messageHandler() {
	for {
		select {
		case msg := <-c.client.ReceiveMessage:
			switch msg.Type {
			case string(OnStatistics):
				if !c.filterStats {
					log.Printf("Statistics: %v", msg.Data)
				}
				c.stats = NewOnStatisticsData(msg.Data)
			case string(OnChatStart):
				c.onChatStart(msg)
			case string(OnMessage):
				c.onMessage(msg)
			case string(OnChatEnd):
				c.onChatEnd()
			case string(OnSearchingPartner):
				log.Println("Searching for partner")
				go func() { c.MessageEventsChannels.SearchingPartner <- struct{}{} }()
			}
		case <-c.aiResponseTimer.C:
			go c.sendAIMessage()
		}
	}
}

func (c *Chat) onChatStart(msg socketIO.IncomingMessage) {
	data := NewOnChatStartData(msg.Data)
	c.inChat = true
	c.conversationsID = &data.ChatID
	c.partnerGender = &data.PartnerGender
	log.Printf("Chat started with a %s", NewOnChatStartData(msg.Data).PartnerGender)
	go func() { c.MessageEventsChannels.ChatStart <- NewOnChatStartData(msg.Data) }()
}

func (c *Chat) onChatEnd() {
	log.Println("Chat ended")

	msgsJson, err := json.Marshal(struct {
		Id            string       `json:"id"`
		Timestamp     string       `json:"timestamp"`
		Person        Persona      `json:"person"`
		PartnerGender PersonGender `json:"partnerGender"`
		Messages      []Message    `json:"messages"`
	}{
		Id:            *c.conversationsID,
		Timestamp:     time.Now().Format(time.RFC3339),
		Person:        *c.person,
		PartnerGender: *c.partnerGender,
		Messages:      c.messages,
	})
	if err != nil {
		log.Printf("Error marshalling messages %s", err)
	}

	fileName := fmt.Sprintf("data/conversations/%d.json", time.Now().Unix())
	if err := os.WriteFile(fileName, msgsJson, 0644); err != nil {
		log.Printf("Error writing conversation to file %s %s", fileName, err)
	}

	c.conversationsID = nil
	c.inChat = false
	go func() { c.MessageEventsChannels.ChatEnd <- struct{}{} }()
}

func (c *Chat) onMessage(msg socketIO.IncomingMessage) {
	if NewOnMessageData(msg.Data).IsYou == 0 {
		msgTxt := html.UnescapeString(NewOnMessageData(msg.Data).Message)
		log.Printf("Partner: %v", msgTxt)
		c.messages = append(c.messages, Message{Entity: Partner, Msg: msgTxt})
		c.aiResponseTimer.Reset(time.Duration((rand.Int()%5)+5) * time.Second)
		go func() { c.MessageEventsChannels.Message <- NewOnMessageData(msg.Data) }()
	}
}

func (c *Chat) sendAIMessage() {
	if c.inChat {
		if c.typing.TryLock() {
			defer c.typing.Unlock()
			c.client.SendMessage <- socketIO.OutgoingMessage{
				Type: string(Typing),
			}
			msg, err := c.ai.GetAnswer(c.messages)
			if err != nil {
				log.Printf("Error getting answer %s", err)
			} else {
				c.SendMessage(msg, Bot)
				c.client.SendMessage <- socketIO.OutgoingMessage{
					Type: string(DoneTyping),
				}
			}
		}
	}
}
