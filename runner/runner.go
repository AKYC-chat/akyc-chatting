package runner

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AKYC-chat/akyc-chatting/connector"
	"github.com/AKYC-chat/akyc-chatting/session"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

var (
	messageHandler         = connector.SqsGetConnection()
	databaseHandler        = connector.DynamoDBGetConnection()
	sessionStorage         = session.SessionStorage{Database: databaseHandler}
	queueUrls, queueUrlErr = messageHandler.GetQueueList()
)

func ReceiveMessageFromMessageQueue() {
	for {
		messages, err := messageHandler.ReceiveMessage(queueUrls[0])

		if err != nil {
			fmt.Println("PANIC!!!!!!")
			panic(err)
		}

		if len(messages) != 0 {
			for _, m := range messages {
				switch m.Body {
				case "create table":

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
		}

		switch frame.Opcode {
		case websocket.OPCODE_PONG:
			log.Println(frame.Text() + " PONG")
		case websocket.OPCODE_CLOSE:
			sessionStorage.DeleteSession(ws.SessionId)
			sessionStorage.Print()
			ws.Close()
			return
		case websocket.OPCODE_BINARY, websocket.OPCODE_FOR_TEXT:
			messageId, err := messageHandler.SendMessage(frame.Text(), queueUrls[0], "test4546345345")

			if err != nil {
				log.Println("Message Queue에 정상적으로 전송되지 않았습니다")
			}
			fmt.Println(messageId)
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

		sessionId := sessionStorage.Append(*ws)
		log.Println("Connect Session id : ", sessionId)
		sessionStorage.Print()

		if err != nil {
			fmt.Println("Websocket 생성 실패")
			panic(err)
		}
	})

	http.ListenAndServe(":5050", nil)
}
