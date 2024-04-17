package session_test

import (
	"log"
	"testing"
	"time"

	"github.com/AKYC-chat/akyc-chatting/session"
	"github.com/AKYC-chat/akyc-chatting/util"
	"github.com/stretchr/testify/assert"
)

var (
	sessionDatabase = session.SessionDatabase{}
)

func TestCreateSession(t *testing.T) {
	sessionId := util.SessionIdGenerator()
	userId := util.SessionIdGenerator()
	time := time.Now().String()

	sessionExpectEntity := session.SessionEntity{
		UserId: userId, SessionId: sessionId, CreateAt: time,
	}

	err := sessionDatabase.CreateSession(sessionExpectEntity)
	if err != nil {
		log.Println("Create Sessoin Fail!")
		log.Fatal(err)
	}

	sessionAssertEntity, err := sessionDatabase.GetSession(sessionExpectEntity)
	if err != nil {
		log.Printf("Get Sessoin by %s Fail!\n", sessionExpectEntity.UserId)
		log.Fatal(err)
	}

	assert.Equal(t, sessionExpectEntity, sessionAssertEntity)

	err = sessionDatabase.DeleteSession(sessionAssertEntity)
	if err != nil {
		log.Println("Delete Sessoin Fail! userId: " + sessionExpectEntity.UserId)
		log.Fatal(err)
	}
}

func TestDeleteSession(t *testing.T) {
	sessionId := util.SessionIdGenerator()
	userId := util.SessionIdGenerator()
	time := time.Now().String()

	sessionExpectEntity := session.SessionEntity{
		UserId: userId, SessionId: sessionId, CreateAt: time,
	}

	err := sessionDatabase.CreateSession(sessionExpectEntity)
	if err != nil {
		log.Println("Create Sessoin Fail!")
		log.Fatal(err)
	}

	sessionAssertEntity, err := sessionDatabase.GetSession(sessionExpectEntity)
	if err != nil {
		log.Printf("Get Sessoin by %s Fail!\n", sessionExpectEntity.UserId)
		log.Fatal(err)
	}

	assert.Equal(t, sessionExpectEntity, sessionAssertEntity)

	err = sessionDatabase.DeleteSession(sessionAssertEntity)
	if err != nil {
		log.Println("Delete Sessoin Fail! userId: " + sessionExpectEntity.UserId)
		log.Fatal(err)
	}
}
