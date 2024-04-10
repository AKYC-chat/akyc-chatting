package message

import "time"

type MessageHandler interface {
	SendMessage(messageBody string, messageUrl string, groupId string)

	ReceiveMessage(messageUrl string) ([]Message, error)
	DeleteMessage()

	CreateQueue(queueName string, isFifoQueue bool) (url string, err error)
	GetQueueList() (queueUrls []string, err error)
}

type Message struct {
	body              *string
	groupId           *string
	sentTimestamp     time.Time
	receivedTimestamp time.Time
}
