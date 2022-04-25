package cetcd

import (
	"fmt"
)

func ExampleNewClientDis() {
	NewClientDis([]string{"127.0.0.1:2379"})
	s, _ := Client.GetOneService("/hello/")
	fmt.Println(s)

	// output:
	//
}
