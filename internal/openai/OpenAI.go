package openai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func GenerateReply(httpClient *http.Client, content string) (string, error) {

	openAIToken := os.Getenv("OPENAI_TOKEN")
	assistantID := os.Getenv("OPENAI_ASSISTANT_ID")
	if openAIToken == "" || assistantID == "" {
		fmt.Println("Error finding openai api token or assistant ID")
		return "", errors.New("error finding OpenAI token: ")
	}

	openAIClient := openai.NewClient(openAIToken)

	openAIThread, err := openAIClient.CreateThread(context.Background(), openai.ThreadRequest{})
	if err != nil {
		return "", errors.New("Error creating OpenAI Thread: " + err.Error())
	}

	_, err = openAIClient.CreateMessage(
		context.Background(),
		openAIThread.ID,
		openai.MessageRequest{
			Role:    "user",
			Content: content,
		},
	)
	if err != nil {
		return "", errors.New("Error creasting OpenAI message: " + err.Error())
	}

	runresp, err := openAIClient.CreateRun(
		context.Background(),
		openAIThread.ID,
		openai.RunRequest{
			AssistantID: assistantID,
			Model:       "gpt-4",
		},
	)
	if err != nil {
		return "", errors.New("Error running message: " + err.Error())
	}

	threadStatus := "queued"
	numberOfStatusChecks := 0

	for threadStatus == "in_progress" || threadStatus == "queued" {
		time.Sleep(1 * time.Second)
		runresp, err = openAIClient.RetrieveRun(
			context.Background(),
			openAIThread.ID,
			runresp.ID,
		)
		if err != nil {
			return "", errors.New("Error retrieving run status: " + err.Error())
		}
		threadStatus = string(runresp.Status)
		fmt.Println("GPT run status:", runresp.Status)
		numberOfStatusChecks += 1
		if numberOfStatusChecks > 15 {
			return "", errors.New("max number of run status checks reached")
		}
	}

	messages, err := openAIClient.ListMessage(
		context.Background(),
		openAIThread.ID,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return "", errors.New("Error lising openai messages: " + err.Error())
	}

	if len(messages.Messages) > 0 && len(messages.Messages[0].Content) > 0 &&
		messages.Messages[0].Content[0].Text != nil {
		return messages.Messages[0].Content[0].Text.Value, nil
	} else {
		return "", errors.New("error getting generated comment content")
	}

}
