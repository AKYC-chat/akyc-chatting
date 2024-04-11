package session

import (
	"errors"
	"fmt"

	"github.com/AKYC-chat/akyc-chatting/util"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

type SessionStorage struct {
	sessions []Session
}

type Session struct {
	SessionId        string
	SessionWebsocket websocket.Websocket
}

func (sessionStorage *SessionStorage) Append(ws websocket.Websocket) string {
	sessionId := util.SessionIdGenerator()
	session := Session{SessionId: sessionId, SessionWebsocket: ws}
	sessionStorage.sessions = append(sessionStorage.sessions, session)

	// TODO: Session DB에 세션 정보 저장
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
