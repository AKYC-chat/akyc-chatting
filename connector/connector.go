package connector

import (
	"context"

	"github.com/AKYC-chat/akyc-chatting/database"
	"github.com/AKYC-chat/akyc-chatting/message"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SqsGetConnection() message.MessageHandler {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("local"))
	if err != nil {
		panic(err)
	}
	client := message.SqsMessageHandler{Client: sqs.NewFromConfig(cfg)}
	return &client
}

func DynamoDBGetConnection() database.DatabaseHandler {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("local"))
	if err != nil {
		panic(err)
	}

	client := database.DynamoDBHandler{Client: dynamodb.NewFromConfig(cfg)}

	return &client
}
