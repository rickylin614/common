package cmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoDBWrapper interface {
	Connect(uri string) error

	// CRUD operations
	Insert(ctx context.Context, collection string, document interface{}) error
	Find(ctx context.Context, collection string, query interface{}) ([]bson.M, error)
	Update(ctx context.Context, collection string, query interface{}, update interface{}) error
	Delete(ctx context.Context, collection string, query interface{}) error
	// UpdateBatch(collection string, updates []mongo.WriteModel)
	// DeleteBatch(collection string, deletes []mongo.WriteModel) error

	// Fluent interface methods
	Where(field string, comparison string, value interface{}) MongoDBWrapper
	Sort(field string, order string) MongoDBWrapper
	Limit(limit int64) MongoDBWrapper
	Offset(offset int64) MongoDBWrapper
	GroupBy(field string) MongoDBWrapper
	Having(condition bson.M) MongoDBWrapper
}
