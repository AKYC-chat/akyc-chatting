package session_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AKYC-chat/akyc-chatting/session"
	"github.com/AKYC-chat/akyc-chatting/websocket"
	"github.com/stretchr/testify/assert"
)

var (
	mockSessionStorage = session.SessionStorage{}
)

func TestAppend(t *testing.T) {
	expect := 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockWs, err := websocket.New(w, r)
		userId, sessionId := mockSessionStorage.Append(*mockWs)

		if err != nil {
			log.Fatalf("Couldn't create websocket. Here's why: %v", err)
		}
		assert.Equal(t, expect, mockSessionStorage.GetSessionCount())
		session := mockSessionStorage.GetSessionByUserId(userId)
		assert.Equal(t, sessionId, session.SessionWebsocket.SessionId)

		mockSessionStorage.CloseCurrSessions()
		expect = 0
		assert.Equal(t, expect, mockSessionStorage.GetSessionCount())
	}))
	defer server.Close()

	http.Get(server.URL + "/ws")
}
