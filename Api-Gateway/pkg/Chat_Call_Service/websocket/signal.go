package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// AllRooms is the global hashmap for the server
var AllRooms RoomMap

// CreateRoomRequestHandler creates a Room and returns roomID
func CreateRoomRequestHandler(c *gin.Context) {
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)
	c.JSON(200, resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var (
	broadcast = make(chan broadcastMsg)
	clients   = make(map[*websocket.Conn]bool)
	mutex     sync.Mutex
)

func broadcaster() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					log.Println("Write error:", err)
					client.Conn.Close()
					delete(clients, client.Conn)
				}
			}
		}
		mutex.Unlock()
	}
}

func init() {
	go broadcaster()
}

// JoinRoomRequestHandler will join the client in a particular room
func JoinRoomRequestHandler(c *gin.Context) {
	roomID := c.Query("roomID")
	if roomID == "" {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Web Socket Upgrade Error:", err)
		return
	}

	mutex.Lock()
	AllRooms.InsertIntoRoom(roomID, false, ws)
	clients[ws] = true
	mutex.Unlock()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Println("Read Error:", err)
			break
		}

		msg.Client = ws
		msg.RoomID = roomID

		log.Println(msg.Message)

		broadcast <- msg
	}

	mutex.Lock()
	delete(clients, ws)
	mutex.Unlock()
	ws.Close()
}