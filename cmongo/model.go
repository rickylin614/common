package cmongo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	client   *mongo.Client
	database string
}

func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

func (m *MongoDB) Connect(ctx context.Context, uri, database string) (*MongoDB, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// 獲取數據庫列表
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// 檢查指定的數據庫是否存在
	dbExists := false
	for _, db := range databases {
		if db == database {
			dbExists = true
			break
		}
	}

	if !dbExists {
		return nil, fmt.Errorf("database %s does not exist", database)
	}

	m.client = client
	m.database = database

	return m, nil
}

func (m *MongoDB) Insert(ctx context.Context, collection string, document interface{}) error {
	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.InsertOne(ctx, document)
	return err
}

func (m *MongoDB) InsertBatch(ctx context.Context, collection string, documents []interface{}) error {
	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.InsertMany(ctx, documents)
	return err
}

func (m *MongoDB) Update(ctx context.Context, collection string, qb *QueryBuilder, update interface{}) error {
	filter, _, _ := qb.Build()
	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

type UpdateModel struct {
	Filter *QueryBuilder
	Update interface{}
}

func (m *MongoDB) UpdateBatch(ctx context.Context, collection string, updates []UpdateModel) error {
	coll := m.client.Database(m.database).Collection(collection)
	models := make([]mongo.WriteModel, len(updates))
	for i, update := range updates {
		filter, _, _ := update.Filter.Build()
		models[i] = mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update.Update)
	}
	_, err := coll.BulkWrite(ctx, models)
	return err
}

func (m *MongoDB) Delete(ctx context.Context, collection string, qb *QueryBuilder) error {
	filter, _, _ := qb.Build()
	coll := m.client.Database(m.database).Collection(collection)
	_, err := coll.DeleteOne(ctx, filter)
	return err
}

func (m *MongoDB) DeleteBatch(ctx context.Context, collection string, deletes []*QueryBuilder) error {
	coll := m.client.Database(m.database).Collection(collection)
	models := make([]mongo.WriteModel, len(deletes))
	for i, delete := range deletes {
		filter, _, _ := delete.Build()
		models[i] = mongo.NewDeleteOneModel().SetFilter(filter)
	}
	_, err := coll.BulkWrite(ctx, models)
	return err
}

func (m *MongoDB) Find(ctx context.Context, table string, qb *QueryBuilder, results any) error {
	filter, group, opts := qb.Build()
	collection := m.client.Database(m.database).Collection(table)
	if len(group) > 0 {
		pipeline := mongo.Pipeline{
			{{"$match", filter}},
			{{"$group", group}},
		}
		if len(qb.having) > 0 {
			pipeline = append(pipeline, bson.D{{"$match", qb.having}})
		}
		cursor, err := collection.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)
		return cursor.All(ctx, results)
	} else {
		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)
		return cursor.All(ctx, results)
	}
}

func (m *MongoDB) Count(ctx context.Context, table string, qb *QueryBuilder) (int64, error) {
	filter, group, _ := qb.Build()
	collection := m.client.Database(m.database).Collection(table)
	if len(group) > 0 {
		pipeline := mongo.Pipeline{
			{{"$match", filter}},
			{{"$group", group}},
		}
		if len(qb.having) > 0 {
			pipeline = append(pipeline, bson.D{{"$match", qb.having}})
		}
		pipeline = append(pipeline, bson.D{{"$count", "count"}})

		cursor, err := collection.Aggregate(ctx, pipeline)
		if err != nil {
			return 0, err
		}
		defer cursor.Close(ctx)
		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			return 0, err
		}
		if len(results) > 0 {
			return results[0]["count"].(int64), nil
		}
		return 0, nil
	} else {
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
}

type QueryBuilder struct {
	filter    bson.D
	sort      bson.D
	limit     int64
	offset    int64
	group     bson.D
	sumFields map[string]bool
	having    bson.D
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{sumFields: map[string]bool{}}
}

