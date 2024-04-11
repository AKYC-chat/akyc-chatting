package util

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func ParseTimestamp(timestampStr string) (*time.Time, error) {
	log.Println(timestampStr)

	// 시간 문자열 조합
	millisec, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		log.Println("Error Parsing Timestamp")
		return &time.Time{}, err
	}
	parseMillisec := time.Unix(0, millisec*int64(time.Millisecond)).String()
	parseDateTime, err := time.Parse(time.DateTime, strings.Split(parseMillisec, " +")[0])
	if err != nil {
		return &time.Time{}, err
	}
	log.Println(parseDateTime)
	return &parseDateTime, nil

}

// MessageDateTime MessageDeduplicationId: 메세지 생성 날짜
func MessageDateTime() string {
	date := strings.Join(strings.Split(time.Now().Format(time.DateTime), " "), "/")
	defer log.Printf("message duplicationId: %v \n", date)
	return date
}
