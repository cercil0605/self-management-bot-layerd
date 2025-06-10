package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

var ollamaCmd *exec.Cmd

type LLMRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func StartLLM() error {
	if os.Getenv("LLM_PROVIDER") == "gemini" {
		return nil
	}
	ollamaCmd = exec.Command("ollama", "serve")
	return ollamaCmd.Start()
}

func StopLLM() error {
	if os.Getenv("LLM_PROVIDER") == "gemini" {
		return nil
	}
	if ollamaCmd == nil {
		return nil
	}
	return ollamaCmd.Process.Kill()
}

func GetLLMResponse(prompt string) (string, error) {
	if os.Getenv("LLM_PROVIDER") == "gemini" {
		return getGeminiResponse(prompt)
	}
	return getElyzaResponse(prompt)
}

func getElyzaResponse(prompt string) (string, error) {
	req := LLMRequest{
		Model:  "Llama-3-ELYZA-JP",
		Prompt: prompt,
		Stream: false,
	}
	body, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var llmResponse struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &llmResponse); err != nil {
		return "", err
	}
	return llmResponse.Response, nil
}

func getGeminiResponse(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set")
	}
	req := struct {
		Contents []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"contents"`
	}{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{Parts: []struct {
				Text string `json:"text"`
			}{{Text: prompt}}},
		},
	}
	body, _ := json.Marshal(req)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var res struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return res.Candidates[0].Content.Parts[0].Text, nil
}
