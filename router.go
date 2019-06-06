package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
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

func giftsNonFreeHandler(c *gin.Context) {
	var r []int64
	var start time.Time
	var end time.Time
	id := c.Param("id")
	startString := c.Query("start")
	endString := c.Query("end")
	if startString == "" {
		start = time.Now().Add(-time.Hour * 1)
	} else {
		startInt, err := strconv.ParseInt(startString, 10, 64)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		start = time.Unix(startInt, 0)
	}
	if startString == "" {
		end = time.Now()
	} else {
		endInt, err := strconv.ParseInt(endString, 10, 64)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		end = time.Unix(endInt, 0)
	}
	db.Table("gifts").Select("sum(price) as n").Where("up = ? AND type = 1 AND created_at > ? AND created_at < ?", id, start, end).Pluck("n", &r)
	c.JSON(200, r[0])
}

func giftsFreeHandler(c *gin.Context) {
	var r []int64
	var start time.Time
	var end time.Time
	id := c.Param("id")
	startString := c.Query("start")
	endString := c.Query("end")
	if startString == "" {
		start = time.Now().Add(-time.Hour * 1)
	} else {
		startInt, err := strconv.ParseInt(startString, 10, 64)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		start = time.Unix(startInt, 0)
	}
	if startString == "" {
		end = time.Now()
	} else {
		endInt, err := strconv.ParseInt(endString, 10, 64)
		if err != nil {
			c.JSON(500, gin.H{})
			return
		}
		end = time.Unix(endInt, 0)
	}
	db.Table("gifts").Select("sum(price) as n").Where("up = ? AND type = 0 AND created_at > ? AND created_at < ?", id, start, end).Pluck("n", &r)
	c.JSON(200, r[0])
}
