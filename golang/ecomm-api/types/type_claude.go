package types


type AnthropicRequest struct {
	Model               string   `json:"model"`
	Prompt              string   `json:"prompt"`
	MaxTokensToSample   int      `json:"max_tokens_to_sample"`
	StopSequences       []string `json:"stop_sequences"`
	Temperature         float64  `json:"temperature"`
}

type AnthropicResponse struct {
	Completion string `json:"completion"`
}