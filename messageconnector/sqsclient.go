package messageconnector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
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

func (s SqsClient) CreateQueue(queueName string, isFifoQueue bool) (url string, err error) {
	var queueUrl string
	queueAttributes := map[string]string{}
	if isFifoQueue {
		queueAttributes["FifoQueue"] = "true"
	}

	queue, err := s.client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName:  aws.String(queueName),
		Attributes: queueAttributes,
	})
	if err != nil {
		log.Fatalf("Couldn't create queue %v. caused by: %v\n", queueName, err)
		return "", err
	}
	queueUrl = *queue.QueueUrl
	return queueUrl, err
}
