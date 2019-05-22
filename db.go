package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// Room 需要监控的直播间 ID
type Room struct {
	ID   uint
	Name string
}

// User 用户
type User struct {
	ID   uint
	Name string `gorm:"size=16"`
}

// Rqz 人气
type Rqz struct {
	ID        uint
	UP        uint `gorm:"index:up"`
	Value     uint
	CreatedAt time.Time `gorm:"index:time"`
}

// Danmaku 弹幕
type Danmaku struct {
	ID        uint
	UP        uint `gorm:"index:up"`
	User      User `gorm:"foreignkey:UserRefer"`
	UserRefer uint `gorm:"index:user"`
	Kimo      bool
	Comment   string
	CreatedAt time.Time `gorm:"index:time"`
}

// Gift 礼物
type Gift struct {
	ID        uint
	UP        uint `gorm:"index:up"`
	User      User `gorm:"foreignkey:UserRefer"`
	UserRefer uint `gorm:"index:user"`
	Type      bool // true 花钱买的 - false 免费的和白嫖的
	Name      string
	Number    uint
	Price     uint
	Remain    uint
	CreatedAt time.Time `gorm:"index:time"`
}

func init() {
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2fShanghai", os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBNAME")))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Room{}, &User{}, &Rqz{}, &Danmaku{}, &Gift{})
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}
