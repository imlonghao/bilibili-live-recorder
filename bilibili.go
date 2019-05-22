package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsHeartbeatSent  = 2
	wsHeartbeatReply = 3
	wsMessage        = 5
	wsAuth           = 7
	wsAuthSuccess    = 8
)

type senderHeader struct {
	A uint32
	B uint16
	C uint16
	D uint32
	E uint32
}

type recvHeader struct {
	Length uint32
	Two    uint32
	Type   uint32
	Four   uint32
}

type renqizhi struct {
	Value uint32
}

type cmdJSON struct {
	Info [][]interface{} `json:"info"`
	Cmd  string          `json:"cmd"`
}

type danmakuJSON struct {
	Info []string `json:"info"`
}

type giftJSON struct {
	Data struct {
		GiftName  string `json:"giftName"`
		Num       uint   `json:"num"`
		Uname     string `json:"uname"`
		UID       uint   `json:"uid"`
		Remain    uint   `json:"remain"`
		CoinType  string `json:"coin_type"`
		TotalCoin uint   `json:"total_coin"`
	} `json:"data"`
}

type bilibiliClient struct {
	room uint
	ws   *websocket.Conn
}

func newBilibiliClient(room uint) (*bilibiliClient, error) {
	b := bilibiliClient{
		room: room,
	}
	c, _, err := websocket.DefaultDialer.Dial("wss://broadcastlv.chat.bilibili.com/sub", nil)
	if err != nil {
		return nil, err
	}
	b.ws = c
	return &b, nil
}

func (c *bilibiliClient) send(data []byte, operation uint32) {
	h := senderHeader{uint32(binary.Size(data) + 16), 16, 1, operation, 1}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, h)
	c.ws.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprintf("%s%s", buf, data)))
}

func (c *bilibiliClient) messageWorker(message []byte) {
	h := message[:16]
	var head recvHeader
	buf := bytes.NewReader(h)
	binary.Read(buf, binary.BigEndian, &head)
	body := message[16:head.Length]
	switch head.Type {
	case wsHeartbeatReply:
		var rqz renqizhi
		binary.Read(bytes.NewReader(body), binary.BigEndian, &rqz)
		if rqz.Value > 1 {
			go db.Create(&Rqz{
				UP:    c.room,
				Value: uint(rqz.Value),
			})
		}
	case wsMessage:
		var cmd cmdJSON
		_ = json.Unmarshal(body, &cmd)
		if cmd.Cmd != "DANMU_MSG" {
			go sendLog(body)
		}
		switch cmd.Cmd {
		case "DANMU_MSG":
			var danmaku danmakuJSON
			json.Unmarshal(body, &danmaku)
			log.Printf("%d - %s: %s\n", c.room, cmd.Info[2][1], danmaku.Info[1])
			go sendLog([]byte(fmt.Sprintf("{\"CMD\": \"DANMU_MSG\", \"kimo\": %t, \"roomid\": %d, \"user\": %d, \"message\": \"%s\"}", cmd.Info[0][9].(uint) == 1, c.room, cmd.Info[2][0].(uint), danmaku.Info[1])))
			user := User{
				ID:   cmd.Info[2][0].(uint),
				Name: cmd.Info[2][1].(string),
			}
			go db.Save(&user)
			d := Danmaku{
				UP:        c.room,
				UserRefer: user.ID,
				Kimo:      cmd.Info[2][0].(uint) == 1,
				Comment:   danmaku.Info[1],
			}
			go db.Create(&d)
		case "SEND_GIFT":
			var g giftJSON
			json.Unmarshal(body, &g)
			log.Printf("%d - %s: %s (%s) x %d\n", c.room, g.Data.Uname, g.Data.GiftName, g.Data.CoinType, g.Data.Num)
			user := User{
				ID:   g.Data.UID,
				Name: g.Data.Uname,
			}
			go db.Save(&user)
			gift := Gift{
				UP:        c.room,
				UserRefer: g.Data.UID,
				Type:      g.Data.CoinType == "gold",
				Name:      g.Data.GiftName,
				Number:    g.Data.Num,
				Price:     g.Data.TotalCoin,
				Remain:    g.Data.Remain,
			}
			go db.Create(&gift)
		}
	}
	next := message[head.Length:]
	if binary.Size(next) != 0 {
		c.messageWorker(next)
	}
}

func wsWorker(room uint) {
	log.Printf("adding %d to monitoring list", room)
	c, err := newBilibiliClient(room)
	if err != nil {
		log.Printf("%d - fail to create a new bilibili client, err: %s", room, err)
		time.Sleep(5 * time.Second)
		go wsWorker(room)
		return
	}
	auth := []byte(fmt.Sprintf("{\"uid\":0,\"roomid\":%d,\"protover\":1,\"platform\":\"web\",\"clientver\":\"1.4.0\"}", room))
	c.send(auth, 7)
	go func() {
		for {
			c.send([]byte(""), 2)
			time.Sleep(30 * time.Second)
		}
	}()
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("%d - fail to connect to server, restarting, err: %s", room, err)
			c.ws.Close()
			time.Sleep(5 * time.Second)
			go wsWorker(room)
			return
		}
		go c.messageWorker(message)
	}
}
