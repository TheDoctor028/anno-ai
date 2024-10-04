package annotalk

import (
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"
	"log"
)
import "github.com/TheDoctor028/annotalk-chatgpt/pkg/utils"

type Chat struct {
	filterStats    bool
	client         *socketIO.Client
	alreadyHadChat bool
}

func NewChat(filterStats bool, client *socketIO.Client) *Chat {
	return &Chat{

		alreadyHadChat: false,
		filterStats:    true,
		client:         client,
	}
}

func (c *Chat) StartNewChat(self Person) {
	log.Printf("Starting new chat as %s(%d %s) to talk with %s", self.Name, self.Age, self.Gender, self.InterestedInGender)
	c.client.SendMessage <- socketIO.Message{
		Type: InitChat,
		Data: InitChatData{
			Gender:        self.Gender,
			PartnerGender: self.InterestedInGender,
			CaptchaID:     utils.RandStringRunes(20), // TODO investigate this
		},
	}
}
