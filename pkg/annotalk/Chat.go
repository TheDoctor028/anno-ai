package annotalk

import "github.com/TheDoctor028/annotalk-chatgpt/pkg/socketIO"

type Chat struct {
	filterStats bool
	client      *socketIO.Client
}

func NewChat(filterStats bool, client *socketIO.Client) *Chat {
	return &Chat{
		filterStats: true,
		client:      client,
	}
}

func (c *Chat) Start(self Person) {

}
