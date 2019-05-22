package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func addRoomHandler(c *gin.Context) {
	room := c.PostForm("room")
	name := c.PostForm("name")
	var r Room
	roomint, err := strconv.Atoi(room)
	if err != nil {
		c.JSON(500, gin.H{
			"room": room,
			"name": name,
		})
		return
	}
	db.Where("id = ?", roomint).Take(&r)
	if r.Name == "" {
		go db.Save(&Room{
			ID:   uint(roomint),
			Name: name,
		})
		go wsWorker(uint(roomint))
	}
	c.JSON(200, gin.H{
		"room": room,
		"name": name,
	})
}

func listRoomHandler(c *gin.Context) {
	var rooms []Room
	db.Find(&rooms)
	c.JSON(200, rooms)
}
