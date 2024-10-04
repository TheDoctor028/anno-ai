package annotalk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
	persona Persona

	instructions string
}

func NewAI(persona Persona, partnerGender PersonGender) (*AI, error) {
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
		PartnerGender PersonGender
		Persona
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

func (a *AI) GetAnswer(messages []Message) (string, error) {
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

func mapMessagesToOpenAIMessages(messages []Message) []openai.ChatCompletionMessage {
	var openAIMessages []openai.ChatCompletionMessage
	for _, msg := range messages {
		if len(msg.Msg) == 0 {
			continue
		}

		if msg.Entity == Partner {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: msg.Msg,
			})
		}
		if msg.Entity == Bot || msg.Entity == User {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: msg.Msg,
			})
		}
	}

	return openAIMessages
}
