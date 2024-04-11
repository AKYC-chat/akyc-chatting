package message

import (
	"context"
	"github.com/AKYC-chat/akyc-chatting/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
)

const FifoSuffix = ".fifo"

type SqsMessageHandler struct {
	Client *sqs.Client
}

func (s SqsMessageHandler) SendMessage(messageBody string, messageUrl string, groupId string) (messageId *string, err error) {
	messageOutPut, err := s.Client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody:            aws.String(messageBody),
		MessageGroupId:         aws.String(groupId),
		QueueUrl:               aws.String(messageUrl),
		MessageDeduplicationId: aws.String(util.MessageDateTime()),
	})
	if err != nil {
		log.Fatalln(err)

		return nil, err
	}
	return messageOutPut.MessageId, nil

}

func (s SqsMessageHandler) ReceiveMessage(messageUrl string) ([]Message, error) {
	client := *s.Client

	// 메세지 가져오기
	messageOutput, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(messageUrl),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
		MessageAttributeNames: []string{
			"All",
		},
		MaxNumberOfMessages: 10,
		VisibilityTimeout:   30,
	})
	if err != nil {
		panic(err)
	}

	// sqs struct -> messageList 변환
	messages := &messageOutput.Messages
	messageList := make([]Message, 0)
	for _, message := range *messages {
		messageBody := aws.ToString(message.Body)
		groupId := message.Attributes["MessageGroupId"]
		receiveTime := util.ParseTimestamp(message.Attributes["ApproximateFirstReceiveTimestamp"])

		sentTime := util.ParseTimestamp(message.Attributes["SentTimestamp"])

		messageList = append(
			messageList,
			Message{body: &messageBody, groupId: &groupId, receivedTimestamp: receiveTime, sentTimestamp: sentTime},
		)
	}

	return messageList, err
}

func (s SqsMessageHandler) DeleteMessage() {

}

func (s SqsMessageHandler) CreateQueue(queueName string, isFifoQueue bool) (url string, err error) {
	var queueUrl string
	queueAttributes := map[string]string{}

	if isFifoQueue {
		queueAttributes = map[string]string{
			"FifoQueue":                 "true",
			"ContentBasedDeduplication": "true",
			"DeduplicationScope":        "messageGroup",
			"FifoThroughputLimit":       "perMessageGroupId",
		}
		queue, err := s.Client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
			QueueName:  aws.String(queueName + FifoSuffix),
			Attributes: queueAttributes,
		})
		if err != nil {
			log.Fatalf("Couldn't create queue %v. caused by: %v\n", queueName, err)
			return "", err
		}
		queueUrl = *queue.QueueUrl

		return queueUrl, err
	}
	queue, err := s.Client.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
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
func (s SqsMessageHandler) GetQueueList() (queueUrls []string, err error) {
	paginator := sqs.NewListQueuesPaginator(s.Client, &sqs.ListQueuesInput{
		QueueNamePrefix: aws.String("Aykc-Chat"),
	})
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

func deleteMessage(client *sqs.Client) {
	client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{})
}
