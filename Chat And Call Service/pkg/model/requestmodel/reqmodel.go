package requestmodel

import (
	"time"

	"github.com/gorilla/websocket"
)

type OneToOneChatRequest struct {
	SenderID    string
	RecipientID string
	Content     string
	TimeStamp   time.Time
	Status      string
}

type NewGroupInfo struct {
	GroupName    string
	GroupMembers []uint64
	CreatorID    string
	CreateAt     time.Time
}

type OneToManyMessageRequest struct {
	SenderID  string
	GroupID   string
	Content   string
	TimeStamp time.Time
	Status    string
}

type AddNewMemberToGroup struct {
	UserID       string
	GroupID      string
	GroupMembers []uint64
}

type RemoveMemberFromGroup struct {
	UserID   string
	GroupID  string
	MemberID string
}

// Participant represents a user in a chat room.
type Participant struct {
	Host bool            `bson:"host"`
	Conn *websocket.Conn `bson:"-"`
}

// Room represents a chat room with a unique ID and participants.
type Room struct {
	ID           string        `bson:"_id,omitempty"`
	Participants []Participant `bson:"participants"`
}
