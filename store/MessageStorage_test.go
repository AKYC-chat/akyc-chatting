package store_test

import (
	"github.com/AKYC-chat/akyc-chatting/store"
	"github.com/AKYC-chat/akyc-chatting/util"
	"log"
	"testing"
	"time"
)

var (
	messageStorage = store.MessageStorage{TableName: "AYKC_CHAT"}
)

func TestInsertChat(t *testing.T) {
	roomId := util.RoomIdGenerator()
	userId := util.UserIdGenerator()
	createAt := time.Now().UTC().String()
	message := "test_meassage!@#$"

	messageEntity := store.MessageEntity{
		userId, createAt, roomId, message,
	}

	err := messageStorage.InsertMessage(messageEntity)
	if err != nil {
		log.Println("InsertMessage Fail!")
		log.Fatal(err)
	}

	FindByRoomId, err := messageStorage.FindByRoomId(roomId)
	if err != nil {
		log.Println("findByRoomId Message Fail!")
		log.Fatal(err)
	}
	log.Println(FindByRoomId)

	FindByUserId, err := messageStorage.FindByRoomId(userId)
	if err != nil {
		log.Println("findByUserId Message Fail!")
		log.Fatal(err)
	}
	log.Println(FindByUserId)

}

func TestMessageStorage_FindAllRooms(t *testing.T) {
	rooms, err := messageStorage.FindAllRooms()
	if err != nil {
		log.Println("findAllRooms Fail!")
		log.Fatal(err)
	}
	for _, room := range rooms {
		log.Printf("user_id: %v \n", room.UserId)
		log.Printf("room_id: %v \n", room.RoomId)
		log.Printf("createAt: %v \n", room.CreateAt)
		log.Printf("message: %v \n", room.Message)
	}
}
