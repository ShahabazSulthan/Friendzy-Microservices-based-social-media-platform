package requestmodel

import "time"

type OneToOneChatRequest struct {
	SenderID    string
	RecipientID string
	Content     string
	TimeStamp   time.Time
	Status      string
}
