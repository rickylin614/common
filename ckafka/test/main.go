package main

import (
	"fmt"

	"github.com/rickylin614/common/ckafka"
)

var brokers []string = []string{"127.0.0.1:9092"}
var topic = "gameserver"
var groupId = "b"

func main() {
	fmt.Println("ready")
	ckafka.Manage.SetBrokers(brokers)
	chanGateway, _ := ckafka.Manage.NewReader(topic, groupId)
	for {
		n := <-chanGateway
		fmt.Printf("key:%s value:%s\n", n.Key, n.Value)
	}
	// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
	// m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
}
