package cetcd

import (
	"fmt"
	"time"
)

func ExampleNewService() {
	srv, err := NewService("/hello", "你好啊", []string{"127.0.0.1:2379"})
	if err != nil {
		fmt.Println(err)
	}
	go srv.Start()
	time.Sleep(time.Second * 10)

	// output:
	// put key: /hello/你好啊 and value :你好啊
}
