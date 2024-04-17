package session_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllSessions(t *testing.T) {
	t.Run("Join Websocket", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/ws", nil)
		response := httptest.NewRecorder()

		fmt.Print(request)
		fmt.Print(response)

	})
}
