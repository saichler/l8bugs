package triage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultModel  = "claude-sonnet-4-20250514"
	defaultURL    = "https://api.anthropic.com/v1/messages"
	apiVersion    = "2023-06-01"
	maxTokens     = 2048
	clientTimeout = 30 * time.Second
)

type Client struct {
	apiKey  string
	model   string
	baseURL string
	http    *http.Client
}

type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messagesResponse struct {
	Content []contentBlock `json:"content"`
	Error   *apiError      `json:"error,omitempty"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type apiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewClient() *Client {
	apiKey := os.Getenv("L8BUGS_ANTHROPIC_API_KEY")
	model := os.Getenv("L8BUGS_AI_MODEL")
	if model == "" {
		model = defaultModel
	}
	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: defaultURL,
		http:    &http.Client{Timeout: clientTimeout},
	}
}

func (c *Client) Available() bool {
	return c.apiKey != ""
}

func (c *Client) Complete(systemPrompt, userMessage string) (string, error) {
	if !c.Available() {
		return "", fmt.Errorf("anthropic API key not configured")
	}

	reqBody := messagesRequest{
		Model:     c.model,
		MaxTokens: maxTokens,
		System:    systemPrompt,
		Messages: []message{
			{Role: "user", Content: userMessage},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var msgResp messagesResponse
	if err := json.Unmarshal(body, &msgResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if msgResp.Error != nil {
		return "", fmt.Errorf("API error: %s", msgResp.Error.Message)
	}

	for _, block := range msgResp.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}

	return "", fmt.Errorf("no text content in response")
}
