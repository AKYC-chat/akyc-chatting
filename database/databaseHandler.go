package database

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DatabaseHandler interface {
	Insert()
	Delete()
	Get()
	Update()
	CreateTable(tableName string, attributes []DynamoDBTableAttribute) (*types.TableDescription, error)
}
