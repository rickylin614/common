package cmongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var wrapper *MongoDB

func init() {
	// docker run --name mongodb -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=example mongo

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
	ctx := context.Background()
	var err error

	// 測試數據
	testCollection := "testCollection"

	//testDocuments := make([]any, 10)
	//for i := 0; i < 10; i++ {
	//	testDocuments[i] = bson.M{
	//		"name": fmt.Sprintf("TestName%d", i),
	//		"age":  10 + i, // 為每個文檔設置不同的年齡
	//	}
	//}
	//err = wrapper.InsertBatch(ctx, testCollection, testDocuments)
	//if err != nil {
	//	t.Fatalf("InsertBatch method failed: %v", err)
	//}

	// 測試 Find 方法
	// wrapper.Where("name", "John Doe")
	cond := NewQueryBuilder().
		Where("age > ?", 10).
		GroupBy("name").
		Sum("age").
		Having("age > ?", 32)

	results := make([]bson.M, 0)

	err = wrapper.Find(ctx, testCollection, cond, &results)
	if err != nil {
		t.Fatalf("Find method failed: %v", err)
	}

	if len(results) == 0 {
		t.Errorf("Expected to find at least one document, but found none")
	}

	for _, v := range results {
		fmt.Printf("%+v", v)
	}

	wrapper.client.Disconnect(context.Background())
}

func TestQueryBuilder_processSimpleCondition(t *testing.T) {
	qb := NewQueryBuilder()

	d := qb.processSimpleCondition("ABC = ? and BCD > ? or QQQ > ?", false, "字串", 123, time.Now())
	fmt.Println(d)
}
