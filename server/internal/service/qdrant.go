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

type QdrantPoint struct {
	ID      uint64                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

type qdrantSearchResult struct {
	ID     uint64                 `json:"id"`
	Score  float64                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

type QdrantClient struct {
	baseURL  string
	apiKey   string
	client   *http.Client
}

func NewQdrantClient(baseURL, apiKey string) *QdrantClient {
	return &QdrantClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (q *QdrantClient) collectionURL(name string) string {
	return fmt.Sprintf("%s/collections/%s", q.baseURL, name)
}

func (q *QdrantClient) pointsURL(name string) string {
	return fmt.Sprintf("%s/collections/%s/points", q.baseURL, name)
}

func (q *QdrantClient) do(method, url string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if q.apiKey != "" {
		req.Header.Set("api-key", q.apiKey)
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qdrant request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("qdrant error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Result json.RawMessage `json:"result"`
		Status string          `json:"status"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	return []byte(result.Result), nil
}

func (q *QdrantClient) HealthCheck() error {
	req, err := http.NewRequest("GET", q.baseURL, nil)
	if err != nil {
		return err
	}
	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (q *QdrantClient) EnsureCollection(name string, dims uint64) error {
	url := q.collectionURL(name)
	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     dims,
			"distance": "Cosine",
		},
	}

	data, err := q.do("PUT", url, body)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	_ = data
	return nil
}

func (q *QdrantClient) DeleteCollection(name string) error {
	_, err := q.do("DELETE", q.collectionURL(name), nil)
	return err
}

func (q *QdrantClient) UpsertPoints(collection string, points []QdrantPoint) error {
	type namedPoint struct {
		ID      uint64                 `json:"id"`
		Vector  []float32              `json:"vector"`
		Payload map[string]interface{} `json:"payload"`
	}

	named := make([]namedPoint, len(points))
	for i, p := range points {
		named[i] = namedPoint{ID: p.ID, Vector: p.Vector, Payload: p.Payload}
	}

	body := map[string]interface{}{
		"points": named,
		"wait":   true,
	}
	_, err := q.do("PUT", q.pointsURL(collection)+"?wait=true", body)
	return err
}

func (q *QdrantClient) DeletePoints(collection string, ids []uint64) error {
	filter := map[string]interface{}{
		"filter": map[string]interface{}{
			"must": []map[string]interface{}{
				{
					"has_id": ids,
				},
			},
		},
	}
	_, err := q.do("POST", q.pointsURL(collection)+"/delete", filter)
	return err
}

func (q *QdrantClient) Search(collection string, vector []float32, limit int) ([]qdrantSearchResult, error) {
	body := map[string]interface{}{
		"vector": vector,
		"limit":  limit,
		"with_payload": true,
	}
	data, err := q.do("POST", q.pointsURL(collection)+"/search", body)
	if err != nil {
		return nil, err
	}

	var results []qdrantSearchResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (q *QdrantClient) DeletePointsByFilter(collection string, filter map[string]interface{}) error {
	body := map[string]interface{}{
		"filter": filter,
	}
	_, err := q.do("POST", q.pointsURL(collection)+"/delete", body)
	return err
}
