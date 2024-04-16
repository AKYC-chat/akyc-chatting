package database

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DatabaseHandler interface {
	Insert(tableName string, Object interface{}) error
	Delete()
	Get()
	Update()
	CreateTable(tableName string, attributes []DynamoDBTableAttribute) (*types.TableDescription, error)
}
