package message

import "time"

type MessageHandler interface {
	SendMessage(messageBody string, messageUrl string, groupId string) (messageId *string, err error)

	ReceiveMessage(messageUrl string) ([]Message, error)
	DeleteMessage()

	CreateQueue(queueName string, isFifoQueue bool) (url string, err error)
	GetQueueList() (queueUrls []string, err error)
}

type Message struct {
	Body              string
	GroupId           string
	SentTimestamp     *time.Time
	ReceivedTimestamp *time.Time
}
