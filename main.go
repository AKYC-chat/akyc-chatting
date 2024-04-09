package main

import (
	"log"

	"github.com/AKYC-chat/akyc-chatting/connector"
)

func main() {
	messageHandler := connector.GetConnection()
	url, err := messageHandler.CreateQueue("test1", true)
	if err != nil {
		log.Fatalf("CreateQueue err: %v", err)
	}
	log.Println(url)
}
