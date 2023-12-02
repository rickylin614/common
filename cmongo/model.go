package cmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBWrapper struct {
	client     *mongo.Client
	collection *mongo.Collection
	database   string

	queryFilter bson.D
	sort        bson.D
	limit       int64
	offset      int64
	groupBy     bson.D
}

func (m *mongoDBWrapper) Clone() MongoDBWrapper {
	cloned := *m // 淺拷貝
	// 如果有引用類型的字段，需要進行深拷貝
	// 例如，如果 queryFilter 是引用類型，則需要：
	if m.queryFilter != nil {
		cloned.queryFilter = make(bson.D, len(m.queryFilter))
		copy(cloned.queryFilter, m.queryFilter)
	}
	// 對其他可能的引用類型字段進行類似操作...
	return &cloned
}

func (m *mongoDBWrapper) Connect(uri string) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDBWrapper) Insert(ctx context.Context, collection string, document interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDBWrapper) Find(ctx context.Context, collection string, query interface{}) ([]bson.M, error) {
	// 設置集合
	m.collection = m.client.Database(m.database).Collection(collection)

	// 轉換自定義查詢條件為BSON
	var bsonQuery bson.M
	if query != nil {
		bsonBytes, err := bson.Marshal(query)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(bsonBytes, &bsonQuery)
		if err != nil {
			return nil, err
		}
	}

	// 構建查詢過濾條件
	var finalFilter bson.D
	if len(m.queryFilter) > 0 {
		finalFilter = append(finalFilter, bson.E{Key: "$and", Value: bson.A{m.queryFilter, bsonQuery}})
	} else {
		finalFilter = m.queryFilter
	}

	// 設置查詢選項
	findOptions := options.Find()
	if m.limit > 0 {
		findOptions.SetLimit(m.limit)
	}
	if m.offset > 0 {
		findOptions.SetSkip(m.offset)
	}
	if len(m.sort) > 0 {
		findOptions.SetSort(m.sort)
	}

	// 執行查詢
	cursor, err := m.collection.Find(ctx, finalFilter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 解析結果
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mongoDBWrapper) Count(ctx context.Context, collection string) (int64, error) {
	m.collection = m.client.Database(m.database).Collection(collection)
	count, err := m.collection.CountDocuments(ctx, m.queryFilter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *mongoDBWrapper) Update(ctx context.Context, collection string, query interface{}, update interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDBWrapper) UpdateBatch(ctx context.Context, collection string, updates []mongo.WriteModel) error {
	m.collection = m.client.Database(m.database).Collection(collection)
	_, err := m.collection.BulkWrite(ctx, updates)
	return err
}

func (m *mongoDBWrapper) Delete(ctx context.Context, collection string, query interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoDBWrapper) DeleteBatch(ctx context.Context, collection string, deletes []mongo.WriteModel) error {
	m.collection = m.client.Database(m.database).Collection(collection)
	_, err := m.collection.BulkWrite(ctx, deletes)
	return err
}

func (m *mongoDBWrapper) Where(field string, comparison string, value interface{}) MongoDBWrapper {
	var filter bson.E

	switch comparison {
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
		// 更多條件...
	}

	m.queryFilter = append(m.queryFilter, filter)
	return m
}

func (m *mongoDBWrapper) Sort(field string, order string) MongoDBWrapper {
	sortOrder := 1 // Ascending
	if order == "desc" {
		sortOrder = -1 // Descending
	}

	m.sort = append(m.sort, bson.E{Key: field, Value: sortOrder})
	return m
}

func (m *mongoDBWrapper) Limit(limit int64) MongoDBWrapper {
	m.limit = limit
	return m
}

func (m *mongoDBWrapper) Offset(offset int64) MongoDBWrapper {
	m.offset = offset
	return m
}

func (m *mongoDBWrapper) GroupBy(field string) MongoDBWrapper {
	m.groupBy = bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$" + field}}}}
	return m
}

func (m *mongoDBWrapper) Having(condition bson.M) MongoDBWrapper {
	m.groupBy = append(m.groupBy, bson.E{Key: "$match", Value: condition})
	return m
}
