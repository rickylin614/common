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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
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
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) Update(ctx context.Context, collection string, query interface{}, update interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) UpdateBatch(ctx context.Context, collection string, updates []mongo.WriteModel) error {
	c := m.client.Database(m.database).Collection(collection)
	_, err := c.BulkWrite(ctx, updates)
	return err
}

func (m *MongoDB) Delete(ctx context.Context, collection string, query interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDB) DeleteBatch(ctx context.Context, collection string, deletes []mongo.WriteModel) error {
	c := m.client.Database(m.database).Collection(collection)
	_, err := c.BulkWrite(ctx, deletes)
	return err
}

func (m *MongoDB) Find(ctx context.Context, table string, qb *queryBuilder, results any) error {
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

func (m *MongoDB) Count(ctx context.Context, table string, qb *queryBuilder) (int64, error) {
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

type queryBuilder struct {
	filter    bson.D
	sort      bson.D
	limit     int64
	offset    int64
	group     bson.D
	sumFields map[string]bool
	having    bson.D
}

func NewQueryBuilder() *queryBuilder {
	return &queryBuilder{sumFields: map[string]bool{}}
}

func (qb *queryBuilder) Build() (bson.D, bson.D, *options.FindOptions) {
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

// func (qb *QueryBuilder) Where(field string, value interface{}, comparison ...string) *QueryBuilder {
// 	var filter bson.E

// 	compare := ""
// 	if len(comparison) > 0 {
// 		compare = comparison[0]
// 	}

// 	switch compare {
// 	case "=":
// 		filter = bson.E{Key: field, Value: value}
// 	case ">":
// 		filter = bson.E{Key: field, Value: bson.M{"$gt": value}}
// 	case "<":
// 		filter = bson.E{Key: field, Value: bson.M{"$lt": value}}
// 	case ">=":
// 		filter = bson.E{Key: field, Value: bson.M{"$gte": value}}
// 	case "<=":
// 		filter = bson.E{Key: field, Value: bson.M{"$lte": value}}
// 	default:
// 		filter = bson.E{Key: field, Value: value}
// 	}

// 	qb.filter = append(qb.filter, filter)
// 	return qb
// }

func (qb *queryBuilder) Where(condition string) *queryBuilder {
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
	default:
		return qb
	}

	qb.filter = append(qb.filter, filter)
	return qb
}

// Sort default ASC, use `-` prefix as desc, example: "-age" is "age desc"
func (qb *queryBuilder) Sort(fields ...string) *queryBuilder {
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

func (qb *queryBuilder) Limit(limit int64) *queryBuilder {
	qb.limit = limit
	return qb
}

func (qb *queryBuilder) Offset(offset int64) *queryBuilder {
	qb.offset = offset
	return qb
}

// func (qb *QueryBuilder) GroupBy(id interface{}) *QueryBuilder {
// 	qb.group = append(qb.group, bson.E{Key: "_id", Value: id})
// 	return qb
// }

func (qb *queryBuilder) GroupBy(fields ...string) *queryBuilder {
	groupFields := bson.D{}
	for _, field := range fields {
		groupFields = append(groupFields, bson.E{Key: field, Value: "$" + field})
	}
	qb.group = append(qb.group, bson.E{Key: "_id", Value: groupFields})
	return qb
}

func (qb *queryBuilder) Sum(fields ...string) *queryBuilder {
	for _, field := range fields {
		qb.sumFields[field] = true
		qb.group = append(qb.group, bson.E{Key: "total_$" + field, Value: bson.D{{"$sum", field}}})
	}
	return qb
}

func (qb *queryBuilder) Having(condition string) *queryBuilder {
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
	default:
		return qb
	}

	qb.having = append(qb.having, filter)
	return qb
}
