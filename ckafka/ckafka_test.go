package ckafka

import (
	"encoding/json"
	"fmt"
)

var testBrokers []string = []string{"10.1.1.152:9092", "10.1.1.152:9093", "10.1.1.152:9094"}
var topic = "gateway"
var groupId = "b"

func ExampleWrite() {
	Manage.SetLeaderAddr("localhost:9092")
	m := map[string]interface{}{
		"agent":     "google_12345",
		"action":    "gogogo",
		"accountId": "904f5e10-4342-ad87-a0f5-8170db4a960a",
		"msg":       []byte("成功測試"),
	}
	b, _ := json.Marshal(m)
	err := Manage.Write([]byte("sampleMsg"), b, topic)
	fmt.Println(err)
	// output:
	//
}

// func ExampleWrite() {
// 	topic := "gateway"
// 	partition := 0

// 	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
// 	if err != nil {
// 		log.Fatal("failed to dial leader:", err)
// 	}

// 	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
// 	_, err = conn.WriteMessages(
// 		kafka.Message{Value: []byte("one!")},
// 		kafka.Message{Value: []byte("two!")},
// 		kafka.Message{Value: []byte("three!")},
// 	)
// 	if err != nil {
// 		log.Fatal("failed to write messages:", err)
// 	}

// 	if err := conn.Close(); err != nil {
// 		log.Fatal("failed to close writer:", err)
// 	}
// 	// output:
// 	//

// }
