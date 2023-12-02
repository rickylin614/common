package cmongo

// import (
// 	"context"
// 	"fmt"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type QueryBuilder struct {
// 	filter bson.D
// 	sort   bson.D
// 	limit  int64
// 	skip   int64
// }

// func NewQueryBuilder() *QueryBuilder {
// 	return &QueryBuilder{}
// }

// func (qb *QueryBuilder) Where(key string, value interface{}) *QueryBuilder {
// 	qb.filter = append(qb.filter, bson.E{Key: key, Value: value})
// 	return qb
// }

// func (qb *QueryBuilder) Sort(key string, value int) *QueryBuilder {
// 	qb.sort = append(qb.sort, bson.E{Key: key, Value: value})
// 	return qb
// }

// func (qb *QueryBuilder) Limit(value int64) *QueryBuilder {
// 	qb.limit = value
// 	return qb
// }

// func (qb *QueryBuilder) Skip(value int64) *QueryBuilder {
// 	qb.skip = value
// 	return qb
// }

// func (qb *QueryBuilder) Build() *options.FindOptions {
// 	findOptions := options.Find()
// 	if qb.sort != nil {
// 		findOptions.SetSort(qb.sort)
// 	}
// 	if qb.limit > 0 {
// 		findOptions.SetLimit(qb.limit)
// 	}
// 	if qb.skip > 0 {
// 		findOptions.SetSkip(qb.skip)
// 	}
// 	return findOptions
// }

// type MongoDB struct {
// 	client   *mongo.Client
// 	database string
// }

// func NewMongoDB(uri, database string) (*MongoDB, error) {
// 	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &MongoDB{client: client, database: database}, nil
// }

// func (m *MongoDB) Find(ctx context.Context, collection string, qb *QueryBuilder, results interface{}) error {
// 	c := m.client.Database(m.database).Collection(collection)
// 	cursor, err := c.Find(ctx, qb.filter, qb.Build())
// 	if err != nil {
// 		return err
// 	}
// 	defer cursor.Close(ctx)
// 	return cursor.All(ctx, results)
// }

// func main() {
// 	mongoDB, err := NewMongoDB("mongodb://localhost:27017", "test")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	qb := NewQueryBuilder().
// 		Where("age", bson.D{{"$gt", 20}}).
// 		Sort("name", 1).
// 		Limit(10).
// 		Skip(0)

// 	var results []bson.M
// 	err = mongoDB.Find(context.Background(), "users", qb, &results)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println(results)
// }
