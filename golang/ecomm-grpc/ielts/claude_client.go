package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"english-ai-full/ecomm-api/types"
)

const (
	ANTHROPIC_API_URL = "https://api.anthropic.com/v1/completions"
)

type Service struct {
	apiKey string
	client *http.Client
}

func NewService(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (s *Service) CallAnthropic(ctx context.Context, prompt string, maxTokens int) (string, error) {
	reqBody, err := json.Marshal(types.AnthropicRequest{
		Model:             "claude-3-sonnet-20240229",
		Prompt:            prompt,
		MaxTokensToSample: maxTokens,
		StopSequences:     []string{"\n\nHuman:"},
		Temperature:       0.7,
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ANTHROPIC_API_URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var anthropicResp types.AnthropicResponse
	err = json.Unmarshal(body, &anthropicResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	return anthropicResp.Completion, nil
}