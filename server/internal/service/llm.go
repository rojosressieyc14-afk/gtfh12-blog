package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type llmChatProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewDeepSeekLLM(apiKey, baseURL string) LLMProvider {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	return &llmChatProvider{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   "deepseek-chat",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func NewOpenAILLM(apiKey, baseURL string) LLMProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &llmChatProvider{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   "gpt-4o-mini",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (l *llmChatProvider) Name() string {
	return l.model
}

func (l *llmChatProvider) Chat(systemPrompt, userMessage string) (string, error) {
	reqBody := map[string]interface{}{
		"model": l.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userMessage},
		},
		"temperature": 0.3,
		"max_tokens":  2048,
		"stream":      false,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", l.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.apiKey)

	resp, err := l.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("LLM request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("LLM returned empty result")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
