package store

import (
	"context"
	"fmt"
	"github.com/AKYC-chat/akyc-chatting/connections"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"log"
)

type MessageStorage struct {
	TableName string
}

type MessageEntity struct {
	UserId   string `dynamodbav:"user_id"`
	CreateAt string `dynamodbav:"last_post_date_time"`
	RoomId   string `dynamodbav:"room_id"`
	message  string `dynamodbav:"message_content"`
}

func (m *MessageStorage) storeMessage(e MessageEntity) error {
	connection := connections.DatabaseConnection.Conn
	item, err := attributevalue.MarshalMap(e)
	if err != nil {
		panic(err)
	}
	_, err = connection.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(m.TableName),
		Item:      item,
	})
	return err
}

func (m *MessageStorage) findByRoomId(roomId string) (MessageEntity, error) {
	var messageEntity MessageEntity
	connection := connections.DatabaseConnection.Conn
	params, err := attributevalue.MarshalList([]interface{}{roomId})
	if err != nil {
		panic(err)
	}
	response, err := connection.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE room_id=?", m.TableName),
		),
		Parameters: params,
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", roomId, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Items[0], &messageEntity)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	return messageEntity, err
}
