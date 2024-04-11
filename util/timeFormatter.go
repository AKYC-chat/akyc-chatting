package util

import (
	"log"
	"strings"
	"time"
)

func ParseTimestamp(timestampStr string) time.Time {
	parse, err := time.Parse(time.DateTime, timestampStr)
	if err != nil {
		panic(err)
	}
	return parse
}

// MessageDateTime MessageDeduplicationId: 메세지 생성 날짜
func MessageDateTime() string {
	date := strings.Join(strings.Split(time.Now().Format(time.DateTime), " "), "/")
	defer log.Printf("message duplicationId: %v \n", date)
	return date
}
