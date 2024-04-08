package messageconnector

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsConnection struct {
}

func (s SqsConnection) GetConnection() Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := SqsClient{client: sqs.NewFromConfig(cfg)}
	return &client
}
