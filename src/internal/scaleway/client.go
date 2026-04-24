package scaleway

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	model              = "mistral-small-3.2-24b-instruct-2506"
	defaultMaxTokens   = 300
	defaultTemperature = 0.92
)

type Client struct {
	apiURL string
	apiKey string
	http   *http.Client
}

func NewClient(apiURL, apiKey string) *Client {
	return &Client{
		apiURL: apiURL,
		apiKey: apiKey,
		http:   &http.Client{},
	}
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionRequest struct {
	Model       string    `json:"model"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
	Messages    []message `json:"messages"`
}

// Stream calls the Scaleway SSE endpoint and writes content chunks to w as they arrive.
// It returns the full accumulated response for caching purposes.
func (c *Client) Stream(ctx context.Context, systemPrompt, userPrompt string, w io.Writer) (string, error) {
	body, err := json.Marshal(completionRequest{
		Model:       model,
		MaxTokens:   defaultMaxTokens,
		Temperature: defaultTemperature,
		Stream:      true,
		Messages: []message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("scaleway: unexpected status %d", resp.StatusCode)
	}

	var full strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil || len(event.Choices) == 0 {
			continue
		}

		chunk := event.Choices[0].Delta.Content
		if chunk == "" {
			continue
		}

		full.WriteString(chunk)
		_, _ = io.WriteString(w, chunk)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	return full.String(), scanner.Err()
}
