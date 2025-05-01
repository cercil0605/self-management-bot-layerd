package client

import (
	"bytes"
	"encoding/json"
	"io"
	// "fmt"
	"net/http"
	"os/exec"
)

var ollamaCmd *exec.Cmd

type LLMRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type LLMResponse struct {
	Response string `json:"response"`
}

// StartLLM Ollama serverの起動または終了
func StartLLM() error {
	ollamaCmd = exec.Command("ollama", "serve")
	return ollamaCmd.Start()
}
func StopLLM() error {
	return ollamaCmd.Process.Kill()
}

// GetLLMResponse Response from LLM
func GetLLMResponse(prompt string) (string, error) {
	req := LLMRequest{
		Model:  "Llama-3-ELYZA-JP",
		Prompt: prompt,
		Stream: false,
	}
	body, _ := json.Marshal(req)
	// debug prompt
	// fmt.Println("送信プロンプト:", prompt)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	// close connection
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type LLMResponse struct {
		Response string `json:"response"`
	}
	var llmResponse LLMResponse
	if err := json.Unmarshal(body, &llmResponse); err != nil {
		return "", err
	}
	return llmResponse.Response, nil
}
