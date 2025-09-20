package client

import (
	"context"
	"fmt"

	"google.golang.org/genai"
	"self-management-bot/config"
)

func GetGeminiResponse(prompt string) (string, error) {
	apiKey := config.Cfg.GeminiApiKey
	if apiKey == "" {
		return "", fmt.Errorf("Gemini API key is not set")
	}

	ctx := context.Background()
	cl, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("genai client init: %w", err)
	}

	res, err := cl.Models.GenerateContent(ctx, "gemini-2.5-flash",
		[]*genai.Content{{Parts: []*genai.Part{{Text: prompt}}}},
		nil)
	if err != nil {
		return "", fmt.Errorf("generate: %w", err)
	}
	txt := res.Text()
	if txt == "" {
		return "", fmt.Errorf("empty response")
	}
	return txt, nil
}
