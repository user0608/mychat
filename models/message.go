package models

type Message struct {
	ID             int    `json:"id,omitempty"`
	SenderUserName string `json:"sender"`
	TargetUserName string `json:"target"`
	Body           string `json:"body"`
}

func NewMessage(sender, target string, body string) *Message {
	return &Message{
		SenderUserName: sender,
		TargetUserName: target,
		Body:           body,
	}
}
