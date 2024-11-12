package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Videocallroutes() {
	// Create a new Gin router instance
	r := gin.New()

	// Attach Logger and Recovery middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default()) // CORS middleware to handle cross-origin requests

	// Load HTML templates from the specified directory
	r.LoadHTMLGlob("pkg/template/*")

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/room/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		c.HTML(http.StatusOK, "room.html", gin.H{"RoomID": roomID})
	})

	// POST route to create a room
	r.POST("/createRoom", CreateRoomRequestHandler)

	// GET route to join a room
	r.GET("/join", JoinRoomRequestHandler)
}
