package runner

import (
	"fmt"
	"net/http"

	"github.com/AKYC-chat/akyc-chatting/connector"
	"github.com/AKYC-chat/akyc-chatting/message"
	"github.com/AKYC-chat/akyc-chatting/session"
	"github.com/AKYC-chat/akyc-chatting/websocket"
)

func ReceiveMessageFromMessageQueue(messageHandler message.MessageHandler, messageUrl string) {
	for {
		messages, err := messageHandler.ReceiveMessage(messageUrl)

		if err != nil {
			panic(err)
		}

		fmt.Println(messages)
	}
}

func Run() {
	sessionStorage := session.SessionStorage{}
	messageHandler := connector.SqsGetConnection()
	queueUrls, err := messageHandler.GetQueueList()

	// go ReceiveMessageFromMessageQueue(messageHandler, queueUrls[0])

	if err != nil {
		fmt.Println("SQS에서 Queue url 정보를 가져 올 수 없습니다")
		panic(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.New(w, r)
		sessionId := sessionStorage.Append(*ws)
		fmt.Println("Connect Session id : ", sessionId)

		sessionStorage.Print()
		if err != nil {
			fmt.Println("Websocket 생성 실패")
			panic(err)
		}

		for {
			frame, err := ws.Recv()

			if err != nil {
				fmt.Println("옳바르지 않은 Frame양식 입니다")
				panic(err)
			}

			messageId, err := messageHandler.SendMessage(frame.Text(), queueUrls[0], "test4546345345")

			if err != nil {
				fmt.Println("Message Queue에 정상적으로 전송되지 않았습니다")
				panic(err)
			}

			fmt.Println(messageId)
		}

	})

	http.ListenAndServe(":5050", nil)
}
