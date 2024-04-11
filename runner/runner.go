package runner

import (
	"fmt"
	"net/http"

	"github.com/AKYC-chat/akyc-chatting/connector"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

type WebSocketSession struct {
	sessionId string
	ws        websocket.Websocket
}

func Run() {
	sessionList := make([]WebSocketSession, 10)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.New(w, r)
		session := WebSocketSession{sessionId: r.Header.Get("Sec-Websocket-Key"), ws: *ws}

		sessionList = append(sessionList, session)

		if err != nil {
			fmt.Println(err)
		}

		messageHandler := connector.SqsGetConnection()
		queueUrls, err := messageHandler.GetQueueList()

		if err != nil {
			fmt.Println("SQS에서 Queue url 정보를 가져 올 수 없습니다.")
			panic(err)
		}

		for {
			frame, err := ws.Recv()

			if err != nil {
				fmt.Println("옳바르지 않은 Frame양식 입니다.")
				panic(err)
			}

			for _, s := range sessionList {
				if len(s.sessionId) > 0 {
					s.ws.Send(frame)
				}
			}
			messageHandler.SendMessage(frame.Text(), queueUrls[0], "test4546345345")
		}

	})

	http.ListenAndServe(":5050", nil)
}
