package requestmodel

import "time"

type KafkaNotification struct {
	UserID      string
	ActorID     string
	ActionType  string
	TargetID    string
	TargetType  string
	CommentText string
	CreatedAt   time.Time
}