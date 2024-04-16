package runner

import (
	"fmt"
	"net/http"

	"github.com/AKYC-chat/akyc-chatting/connector"
	"github.com/AKYC-chat/akyc-chatting/session"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

var (
	databaseHandler        = connector.DynamoDBGetConnection()
	sessionStorage         = session.SessionStorage{Database: databaseHandler}
	messageHandler         = connector.SqsGetConnection()
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

func Run() {
	go ReceiveMessageFromMessageQueue()

	if queueUrlErr != nil {
		fmt.Println("SQS에서 Queue url 정보를 가져 올 수 없습니다")
		panic(queueUrlErr)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.New(w, r)
		sessionId := sessionStorage.Append(*ws)
		fmt.Println("Connect Session id : ", sessionId)
		ws.Send(websocket.Frame{Payload: []byte(sessionId), PayloadLength: len(sessionId), Opcode: websocket.OPCODE_FOR_TEXT})

		if err != nil {
			fmt.Println("Websocket 생성 실패")
			panic(err)
		}

		for {
			frame, err := ws.Recv()

			switch frame.Opcode {
			case websocket.OPCODE_CLOSE:
				// TODO: 해당 유저의 session id를 가져와야함
				sessionStorage.Print()
				sessionStorage.DeleteSession("test")
			}

			if err != nil {
				fmt.Println("옳바르지 않은 Frame양식 입니다")
				panic(err)
			}

			messageId, err := messageHandler.SendMessage(frame.Text(), queueUrls[0], "test4546345345")

			if err != nil {
				fmt.Println("Message Queue에 정상적으로 전송되지 않았습니다")
				// panic(err)
			}

			fmt.Println(messageId)
		}

	})

	http.ListenAndServe(":5050", nil)
}
