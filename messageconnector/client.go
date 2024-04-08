package messageconnector

import "github.com/aws/aws-sdk-go-v2/service/sqs"

type Client interface {
	SendMessage(messageBody string)
	ReceiveMessage() (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(*sqs.SendMessageOutput)
}
