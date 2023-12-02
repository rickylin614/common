package cmongo

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var wrapper *MongoDB

func init() {
	// 創建 MongoDB 實例
	w, err := NewMongoDB().Connect(context.Background(), "mongodb://root:example@localhost:27017", "demo")
	if err != nil {
		panic(err)
	}

	wrapper = w
}

type user struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func TestMongoDBWrapper_Find(t *testing.T) {

	// 測試數據
	testCollection := "testCollection"

	// 測試 Find 方法
	// wrapper.Where("name", "John Doe")
	cond := NewQueryBuilder().
		Where("age > 10").
		GroupBy("name").
		Sum("age").
		Having("name = 'BB'")

	results := make([]bson.M, 0)

	err := wrapper.Find(context.Background(), testCollection, cond, &results)
	if err != nil {
		t.Fatalf("Find method failed: %v", err)
	}

	if len(results) == 0 {
		t.Errorf("Expected to find at least one document, but found none")
	}

	fmt.Println(results)

	wrapper.client.Disconnect(context.Background())
}
