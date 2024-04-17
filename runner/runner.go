package runner

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AKYC-chat/akyc-chatting/connections"
	"github.com/AKYC-chat/akyc-chatting/session"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

var (
	SessionStorage         = session.SessionStorage{}
	queueUrls, queueUrlErr = connections.MessageHandler.GetQueueList()
)

func ReceiveMessageFromMessageQueue() {
	for {
		messages, err := connections.MessageHandler.ReceiveMessage(queueUrls[0])

		if err != nil {
			log.Println(err)
		}

		if len(messages) != 0 {
			for _, m := range messages {
				switch m.Body {
				case "clear":
					SessionStorage.CloseCurrSessions()
				}
			}
		}
	}
}

func ReceiveWebsocket(ws *websocket.Websocket) {
	for {
		frame, err := ws.Recv()
		if err != nil {
			log.Println("옳바르지 않은 Frame양식 입니다")
			log.Println(err)
			return
		}

		switch frame.Opcode {
		case websocket.OPCODE_PONG:
			log.Println(frame.Text() + " PONG")
		case websocket.OPCODE_CLOSE:
			SessionStorage.DeleteSession(ws.SessionId)
			ws.Close()
			return
		case websocket.OPCODE_BINARY, websocket.OPCODE_FOR_TEXT:
			_, err := connections.MessageHandler.SendMessage(frame.Text(), queueUrls[0], "test4546345345")

			if err != nil {
				log.Println("Message Queue에 정상적으로 전송되지 않았습니다")
			}
		}

	}
}

func Run() {
	go ReceiveMessageFromMessageQueue()

	if queueUrlErr != nil {
		log.Println("SQS에서 Queue url 정보를 가져 올 수 없습니다")
		panic(queueUrlErr)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.New(w, r)
		go ReceiveWebsocket(ws)

		_, sessionId := SessionStorage.Append(*ws)
		log.Println("Connect Session id : ", sessionId)

		if err != nil {
			fmt.Println("Websocket 생성 실패")
			panic(err)
		}
	})

	http.ListenAndServe(":5050", nil)
}