func (qb *QueryBuilder) Build() (bson.D, bson.D, *options.FindOptions) {
	findOptions := options.Find()
	if qb.limit > 0 {
		findOptions.SetLimit(qb.limit)
	}
	if qb.offset > 0 {
		findOptions.SetSkip(qb.offset)
	}
	if len(qb.sort) > 0 {
		findOptions.SetSort(qb.sort)
	}
	return qb.filter, qb.group, findOptions
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	parts := strings.Fields(condition)
	if len(parts) != 3 {
		return qb
	}

	field := parts[0]
	operator := parts[1]
	var value any
	value = parts[2]
	if v, ok := value.(string); ok {
		i, err := strconv.Atoi(v)
		if strings.HasPrefix(v, `'`) && strings.HasSuffix(v, `'`) {
			value = v[1 : len(v)-1]
		} else if err == nil {
			value = i
		}
	}

	var filter bson.E
	switch operator {
	case "=":
		filter = bson.E{Key: field, Value: value}
	case ">":
		filter = bson.E{Key: field, Value: bson.M{"$gt": value}}
	case "<":
		filter = bson.E{Key: field, Value: bson.M{"$lt": value}}
	case ">=":
		filter = bson.E{Key: field, Value: bson.M{"$gte": value}}
	case "<=":
		filter = bson.E{Key: field, Value: bson.M{"$lte": value}}
	case "in":
		values := strings.Split(parts[2], ",")
		filter = bson.E{Key: field, Value: bson.M{"$in": values}}
	case "like":
		likeValue := parts[2]
		regexPattern := ""
		// 檢測並轉換不同的模式
		if strings.HasPrefix(likeValue, "%") && strings.HasSuffix(likeValue, "%") {
			// 匹配任意位置的子串（%value%）
			regexPattern = strings.Trim(likeValue, "%")
		} else if strings.HasPrefix(likeValue, "%") {
			// 匹配結尾的子串（%value）
			regexPattern = strings.TrimLeft(likeValue, "%") + "$"
		} else if strings.HasSuffix(likeValue, "%") {
			// 匹配開頭的子串（value%）
			regexPattern = "^" + strings.TrimRight(likeValue, "%")
		}
		regex := bson.M{"$regex": regexPattern, "$options": "i"} // 使用 'i' 選項實現不區分大小寫的匹配
		filter = bson.E{Key: field, Value: regex}
	default:
		return qb
	}

	qb.filter = append(qb.filter, filter)
	return qb
}

// Sort default ASC, use `-` prefix as desc, example: "-age" is "age desc"
func (qb *QueryBuilder) Sort(fields ...string) *QueryBuilder {
	for _, field := range fields {
		order := 1 // Ascending
		if strings.HasPrefix(field, "-") {
			order = -1 // Descending
			field = strings.TrimPrefix(field, "-")
		}
		qb.sort = append(qb.sort, bson.E{Key: field, Value: order})
	}
	return qb
}

func (qb *QueryBuilder) Limit(limit int64) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int64) *QueryBuilder {
	qb.offset = offset
	return qb
}

// func (qb *QueryBuilder) GroupBy(id interface{}) *QueryBuilder {
// 	qb.group = append(qb.group, bson.E{Key: "_id", Value: id})
// 	return qb
// }

func (qb *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	groupFields := bson.D{}
	for _, field := range fields {
		groupFields = append(groupFields, bson.E{Key: field, Value: "$" + field})
	}
	qb.group = append(qb.group, bson.E{Key: "_id", Value: groupFields})
	return qb
}

func (qb *QueryBuilder) Sum(fields ...string) *QueryBuilder {
	for _, field := range fields {
		qb.sumFields[field] = true
		qb.group = append(qb.group, bson.E{Key: "total_$" + field, Value: bson.D{{"$sum", field}}})
	}
	return qb
}

func (qb *QueryBuilder) Having(condition string) *QueryBuilder {
	parts := strings.Fields(condition)
	if len(parts) != 3 {
		return qb
	}

	field := parts[0]
	// sum field with prefix total_$
	if _, ok := qb.sumFields[field]; ok {
		field = "total_$" + field
	}
	operator := parts[1]
	var value any
	value = parts[2]
	if v, ok := value.(string); ok {
		i, err := strconv.Atoi(v)
		if strings.HasPrefix(v, `'`) && strings.HasSuffix(v, `'`) {
			value = v[1 : len(v)-1]
		} else if err == nil {
			value = i
		}
	}

	var filter bson.E
	switch operator {
	case "=":
		filter = bson.E{Key: field, Value: value}
	case ">":
		filter = bson.E{Key: field, Value: bson.M{"$gt": value}}
	case "<":
		filter = bson.E{Key: field, Value: bson.M{"$lt": value}}
	case ">=":
		filter = bson.E{Key: field, Value: bson.M{"$gte": value}}
	case "<=":
		filter = bson.E{Key: field, Value: bson.M{"$lte": value}}
	case "in":
		values := strings.Split(parts[2], ",")
		filter = bson.E{Key: field, Value: bson.M{"$in": values}}
	case "like":
		likeValue := parts[2]
		regexPattern := ""
		// 檢測並轉換不同的模式
		if strings.HasPrefix(likeValue, "%") && strings.HasSuffix(likeValue, "%") {
			// 匹配任意位置的子串（%value%）
			regexPattern = strings.Trim(likeValue, "%")
		} else if strings.HasPrefix(likeValue, "%") {
			// 匹配結尾的子串（%value）
			regexPattern = strings.TrimLeft(likeValue, "%") + "$"
		} else if strings.HasSuffix(likeValue, "%") {
			// 匹配開頭的子串（value%）
			regexPattern = "^" + strings.TrimRight(likeValue, "%")
		}
		regex := bson.M{"$regex": regexPattern, "$options": "i"} // 使用 'i' 選項實現不區分大小寫的匹配
		filter = bson.E{Key: field, Value: regex}
	default:
		return qb
	}

	qb.having = append(qb.having, filter)
	return qb
}
