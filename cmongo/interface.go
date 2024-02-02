package cmongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	Connect(ctx context.Context, database string, opts ...*options.ClientOptions) (*MongoDB, error)

	Insert(ctx context.Context, collection string, document interface{}) error
	InsertBatch(ctx context.Context, collection string, documents []any) error
	Update(ctx context.Context, collection string, qb *QueryBuilder, update interface{}) error
	UpdateBatch(ctx context.Context, collection string, updates []UpdateModel) error
	Delete(ctx context.Context, collection string, qb *QueryBuilder) error
	DeleteBatch(ctx context.Context, collection string, deletes []*QueryBuilder) error

	Find(ctx context.Context, table string, qb *QueryBuilder, results any) error
	Count(ctx context.Context, table string, qb *QueryBuilder) (int64, error)
}
