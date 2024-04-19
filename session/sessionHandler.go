package session

import (
	"errors"
	"log"
	"time"

	"github.com/AKYC-chat/akyc-chatting/util"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

var (
	sessionDatabase = SessionDatabase{TableName: "AYKC_SESSION"}
)

type SessionStorage struct {
	sessions []Session
}

type Session struct {
	UserId           string
	CreateAt         string
	SessionWebsocket websocket.Websocket
}

func (sessionStorage *SessionStorage) Append(ws websocket.Websocket) (string, string) {
	sessionId := util.SessionIdGenerator()
	userId := util.SessionIdGenerator()
	time := time.Now().UTC().String()

	session := Session{UserId: userId, SessionWebsocket: ws, CreateAt: time}
	sessionStorage.sessions = append(sessionStorage.sessions, session)

	sessionEntity := SessionEntity{UserId: userId, CreateAt: time, SessionId: sessionId}
	err := sessionDatabase.CreateSession(sessionEntity)
	if err != nil {
		log.Panicf("Couldn't add item to table. Here's why: %v\n", err)
	}

	return userId, sessionId
}

func (sessionStorage *SessionStorage) DeleteSession(sessionId string) error {
	idx, err := sessionStorage.indexOf(sessionId)
	if err != nil {
		return err
	}

	session := sessionStorage.sessions[idx]

	util.DeleteElement(sessionStorage.sessions, idx)

	sessionEntity := SessionEntity{UserId: session.UserId, SessionId: sessionId, CreateAt: session.CreateAt}
	err = sessionDatabase.DeleteSession(sessionEntity)
	if err != nil {
		log.Println(err)
	}

	return err
}

// 현재 접속중인 세션들 모두 종료
func (sessionStorage *SessionStorage) CloseCurrSessions() error {
	sessions := sessionStorage.sessions

	for _, s := range sessions {
		err := sessionStorage.DeleteSession(s.SessionWebsocket.SessionId)
		s.SessionWebsocket.Send(websocket.Frame{Opcode: websocket.OPCODE_CLOSE, Payload: []byte("close"), PayloadLength: len("close")})
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (sessionStorage *SessionStorage) GetSessions() []SessionEntity {
	sessionEntities, err := sessionDatabase.GetAllSessions()

	if err != nil {
		log.Println(err)
	}
	return sessionEntities
}

func (sessionStorage *SessionStorage) GetSessionByUserId(userId string) Session {
	sessionEntity, err := sessionDatabase.GetSessionByUserId(userId)
	if err != nil {
		panic(err)
	}

	idx, err := sessionStorage.indexOf(sessionEntity.SessionId)
	if err != nil {
		log.Println("Session id에 해당하는 유저가 없습니다.")
	}

	return sessionStorage.sessions[idx]
}

func (sessionStorage *SessionStorage) GetSessionCount() int {
	return len(sessionStorage.sessions)
}

func (sessionStorage *SessionStorage) indexOf(sessionId string) (int, error) {
	for i, s := range sessionStorage.sessions {
		if s.SessionWebsocket.SessionId == sessionId {
			return i, nil
		}
	}

	return -1, errors.New("일치하는 SessionId가 없습니다")
}
