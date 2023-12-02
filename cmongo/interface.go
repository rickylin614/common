package cmongo

// type MongoDBWrapper interface {
// 	Connect(ctx context.Context, uri string, database string) error

// 	// CRUD operations
// 	Insert(ctx context.Context, collection string, document interface{}) error
// 	Find(ctx context.Context, collection string, target any) error
// 	Update(ctx context.Context, collection string, query interface{}, update interface{}) error
// 	Delete(ctx context.Context, collection string, query interface{}) error
// 	// UpdateBatch(collection string, updates []mongo.WriteModel)
// 	// DeleteBatch(collection string, deletes []mongo.WriteModel) error

// 	// Fluent interface methods
// 	Where(field string, value interface{}, comparison ...string) MongoDBWrapper
// 	Sort(field string, order string) MongoDBWrapper
// 	Limit(limit int64) MongoDBWrapper
// 	Offset(offset int64) MongoDBWrapper
// 	GroupBy(field string) MongoDBWrapper
// 	Having(condition bson.M) MongoDBWrapper
// }
