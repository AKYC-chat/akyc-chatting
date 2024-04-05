package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// Connection
type MessageQueueConnection struct {
}

type MessageQueueConnector interface {
	MessageQueueConnect(configuration interface{}) MessageQueueHandler
}

type SqsMessageQueueConnector struct {
}

func (conn *SqsMessageQueueConnector) MessageQueueConnect(configuration SqsConfiguration) SqsMessageQueueHandler {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(configuration.profile))
	if err != nil {
		log.Fatal(err)
	}
	client := sqs.NewFromConfig(cfg)

	return SqsMessageQueueHandler{Client: client}
}

// Config
type SqsConfiguration struct {
	profile string
}

// Handler
type MessageQueueHandler interface {
	Send()
	Recv()
	Delete()
}

type SqsMessageQueueHandler struct {
	Client *sqs.Client
}

func (handler *SqsMessageQueueHandler) Send() {

}

func (handler *SqsMessageQueueHandler) Recv() {

}

func (handler *SqsMessageQueueHandler) Delete() {

}

func (handler *SqsMessageQueueHandler) GetQueueUrls() []string {
	var queueUrls []string
	paginator := sqs.NewListQueuesPaginator(handler.Client, &sqs.ListQueuesInput{})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Printf("Couldn't get queues. Here's why: %v\n", err)
			break
		} else {
			queueUrls = append(queueUrls, output.QueueUrls...)
		}
	}
	if len(queueUrls) == 0 {
		fmt.Println("You don't have any queues!")
	} else {
		for _, queueUrl := range queueUrls {
			fmt.Printf("\t%v\n", queueUrl)
		}
	}

	return queueUrls
}

func Run(connector MessageQueueConnector) {
	handler := connector.MessageQueueConnect(SqsConfiguration{profile: "local"})

	handler.Send()
	handler.Recv()
}

func main() {

}
