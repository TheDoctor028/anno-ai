package ai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/TheDoctor028/annotalk-chatgpt/pkg/annotalk"
	"github.com/sashabaranov/go-openai"
	"os"
	"text/template"
)

const baseInstructionsHUN = `
Légy a beszélgető partnerem, mint ha cheten beszélgetnénk úgy, 
hogy nem ismerjuk egymást és semmit nem tudunk egymásról. 
Te egy {{ .Age }} éves {{ .Gender }} vagy aki {{ .InterestedInGender }} szeretne beszélgetni. 
A partneredről nem tudsz semmit csak a nemét ami {{ .PartnerGender }}. 
A válaszaid olyanok legyenek, mint ha chaten beszélnénk.
Csak egyszeru smiley-kat használj pl. :) vagy :D stb.
{{- if .Description }} 
Egy rövid leírás rólad:
{{ .Description }}
{{- end }}
`

type AI struct {
	client  *openai.Client
	persona annotalk.Persona

	instructions string
}

func NewAI(persona annotalk.Persona, partnerGender annotalk.PersonGender) (*AI, error) {
	token := os.Getenv("CHAT_GPT_TOKEN")
	if token == "" {
		return nil, errors.New("CHAT_GPT_TOKEN is not provided")
	}

	client := openai.NewClient(token)

	baseTmplt, err := template.New("baseInstructions").Parse(baseInstructionsHUN)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	err = baseTmplt.Execute(buf, struct {
		PartnerGender annotalk.PersonGender
		annotalk.Persona
	}{
		Persona:       persona,
		PartnerGender: partnerGender,
	})
	if err != nil {
		return nil, err
	}

	return &AI{
		client:       client,
		persona:      persona,
		instructions: buf.String(),
	}, nil
}

func (a *AI) GetAnswer(messages []annotalk.Message) (string, error) {
	resp, err := a.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini, // TODO might change this
			Messages: append([]openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: a.instructions,
				},
			}, mapMessagesToOpenAIMessages(messages)...),
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func mapMessagesToOpenAIMessages(messages []annotalk.Message) []openai.ChatCompletionMessage {
	var openAIMessages []openai.ChatCompletionMessage
	for _, msg := range messages {
		if len(msg.Msg) == 0 {
			continue
		}

		if msg.Entity == annotalk.Partner {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: msg.Msg,
			})
		}
		if msg.Entity == annotalk.Bot || msg.Entity == annotalk.User {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: msg.Msg,
			})
		}
	}

	return openAIMessages
}
