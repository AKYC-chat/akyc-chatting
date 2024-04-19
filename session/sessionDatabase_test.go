package session_test

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/AKYC-chat/akyc-chatting/session"
	"github.com/AKYC-chat/akyc-chatting/util"
	"github.com/stretchr/testify/assert"
)

var (
	sessionDatabase = session.SessionDatabase{TableName: "AYKC_SESSION"}
)

func TestCreateSessionDatabase(t *testing.T) {
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

func TestDeleteSessionDatabase(t *testing.T) {
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

	err = sessionDatabase.DeleteSession(sessionExpectEntity)
	if err != nil {
		log.Println("Delete Sessoin Fail! userId: " + sessionExpectEntity.UserId)
		log.Fatal(err)
	}

	sessionAssertEntity, err := sessionDatabase.GetSession(sessionExpectEntity)
	if err != nil {
		log.Printf("Get Sessoin by %s Fail!\n", sessionExpectEntity.UserId)
		log.Fatal(err)
	}

	assert.Equal(t, session.SessionEntity{}, sessionAssertEntity)
}

func TestGetAllSessionsDatabase(t *testing.T) {
	count := rand.Intn(10)
	for i := 0; i < count; i++ {
		sessionId := util.SessionIdGenerator()
		userId := util.SessionIdGenerator()
		time := time.Now().String()

		sessionExpectEntity := session.SessionEntity{
			UserId: userId, SessionId: sessionId, CreateAt: time,
		}

		err := sessionDatabase.CreateSession(sessionExpectEntity)
		if err != nil {
			log.Printf("Create Sessoin Fail! count: %v\n", count)
			log.Fatal(err)
		}
	}

	sessionEntities, err := sessionDatabase.GetAllSessions()
	if err != nil {
		log.Println("Get all sessions Fail!")
		log.Fatal(err)
	}

	assert.Equal(t, count, len(sessionEntities))

	for _, s := range sessionEntities {
		sessionDatabase.DeleteSession(s)
	}

	expect := 0
	sessionEntities, err = sessionDatabase.GetAllSessions()
	if err != nil {
		log.Println("Get all sessions Fail!")
		log.Fatal(err)
	}

	assert.Equal(t, expect, len(sessionEntities))
}
