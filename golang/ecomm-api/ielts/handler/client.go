package anthropic

import (
	"bytes"
	"encoding/json"
	"english-ai-full/ecomm-api/types"

	"io/ioutil"
	"net/http"
)

const (
	ANTHROPIC_API_KEY = "your-api-key-here"
	ANTHROPIC_API_URL = "https://api.anthropic.com/v1/completions"
)

func CallAnthropic(prompt string, maxTokens int) (string, error) {
	reqBody, err := json.Marshal(types.AnthropicRequest{
		Model:             "claude-3-sonnet-20240229",
		Prompt:            prompt,
		MaxTokensToSample: maxTokens,
		StopSequences:     []string{"\n\nHuman:"},
		Temperature:       0.7,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ANTHROPIC_API_URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", ANTHROPIC_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var anthropicResp types.AnthropicResponse
	err = json.Unmarshal(body, &anthropicResp)
	if err != nil {
		return "", err
	}

	return anthropicResp.Completion, nil
}
