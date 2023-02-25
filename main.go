package main

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	OpenAIApiKey        string   `yaml:"openai_api_key"`
	GptModel            string   `yaml:"gpt_model"`
	ActivationPhrase    string   `yaml:"activation_phrase"`
	ActivationResponses []string `yaml:"activation_responses"`
	VoiceName           string   `yaml:"voice_name"`
}

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	fmt.Println("Listening for activation phrase...")

	say(getRandomActivationResponse(config))

	gptClient := gogpt.NewClient(config.OpenAIApiKey)
	ctx := context.Background()

	gptResponse := gptReq(gptClient, ctx, config, "What is golang?")

	say(gptResponse)

}

func readConfig() (*Config, error) {
	content, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func say(text string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("say '%s'", text))
	return cmd.Run()
}

func gptReq(c *gogpt.Client, ctx context.Context, config *Config, prompt string) string {
	var model = (*config).GptModel

	req := gogpt.CompletionRequest{
		Model:     model,
		MaxTokens: 1024,
		Prompt:    prompt,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "Error"
	}

	fmt.Println("GPT Response: ", resp)

	var rl = len(resp.Choices)
	if rl < 1 {
		return "No result"
	}
	choices := make([]string, rl)

	for j, ch := range resp.Choices {
		choices[j] = ch.Text
	}

	return strings.Join(choices, ". ")
}

func getRandomActivationResponse(config *Config) string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len((*config).ActivationResponses))
	return (*config).ActivationResponses[randomIndex]
}
