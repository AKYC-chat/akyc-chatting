package session

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	session := Session{UserId: userId, SessionWebsocket: ws}
	sessionStorage.sessions = append(sessionStorage.sessions, session)

	// TODO: Session DB에 세션 정보 저장
	sessionEntity := SessionEntity{UserId: userId, CreateAt: time.Now().String(), SessoinId: sessionId}
	sessionStorage.Database.Insert("AYKC_SESSION", sessionEntity)
	return sessionId
}

func (sessionStorage *SessionStorage) DeleteSession(sessionId string) error {
	idx, err := sessionStorage.indexOf(sessionId)
	fmt.Println(idx)
	if err != nil {
		return err
	}

	sessionStorage.sessions = append(sessionStorage.sessions[:idx], sessionStorage.sessions[idx+1:]...)

	return nil
}

// ex)
// ctx, cancle := context.WithCancel(context.Background())
// defer cancle()
// go sessionStorage.CheckAlive(ctx)
func (sessionStorage *SessionStorage) CheckAlive(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, s := range sessionStorage.sessions {
					log.Println(s.SessionWebsocket.SessionId + " PING")
					s.SessionWebsocket.Ping(s.SessionWebsocket.SessionId)
				}
			}
		}
	}()

}

func (sessionStorage *SessionStorage) Print() {
	fmt.Println(sessionStorage.sessions)
	fmt.Println(len(sessionStorage.sessions))
}

func (sessionStorage *SessionStorage) indexOf(sessionId string) (int, error) {
	for i, s := range sessionStorage.sessions {
		if s.SessionWebsocket.SessionId == sessionId {
			return i, nil
		}
	}

	return -1, errors.New("일치하는 SessionId가 없습니다")
}
