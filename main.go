package main

import (
	"fmt"
	"log"

	"github.com/AKYC-chat/akyc-chatting/connector"
)

func main() {
	messageHandler := connector.GetConnection()
	queueUrls, err := messageHandler.GetQueueList()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(queueUrls)
}
