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

type DynamoDBHandler struct {
	Client *dynamodb.Client
}

// ex)
//
//	databaseHandler.CreateTable("test33", []database.DynamoDBTableAttribute{
//		{AttributeName: "part", AttributeType: types.ScalarAttributeTypeS, DynamoDBKeyType: database.DynamoDBKeyTypePart},
//		{AttributeName: "sort", AttributeType: types.ScalarAttributeTypeS, DynamoDBKeyType: database.DynamoDBKeyTypeSort},
//	})
func convertAttributeTypeToDynamoDBAttributeType(tableAttribute TableAttribute) types.AttributeDefinition {
	var attrType types.ScalarAttributeType
	switch tableAttribute.AttributeType {
	case TINYINT, BIT, BOOL, SMALLINT, MEDIUMINT, INTEGER, BIGINT,
		DECIMAL, FLOAT, DOUBLE:
		attrType = types.ScalarAttributeTypeN
	default:
		attrType = types.ScalarAttributeTypeS
	}

	return types.AttributeDefinition{
		AttributeName: aws.String(tableAttribute.AttributeName),
		AttributeType: attrType,
	}
}

func (d DynamoDBHandler) CreateTable(tableName string, attributes []TableAttribute) (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	var attributeDefinitions []types.AttributeDefinition
	var keySchemas []types.KeySchemaElement

	for _, v := range attributes {
		attributeDefinitions = append(attributeDefinitions, convertAttributeTypeToDynamoDBAttributeType(v))

		if v.PrimaryKey {
			keySchemas = append(keySchemas, types.KeySchemaElement{
				AttributeName: aws.String(v.AttributeName),
				KeyType:       types.KeyTypeHash,
			})
		}
		if v.SortKey {
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

func (d DynamoDBHandler) DeleteTable(tableName string) error {
	_, err := d.Client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName)})
	if err != nil {
		log.Printf("Couldn't delete table %v. Here's why: %v\n", tableName, err)
	}
	return err
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

func (d DynamoDBHandler) Delete(tableName string, data map[string]interface{}) error {
	var deleteValue map[string]types.AttributeValue
	for k, v := range data {
		val, err := attributevalue.Marshal(v)
		if err != nil {
			log.Println(err)
			return err
		}
		deleteValue[k] = val
	}
	_, err := d.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName), Key: deleteValue,
	})
	if err != nil {
		log.Printf("Couldn't delete item from the table. Here's why: %v\n", err)
	}
	return err
}

func (d DynamoDBHandler) Get() {

}

func (d DynamoDBHandler) Update() {

}
