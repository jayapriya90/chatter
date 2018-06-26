package main

import (
	"flag"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jayapriya90/chatter/backend"
	log "github.com/sirupsen/logrus"
)

var port = flag.Int("port", 8888, "port to run the chatter service")

func GetChatRooms(c *gin.Context) {

}
func CreateChatRoom(c *gin.Context) {

}
func DeleteChatRoom(c *gin.Context) {

}
func GetUsersInChatRoom(c *gin.Context) {

}
func JoinChatRoom(c *gin.Context) {

}
func LeaveChatRoom(c *gin.Context) {

}
func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetLevel(log.DebugLevel)
	flag.Parse()
	server := backend.NewServer()
	go server.Run()
	router := gin.Default()
	group := router.Group("/v1")
	{
		group.GET("/websocket", func(c *gin.Context) {
			backend.ServeWebSocket(server, c.Writer, c.Request)
		})
		group.GET("/chatroom", GetChatRooms)
		group.POST("/chatroom", CreateChatRoom)
		group.GET("/chatroom/:id", GetUsersInChatRoom)
		group.POST("/chatroom/join", JoinChatRoom)
		group.POST("/chatroom/leave", LeaveChatRoom)
		group.DELETE("/chatroom", DeleteChatRoom)
	}
	log.Infof("Listening and serving on %d", *port)
	portStr := strconv.Itoa(*port)
	router.Run(":" + portStr)
}
