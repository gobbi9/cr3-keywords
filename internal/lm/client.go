package lm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// Client is an LM Studio API client used for model discovery and caption requests.
type Client struct {
	httpClient *http.Client
	baseURL    string
	logger     *slog.Logger
}

// NewClient creates a Client configured for a local LM Studio server.
func NewClient(logger *slog.Logger) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    "http://localhost:1234",
		logger:     logger,
	}
}

// DetectBestModel discovers available models and returns the best candidate
// using the internal vision-priority scoring heuristic.
func (c *Client) DetectBestModel(ctx context.Context) (string, error) {
	models, err := c.modelsFromHTTP(ctx)
	if err == nil && len(models) > 0 {
		return pickBestModel(models), nil
	}

	models, err = modelsFromLMSCLI(ctx)
	if err == nil && len(models) > 0 {
		return pickBestModel(models), nil
	}

	return "", errors.New("could not detect loaded LM Studio model via HTTP or lms CLI")
}

// PromptWithImage sends a multimodal chat completion request with prompt text and
// a base64-encoded JPEG image and returns the model response content.
func (c *Client) PromptWithImage(ctx context.Context, model string, prompt string, imageBase64 string) (string, error) {
	payload := map[string]any{
		"model": model,
		"messages": []any{
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{"type": "text", "text": prompt},
					map[string]any{
						"type": "image_url",
						"image_url": map[string]any{
							"url": "data:image/jpeg;base64," + imageBase64,
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal chat payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build chat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send chat request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read chat response body: %w", err)
	}
	if c.logger != nil {
		c.logger.Debug("LM Studio response body", "server", c.baseURL, "status", resp.StatusCode, "body", string(bodyBytes))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("lm studio chat returned status %d", resp.StatusCode)
	}

	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(bodyBytes, &out); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", errors.New("chat response had no choices")
	}

	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

func (c *Client) modelsFromHTTP(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/models", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var out struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	models := make([]string, 0, len(out.Data))
	for _, m := range out.Data {
		if strings.TrimSpace(m.ID) != "" {
			models = append(models, m.ID)
		}
	}
	return models, nil
}

func modelsFromLMSCLI(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "lms", "ls", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var rows []struct {
		ID     string `json:"id"`
		Loaded bool   `json:"loaded"`
	}
	if err := json.Unmarshal(out, &rows); err != nil {
		return nil, err
	}

	var models []string
	for _, r := range rows {
		if r.Loaded && strings.TrimSpace(r.ID) != "" {
			models = append(models, r.ID)
		}
	}
	return models, nil
}

func pickBestModel(models []string) string {
	sorted := append([]string(nil), models...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return scoreModelName(sorted[i]) > scoreModelName(sorted[j])
	})
	return sorted[0]
}

func scoreModelName(name string) int {
	n := strings.ToLower(name)
	score := 0

	if strings.Contains(n, "vl") || strings.Contains(n, "vision") {
		score += 50
	}
	if strings.Contains(n, "qwen2.5") || strings.Contains(n, "qwen") {
		score += 25
	}
	if strings.Contains(n, "gemma") {
		score += 20
	}
	if strings.Contains(n, "llava") {
		score += 15
	}
	if strings.Contains(n, "7b") || strings.Contains(n, "8b") || strings.Contains(n, "9b") {
		score += 5
	}
	return score
}
