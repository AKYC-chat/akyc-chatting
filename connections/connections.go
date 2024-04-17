package connections

import "github.com/AKYC-chat/akyc-chatting/connector"

var (
	MessageHandler     = connector.SqsGetConnection()
	DatabaseConnection = connector.DynamoDBGetConnection()
)
