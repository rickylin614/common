package utils

import (
	"fmt"
)

func ErrRecover() {
	err := recover()
	if err != nil {
		fmt.Println("Recover error:", err)
	}
}
