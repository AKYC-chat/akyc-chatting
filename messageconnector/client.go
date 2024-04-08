package messageconnector

type Client interface {
	SendMessage(messageBody string, messageUrl string, groupId string)

	ReceiveMessage() error
	DeleteMessage()

	CreateQueue(queueName string, isFifoQueue bool) (url string, err error)
	GetQueueList() (queueUrls []string, err error)
}
