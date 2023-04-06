package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func getResponse(ctx context.Context, client *openai.Client, quesiton, content string) error {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are ChatGPT, Answer as concisely as possible.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("task_description: %s \n\n Content: %s", quesiton, content),
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err
	}

	fmt.Println(resp.Choices[0].Message.Content)
	return nil
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))

	if len(os.Args) < 2 {
		panic("Missing argument: <question>")
	}

	desc := os.Args[1]

	config := openai.DefaultConfig(mustKey())
	if v := orEnv("OPENAI_ENDPOINT", ""); v != "" {
		config.BaseURL = v
	}
	c := openai.NewClientWithConfig(config)

	content, _ := io.ReadAll(os.Stdin)
	if err := getResponse(context.Background(), c, desc, string(content)); err != nil {
		panic(err)
	}
}

func mustKey() string {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		panic("Missing OPENAI_API_KEY=")
	}
	return key
}

func orEnv(k, dv string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return dv
}
