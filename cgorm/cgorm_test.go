package cgorm

import (
	"fmt"
)

/*
	if you never use gorm,
	please check https://gorm.io/docs/index.html , and get the way how to use.
*/

func Init_test() {
	InitDB("10.1.1.152:3306", "gl", "root", "123456", "")
}

func ExampleGetDB() {
	Init_test()
	db := GetDB()
	fmt.Println(db)

	// output:
	// end
}

type Demo struct {
	Id    int64 `gorm:"primaryKey;not null"`
	Name  string
	Phone string
	Age   int
}

func (Demo) TableName() string {
	return "demo"
}
