package requestmodel_chat

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MessageRequest struct {
	SenderID    string    `json:"SenderID" validate:"required"`
	RecipientID string    `json:"RecipientID"`
	Content     string    `json:"Content" `
	TimeStamp   time.Time `json:"TimeStamp"`
	Type        string    `json:"Type" validate:"required"`
	Status      string    `json:"Status"`
	GroupID     string    `json:"GroupID"`
	TypingStat  bool
}

type OnetoOneMessageRequest struct {
	SenderID    string    `json:"SenderID" validate:"required"`
	RecipientID string    `json:"RecipientID" `
	Type        string    `json:"Type"`
	Content     string    `json:"Content" validate:"required"`
	TimeStamp   time.Time `json:"TimeStamp"`
	Status      string    `json:"Status"`
}

type TypingStatusRequest struct {
	SenderID    string `json:"SenderID" `
	RecipientID string `json:"RecipientID"`
	Type        string `json:"Type" `
	TypingStat  bool
}

type OnetoManyMessageRequest struct {
	SenderID  string    `json:"SenderID" validate:"required"`
	GroupID   string    `json:"GroupID" validate:"required"`
	Type      string    `json:"Type"`
	Content   string    `json:"Content" validate:"required"`
	TimeStamp time.Time `json:"TimeStamp"`
	Status    string    `json:"Status"`
}

type NewGroupInfo struct {
	GroupName    string   `json:"GroupName" validate:"required,lte=20"`
	GroupMembers []uint64 `json:"GroupMembers" validate:"required,min=1,max=12,unique,dive,number"`
	CreatorID    string
	CreatedAt    time.Time
}

type AddNewMembersToGroup struct {
	GroupID      string   `json:"GroupID" validate:"required,min=15,lte=35"`
	GroupMembers []uint64 `json:"GroupMembers" validate:"required,min=1,max=12,unique,dive,number"`
}

type RemoveMemberFromGroup struct {
	GroupID  string `json:"GroupID" validate:"required,min=15,lte=35"`
	MemberID string `json:"MemberID" validate:"required,number"`
}

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

// RoomRepo is the main struct implementing IRoomRepo.
type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}
