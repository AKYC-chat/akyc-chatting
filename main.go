package main

import (
	"github.com/AKYC-chat/akyc-chatting/messageconnector"
	"log"
)

func main() {

	//sqsConnection := messageconnector.SqsConnection{}
	//client := sqsConnection.GetConnection()
	//client.SendMessage("Hello World")
	client := messageconnector.GetConnection()
	url, err := client.CreateQueue("test-1", true)
	if err != nil {
		log.Fatalf("CreateQueue err: %v", err)
	}
	log.Println(url)
}
