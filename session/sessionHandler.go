package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/AKYC-chat/akyc-chatting/database"
	"github.com/AKYC-chat/akyc-chatting/util"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

type SessionStorage struct {
	Database database.DatabaseHandler
	sessions []Session
}

type Session struct {
	UserId           string
	SessionId        string
	SessionWebsocket websocket.Websocket
}

type SessionEntity struct {
	UserId    string `dynamodbav:"user_id"`
	CreateAt  string `dynamodbav:"create_at"`
	SessoinId string `dynamodbav:"session_id"`
}

func (sessionStorage *SessionStorage) Append(ws websocket.Websocket) string {
	sessionId := util.SessionIdGenerator()
	userId := util.SessionIdGenerator()
	session := Session{SessionId: sessionId, UserId: userId, SessionWebsocket: ws}
	sessionStorage.sessions = append(sessionStorage.sessions, session)

	// TODO: Session DB에 세션 정보 저장
	sessionEntity := SessionEntity{UserId: userId, CreateAt: time.Now().String(), SessoinId: sessionId}
	sessionStorage.Database.Insert("AYKC_SESSION", sessionEntity)
	return sessionId
}

func (sessionStorage *SessionStorage) DeleteSession(sessionId string) error {
	idx, err := sessionStorage.indexOf(sessionId)

	if err != nil {
		return err
	}

	sessionStorage.sessions = append(sessionStorage.sessions[:idx], sessionStorage.sessions[idx:]...)

	return nil
}

func (SessionStorage *SessionStorage) Print() {
	fmt.Println(SessionStorage.sessions)
}

func (sessionStorage *SessionStorage) indexOf(sessionId string) (int, error) {
	for i, s := range sessionStorage.sessions {
		if s.SessionId == sessionId {
			return i, nil
		}
	}

	return -1, errors.New("일치하는 SessionId가 없습니다")
}
