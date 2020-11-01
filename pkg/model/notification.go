package model

type Notification struct {
	ID         string   `json:"id,omitempty"`
	Sender     string   `json:"sender,omitempty"`
	To         []string `json:"to,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	Message    string   `json:"message,omitempty"`
	SentStatus bool     `json:"sent_status,omitempty"`
	CreatedAt  string   `json:"created_at,omitempty"`
}

//easyjson:json
type NotificationList []Notification
