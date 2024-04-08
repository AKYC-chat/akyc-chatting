package messageconnector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"strings"
	"time"
)

const FIFO_SUFFIX = ".fifo"

type SqsClient struct {
	client *sqs.Client
}

func (s SqsClient) SendMessage(messageBody string, messageUrl string, groupId string) {

	date := strings.Join(strings.Split(time.Now().Format(time.DateTime), " "), "/")
	log.Printf("message duplicationId: %v \n", date)
	_, err := s.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody:            aws.String(messageBody),
		MessageGroupId:         aws.String(groupId),
		QueueUrl:               aws.String(messageUrl),
		MessageDeduplicationId: aws.String(date),
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func (s SqsClient) ReceiveMessage() error {
	client := s.client
	sqsUrl := "url"
	_, err := client.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(sqsUrl),

		MaxNumberOfMessages: 1,
		VisibilityTimeout:   30,
	})
	if err != nil {
		panic(err)
	}
	return err
}

func (s SqsClient) DeleteMessage() {

}

func (s SqsClient) CreateQueue(queueName string, isFifoQueue bool) (url string, err error) {
	var queueUrl string
	queueAttributes := map[string]string{}

	if isFifoQueue {
		queueAttributes = map[string]string{
			"FifoQueue":                 "true",
			"ContentBasedDeduplication": "true",
			"DeduplicationScope":        "messageGroup",
			"FifoThroughputLimit":       "perMessageGroupId",
		}
	}

	queue, err := s.client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName:  aws.String(queueName + FIFO_SUFFIX),
		Attributes: queueAttributes,
	})
	if err != nil {
		log.Fatalf("Couldn't create queue %v. caused by: %v\n", queueName, err)
		return "", err
	}
	queueUrl = *queue.QueueUrl
	return queueUrl, err
}
func (s SqsClient) GetQueueList() (queueUrls []string, err error) {
	paginator := sqs.NewListQueuesPaginator(s.client, &sqs.ListQueuesInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("Couldn't get queues. Here's why: %v\n", err)
			return queueUrls, err
		} else {
			queueUrls = append(queueUrls, output.QueueUrls...)
		}
	}
	if len(queueUrls) == 0 {
		log.Println("empty queue")
	}
	return queueUrls, err
}
