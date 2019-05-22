package main

import (
	"github.com/dougEfresh/logzio-go"
	"log"
	"os"
)

var logz *logzio.LogzioSender

func sendLog(data []byte) {
	err := logz.Send(data)
	if err != nil {
		log.Printf("fail to send log to logz, err: %s", err)
	}
}

func init() {
	var err error
	logz, err = logzio.New(os.Getenv("LOGZTOKEN"))
	if err != nil {
		panic(err)
	}
}
