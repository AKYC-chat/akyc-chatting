package messageconnector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsClient struct {
	client *sqs.Client
}

func (s SqsClient) SendMessage(messageBody string) {

}

func (s SqsClient) ReceiveMessage() (*sqs.ReceiveMessageOutput, error) {
	client := s.client
	sqsUrl := "url"
	receiveResult, err := client.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(sqsUrl),

		MaxNumberOfMessages: 1,
		VisibilityTimeout:   30,
	})
	if err != nil {
		panic(err)
	}
	return receiveResult, err
}

func (s SqsClient) DeleteMessage(message *sqs.SendMessageOutput) {}
