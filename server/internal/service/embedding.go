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

type EmbeddingProvider interface {
	Embed(texts []string) ([][]float32, error)
	Name() string
}

type deepSeekEmbedding struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewDeepSeekEmbedding(apiKey, baseURL string) EmbeddingProvider {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	return &deepSeekEmbedding{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   "deepseek-chat",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (d *deepSeekEmbedding) Name() string {
	return "deepseek"
}

func (d *deepSeekEmbedding) Embed(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	reqBody := map[string]interface{}{
		"model": "deepseek-embedding",
		"input": texts,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", d.baseURL+"/embeddings", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	vectors := make([][]float32, len(result.Data))
	for i, item := range result.Data {
		vec := make([]float32, len(item.Embedding))
		for j, v := range item.Embedding {
			vec[j] = float32(v)
		}
		vectors[i] = vec
	}
	return vectors, nil
}

type openAIEmbedding struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewOpenAIEmbedding(apiKey, baseURL string) EmbeddingProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &openAIEmbedding{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   "text-embedding-3-small",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (o *openAIEmbedding) Name() string {
	return "openai"
}

func (o *openAIEmbedding) Embed(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	type embeddingRequest struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}
	reqBody := embeddingRequest{
		Model: o.model,
		Input: texts,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", o.baseURL+"/embeddings", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	vectors := make([][]float32, len(result.Data))
	for i, item := range result.Data {
		vec := make([]float32, len(item.Embedding))
		for j, v := range item.Embedding {
			vec[j] = float32(v)
		}
		vectors[i] = vec
	}
	return vectors, nil
}
