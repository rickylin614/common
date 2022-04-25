package utils

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name  string
	Age   int
	Score int
}

func ExampleNewList() {
	pList := []Person{
		{Name: "a", Age: 10, Score: 10},
		{Name: "a", Age: 90, Score: 10},
		{Name: "a", Age: 50, Score: 10},
		{Name: "a", Age: 10, Score: 90},
		{Name: "a", Age: 10, Score: 40},
		{Name: "a", Age: 10, Score: 65},
		{Name: "c", Age: 10, Score: 10},
		{Name: "c", Age: 90, Score: 10},
		{Name: "c", Age: 50, Score: 10},
		{Name: "c", Age: 10, Score: 90},
		{Name: "c", Age: 10, Score: 40},
		{Name: "c", Age: 10, Score: 65},
		{Name: "b", Age: 10, Score: 10},
		{Name: "b", Age: 90, Score: 10},
		{Name: "b", Age: 50, Score: 10},
		{Name: "b", Age: 10, Score: 90},
		{Name: "b", Age: 10, Score: 40},
		{Name: "b", Age: 10, Score: 65},
	}
	l := NewList(pList).SortBy("Name").ThenSort("Age", true).ThenSort("Score").ToList()
	for _, v := range l {
		b, _ := json.Marshal(v)
		fmt.Printf("%s\n", b)
	}

	// output:
	// {"Name":"a","Age":90,"Score":10}
	// {"Name":"a","Age":50,"Score":10}
	// {"Name":"a","Age":10,"Score":10}
	// {"Name":"a","Age":10,"Score":40}
	// {"Name":"a","Age":10,"Score":65}
	// {"Name":"a","Age":10,"Score":90}
	// {"Name":"b","Age":90,"Score":10}
	// {"Name":"b","Age":50,"Score":10}
	// {"Name":"b","Age":10,"Score":10}
	// {"Name":"b","Age":10,"Score":40}
	// {"Name":"b","Age":10,"Score":65}
	// {"Name":"b","Age":10,"Score":90}
	// {"Name":"c","Age":90,"Score":10}
	// {"Name":"c","Age":50,"Score":10}
	// {"Name":"c","Age":10,"Score":10}
	// {"Name":"c","Age":10,"Score":40}
	// {"Name":"c","Age":10,"Score":65}
	// {"Name":"c","Age":10,"Score":90}
}

func ExampleList_Contains() {
	s := []string{"a", "b", "d", "f"}
	list := NewList(s)
	fmt.Println(list.Contains("a"))
	fmt.Println(list.Contains("b"))
	fmt.Println(list.Contains("c"))
	fmt.Println(list.Contains("d"))
	fmt.Println(list.Contains("e"))
	fmt.Println(list.Contains("f"))
	fmt.Println(list.Contains("g"))

	// output:
	// true
	// true
	// false
	// true
	// false
	// true
	// false
}
