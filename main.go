package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	var rooms []Room
	db.Find(&rooms)
	for _, room := range rooms {
		go wsWorker(room.ID)
	}
	r := gin.Default()
	r.POST("/room", addRoomHandler)
	r.GET("/room", listRoomHandler)
	r.GET("/gifts/nonfree/:id", giftsNonFreeHandler)
	r.GET("/gifts/free/:id", giftsFreeHandler)
	r.Run()
}
