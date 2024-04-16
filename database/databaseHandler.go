package database

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type AttributeType string

const (
	BIT        AttributeType = "BIT"
	BOOL       AttributeType = "BOOL"
	BOOLEAN    AttributeType = "BOOLEAN"
	TINYINT    AttributeType = "TINYINT"
	SMALLINT   AttributeType = "SMALLINT"
	MEDIUMINT  AttributeType = "MEDIUMINT"
	INTEGER    AttributeType = "INTEGER"
	BIGINT     AttributeType = "BIGINT"
	DECIMAL    AttributeType = "DECIMAL"
	FLOAT      AttributeType = "FLOAT"
	DOUBLE     AttributeType = "DOUBLE"
	DATE       AttributeType = "DATE"
	DATETIME   AttributeType = "DATETIME"
	TIMESTAMP  AttributeType = "TIMESTAMP"
	TIME       AttributeType = "TIME"
	YEAR       AttributeType = "YEAR"
	CHAR       AttributeType = "CHAR"
	VARCHAR    AttributeType = "VARCHAR"
	TINYBLOB   AttributeType = "TINYBLOB"
	TEXT       AttributeType = "TEXT"
	MEDIUMTEXT AttributeType = "MEDIUMTEXT"
	LONGTEXT   AttributeType = "LONGTEXT"
	ENUM       AttributeType = "ENUM"
	SET        AttributeType = "SET"
)

type TableAttribute struct {
	// 속성이름(required)
	AttributeName string
	// 자료형(required)
	// - DynamoDB의 경우
	AttributeType AttributeType
	// 기본키 여부(optional)
	// - DynamoDB의 경우 Partition key로 적용
	PrimaryKey bool
	// SortKey(optional)
	// - DynamoDB 전용
	SortKey bool
}

type DatabaseHandler interface {
	Insert(tableName string, Object interface{}) error
	Delete()
	Get()
	Update()
	CreateTable(tableName string, attributes []TableAttribute) (*types.TableDescription, error)
	DeleteTable(tableName string, data map[string]interface{}) error
}
