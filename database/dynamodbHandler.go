package database

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type DynamoDBAttrType string

const (
	DynamoDBKeyTypePart DynamoDBAttrType = "PARTITION"
	DynamoDBKeyTypeSort DynamoDBAttrType = "SORT"
)

type DynamoDBHandler struct {
	Client *dynamodb.Client
}

type DynamoDBTableAttribute struct {
	AttributeName string
	// AttributeType
	// The data type for the attribute, where:
	//   - S - the attribute is of type String
	//   - N - the attribute is of type Number
	//   - B - the attribute is of type Binary
	AttributeType types.ScalarAttributeType
	//   - KeyType - The role that the key attribute will assume:
	//   - HASH - partition key
	//   - RANGE - sort key
	DynamoDBKeyType DynamoDBAttrType
}

// ex)
//
//	databaseHandler.CreateTable("test33", []database.DynamoDBTableAttribute{
//		{AttributeName: "part", AttributeType: types.ScalarAttributeTypeS, DynamoDBKeyType: database.DynamoDBKeyTypePart},
//		{AttributeName: "sort", AttributeType: types.ScalarAttributeTypeS, DynamoDBKeyType: database.DynamoDBKeyTypeSort},
//	})
func (d DynamoDBHandler) CreateTable(tableName string, attributes []DynamoDBTableAttribute) (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	var attributeDefinitions []types.AttributeDefinition
	var keySchemas []types.KeySchemaElement

	for _, v := range attributes {
		attributeDefinitions = append(attributeDefinitions, types.AttributeDefinition{
			AttributeName: aws.String(v.AttributeName),
			AttributeType: v.AttributeType,
		})

		switch v.DynamoDBKeyType {
		case DynamoDBKeyTypePart:
			keySchemas = append(keySchemas, types.KeySchemaElement{
				AttributeName: aws.String(v.AttributeName),
				KeyType:       types.KeyTypeHash,
			})

		case DynamoDBKeyTypeSort:
			keySchemas = append(keySchemas, types.KeySchemaElement{
				AttributeName: aws.String(v.AttributeName),
				KeyType:       types.KeyTypeRange,
			})
		}
	}

	table, err := d.Client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: attributeDefinitions,
		KeySchema:            keySchemas,
		TableName:            aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})

	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", tableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(d.Client)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

func (d DynamoDBHandler) DeleteTable() {

}

func (d DynamoDBHandler) Insert(tableName string, Object interface{}) error {
	item, err := attributevalue.MarshalMap(Object)

	if err != nil {
		panic(err)
	}
	_, err = d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (d DynamoDBHandler) Delete() {

}

func (d DynamoDBHandler) Get() {

}

func (d DynamoDBHandler) Update() {

}
