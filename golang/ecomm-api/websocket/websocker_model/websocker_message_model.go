package websocket_model

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Sender  string `json:"sender"`
}