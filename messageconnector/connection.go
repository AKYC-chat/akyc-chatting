package messageconnector

type Connection interface {
	GetConnection() Client
}
