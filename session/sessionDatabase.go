package session

import (
	"context"
	"fmt"
	"log"

	"github.com/AKYC-chat/akyc-chatting/connections"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	sessionTableName = "AYKC_SESSION"
)

type SessionDatabase struct{}

type SessionEntity struct {
	UserId    string `dynamodbav:"user_id"`
	CreateAt  string `dynamodbav:"create_at"`
	SessionId string `dynamodbav:"session_id"`
}

func (s SessionEntity) GetKey() map[string]types.AttributeValue {
	userId, err := attributevalue.Marshal(s.UserId)
	if err != nil {
		panic(err)
	}

	createAt, err := attributevalue.Marshal(s.CreateAt)
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{"user_id": userId, "create_at": createAt}
}

func (s *SessionDatabase) CreateSession(e SessionEntity) error {
	item, err := attributevalue.MarshalMap(e)

	if err != nil {
		panic(err)
	}
	_, err = connections.DatabaseConnection.Conn.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(sessionTableName), Item: item,
	})

	return err
}

func (s *SessionDatabase) DeleteSession(e SessionEntity) error {
	_, err := connections.DatabaseConnection.Conn.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(sessionTableName), Key: e.GetKey(),
	})
	if err != nil {
		log.Printf("Couldn't delete %v from the table. Here's why: %v\n", e.UserId, err)
	}
	return err
}

func (sessionDatabase *SessionDatabase) GetSession(sessionEntity SessionEntity) (SessionEntity, error) {
	var responseSessionEntity SessionEntity
	response, err := connections.DatabaseConnection.Conn.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: sessionEntity.GetKey(), TableName: aws.String(sessionTableName),
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", sessionEntity.UserId, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &responseSessionEntity)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	return responseSessionEntity, err
}

func (sessionDatabase *SessionDatabase) GetSessionByUserId(userId string) (SessionEntity, error) {
	var sessionEntity SessionEntity

	params, err := attributevalue.MarshalList([]interface{}{userId})
	if err != nil {
		panic(err)
	}

	response, err := connections.DatabaseConnection.Conn.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(
			fmt.Sprintf("SELECT * FROM \"%v\" WHERE user_id=?", sessionTableName),
		),
		Parameters: params,
	})

	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", userId, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Items[0], &sessionEntity)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	return sessionEntity, err
}

func (sessionDatabase *SessionDatabase) GetAllSessions() ([]SessionEntity, error) {
	var sessions []SessionEntity
	var err error
	var response *dynamodb.ScanOutput

	projEx := expression.NamesList(
		expression.Name("user_id"), expression.Name("session_id"), expression.Name("create_at"))
	expr, err := expression.NewBuilder().WithProjection(projEx).Build()
	if err != nil {
		log.Printf("Couldn't build expressions for scan. Here's why: %v\n", err)
	} else {
		scanPaginator := dynamodb.NewScanPaginator(connections.DatabaseConnection.Conn, &dynamodb.ScanInput{
			TableName:                 aws.String(sessionTableName),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
		})
		for scanPaginator.HasMorePages() {
			response, err = scanPaginator.NextPage(context.TODO())
			if err != nil {
				log.Printf("Couldn't scan for sessions. Here's why: %v\n", err)
				break
			} else {
				var sessionsPage []SessionEntity
				err = attributevalue.UnmarshalListOfMaps(response.Items, &sessionsPage)
				if err != nil {
					log.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
					break
				} else {
					sessions = append(sessions, sessionsPage...)
				}
			}
		}
	}
	return sessions, err
}
