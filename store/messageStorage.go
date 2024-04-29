package store

import (
	"context"
	"fmt"
	"github.com/AKYC-chat/akyc-chatting/connections"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
	Message  string `dynamodbav:"message_content"`
}

func (m *MessageStorage) InsertMessage(e MessageEntity) error {
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

func (m *MessageStorage) FindByRoomId(roomId string) ([]MessageEntity, error) {
	var messageEntities []MessageEntity
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
	//log.Printf("reponse: %v ", response.Items)
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", roomId, err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &messageEntities)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	return messageEntities, err
}

func (m *MessageStorage) FindByUserId(userId string) ([]MessageEntity, error) {
	var messageEntities []MessageEntity
	connection := connections.DatabaseConnection.Conn
	params, err := attributevalue.MarshalList([]interface{}{userId})
	if err != nil {
		panic(err)
	}
	response, err := connection.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE user_id=?", m.TableName),
		),
		Parameters: params,
	})
	//log.Printf("reponse: %v ", response.Items)
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", userId, err)
	} else {
		err = attributevalue.UnmarshalListOfMaps(response.Items, &messageEntities)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	return messageEntities, err
}

func (m *MessageStorage) FindAllRooms() ([]MessageEntity, error) {
	var messageEntities []MessageEntity

	projEx := expression.NamesList(
		expression.Name("user_id"),
		expression.Name("room_id"),
		expression.Name("message_content"),
		expression.Name("last_post_date_time"))
	expr, err := expression.NewBuilder().WithProjection(projEx).Build()

	if err != nil {
		panic(err)
	}
	scanPaginator := dynamodb.NewScanPaginator(connections.DatabaseConnection.Conn, &dynamodb.ScanInput{
		TableName:                 aws.String(m.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	})

	for scanPaginator.HasMorePages() {
		response, err := scanPaginator.NextPage(context.TODO())
		if err != nil {
			log.Printf("Couldn't scan for sessions. Here's why: %v\n", err)
			break
		} else {
			var messagePage []MessageEntity
			err = attributevalue.UnmarshalListOfMaps(response.Items, &messagePage)
			if err != nil {
				log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
				break
			} else {
				messageEntities = append(messageEntities, messagePage...)
			}
		}
	}

	return messageEntities, err
}
