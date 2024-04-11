package connector

import (
	"context"

	"github.com/AKYC-chat/akyc-chatting/message"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Connector interface {
	GetConnection() message.MessageHandler
}

func SqsGetConnection() message.MessageHandler {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("local"))
	if err != nil {
		panic(err)
	}
	client := message.SqsMessageHandler{Client: sqs.NewFromConfig(cfg)}
	return &client
}
