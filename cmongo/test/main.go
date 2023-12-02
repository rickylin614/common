package test

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryBuilder struct {
	filter bson.D
	sort   bson.D
	limit  int64
	skip   int64
	group  bson.D
	having bson.D
}

func (qb *QueryBuilder) Group(key string, value interface{}) *QueryBuilder {
	qb.group = append(qb.group, bson.E{Key: key, Value: value})
	return qb
}

func (qb *QueryBuilder) Having(key string, value interface{}) *QueryBuilder {
	qb.having = append(qb.having, bson.E{Key: key, Value: value})
	return qb
}

func (qb *QueryBuilder) Build() (bson.D, bson.D, *options.FindOptions) {
	findOptions := options.Find()
	if qb.sort != nil {
		findOptions.SetSort(qb.sort)
	}
	if qb.limit > 0 {
		findOptions.SetLimit(qb.limit)
	}
	if qb.skip > 0 {
		findOptions.SetSkip(qb.skip)
	}
	return qb.filter, qb.group, findOptions
}

type MongoDB struct {
	client   *mongo.Client
	database string
}

func (m *MongoDB) Find(ctx context.Context, cc string, qb *QueryBuilder, results interface{}) error {
	filter, group, opts := qb.Build()
	collection := m.client.Database(m.database).Collection(cc)
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
