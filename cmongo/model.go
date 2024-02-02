package cmongo

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// QueryCondition 定义一个查询条件的结构
type QueryCondition struct {
	Field    string // 字段名称
	Operator string // 操作符，如 "=", ">", "<", "like" 等
	Value    []any  // 值
}

type MongoDB struct {
	client   *mongo.Client
	database string
}

func NewMongoDB() Client {
	return &MongoDB{}
}

func (m *MongoDB) Connect(ctx context.Context, database string, opts ...*options.ClientOptions) (*MongoDB, error) {
	client, err := mongo.Connect(ctx, opts...)
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

func (m *MongoDB) InsertBatch(ctx context.Context, collection string, documents []any) error {
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
	for i, del := range deletes {
		filter, _, _ := del.Build()
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
			{{Key: "$match", Value: filter}},
			{{Key: "$group", Value: group}},
		}
		if len(qb.having) > 0 {
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: qb.having}})
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
			{{Key: "$match", Value: filter}},
			{{Key: "$group", Value: group}},
		}
		if len(qb.having) > 0 {
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: qb.having}})
		}
		pipeline = append(pipeline, bson.D{{Key: "$count", Value: "count"}})

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
	err       error
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
		qb.group = append(qb.group, bson.E{Key: "total_" + field, Value: bson.D{{Key: "$sum", Value: "$" + field}}})
	}
	return qb
}

func (qb *QueryBuilder) Having(query interface{}, args ...interface{}) *QueryBuilder {
	// 檢查 query 是否為字符串
	if condition, ok := query.(string); ok {
		if strings.Contains(condition, " ") {
			// 處理簡單的字符串查詢
			v := qb.processSimpleCondition(condition, true, args...)
			qb.having = append(qb.having, v...)
		} else if len(args) > 0 {
			qb.having = append(qb.having, bson.E{Key: "total_" + condition, Value: args[0]})
		}
	} else {
		switch v := query.(type) {
		case bson.E:
			// 直接添加 BSON 表達式
			qb.having = append(qb.having, v)
		case map[string]interface{}:
			// 處理映射類型的查詢
			for key, value := range v {
				qb.having = append(qb.having, bson.E{Key: key, Value: value})
			}
		case []bson.E:
			// 處理 BSON 表達式陣列
			qb.having = append(qb.having, v...)
		default:
			// 处理结构体类型的查询
			bsonBytes, err := bson.Marshal(v)
			if err != nil {
				qb.err = err
				return qb
			}
			var bsonDoc bson.D
			err = bson.Unmarshal(bsonBytes, &bsonDoc)
			if err != nil {
				qb.err = err
				return qb
			}
			qb.having = append(qb.having, bsonDoc...)
		}
	}
	return qb
}

func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	// 檢查 query 是否為字符串
	if condition, ok := query.(string); ok {
		if strings.Contains(condition, " ") {
			// 處理簡單的字符串查詢
			v := qb.processSimpleCondition(condition, false, args...)
			qb.filter = append(qb.filter, v...)
		} else if len(args) > 0 {
			qb.filter = append(qb.filter, bson.E{Key: condition, Value: args[0]})
		}
	} else {
		switch v := query.(type) {
		case bson.E:
			// 直接添加 BSON 表達式
			qb.filter = append(qb.filter, v)
		case map[string]interface{}:
			// 處理映射類型的查詢
			for key, value := range v {
				qb.filter = append(qb.filter, bson.E{Key: key, Value: value})
			}
		case []bson.E:
			// 處理 BSON 表達式陣列
			qb.filter = append(qb.filter, v...)
			// 可以根據需要添加更多的條件類型
		default:
			// 处理结构体类型的查询
			bsonBytes, err := bson.Marshal(v)
			if err != nil {
				qb.err = err
				return qb
			}
			var bsonDoc bson.D
			err = bson.Unmarshal(bsonBytes, &bsonDoc)
			if err != nil {
				qb.err = err
				return qb
			}
			qb.filter = append(qb.filter, bsonDoc...)
		}
	}
	return qb
}

func (qb *QueryBuilder) processSimpleCondition(condition string, isHaving bool, args ...interface{}) bson.D {
	// 使用正则表达式分割条件，以支持复杂条件，如 "AND"、"OR"
	conditionParts := regexp.MustCompile(`\s+(AND|OR|and|or)\s+`).Split(condition, -1)
	logicalOperators := regexp.MustCompile(`\s+(AND|OR|and|or)\s+`).FindAllString(condition, -1)

	argIndex := 0
	filters := make([]bson.E, 0)
	for _, part := range conditionParts {
		parts := strings.Fields(part)
		if len(parts) >= 3 {
			field := parts[0]
			operator := parts[1]
			var arrValue []any

			if isHaving == true {
				field = "total_" + field
			}

			// 检查是否有 '?' 占位符，若有则替换为 args 中的相应值
			for strings.Contains(parts[2], "?") && argIndex < len(args) {
				strings.Replace(parts[2], "?", "", 1)
				arrValue = append(arrValue, args[argIndex])
				argIndex++
			}
			if len(arrValue) == 0 {
				arrValue = append(arrValue, parts[2])
			}

			// 创建查询条件
			cond := QueryCondition{Field: field, Operator: operator, Value: arrValue}
			// 处理每个子条件
			filter := processConditionPart(cond)
			filters = append(filters, filter)
		}
	}

	// 在所有条件都被处理完之后，根据逻辑运算符来组合它们
	var result bson.D
	for i, filter := range filters {
		if i != 0 && i-1 < len(logicalOperators) {
			result = append(result, filter)
			op := parseLogic(logicalOperators[i-1])
			result = bson.D{bson.E{Key: op, Value: result}}
		} else {
			result = append(result, filter)
		}
	}
	return result
}

// processConditionPart 處理單個條件部分
func processConditionPart(condition QueryCondition) bson.E {
	// 生成过滤器
	var filter bson.E
	switch condition.Operator {
	case "=", ">", "<", ">=", "<=":
		if len(condition.Value) > 0 {
			filter = bson.E{Key: condition.Field, Value: bson.M{parseOperator(condition.Operator): condition.Value[0]}}
		}
	case "in":
		filter = bson.E{Key: condition.Field, Value: bson.M{parseOperator(condition.Operator): condition.Value}}
	case "like":
		// 假设 Value 是一个字符串
		if str, ok := condition.Value[0].(string); ok {
			filter = buildLikeFilter(condition.Field, str)
		}
	}

	return filter
}

// parseOperator 将常规比较操作符转换为 MongoDB 的操作符
func parseOperator(operator string) string {
	switch operator {
	case "=":
		return "$eq" // 等于
	case ">":
		return "$gt" // 大于
	case "<":
		return "$lt" // 小于
	case ">=":
		return "$gte" // 大于或等于
	case "<=":
		return "$lte" // 小于或等于
	case "in":
		return "$in"
	default:
		return "" // 如果操作符不匹配，返回空字符串
	}
}

func parseLogic(logic string) string {
	logic = strings.ToLower(strings.TrimSpace(logic))
	return "$" + logic
}

// buildLikeFilter 創建用於模糊匹配的過濾器
func buildLikeFilter(field, pattern string) bson.E {
	regexPattern := ""
	if strings.HasPrefix(pattern, "%") && strings.HasSuffix(pattern, "%") {
		regexPattern = strings.Trim(pattern, "%")
	} else if strings.HasPrefix(pattern, "%") {
		regexPattern = strings.TrimLeft(pattern, "%") + "$"
	} else if strings.HasSuffix(pattern, "%") {
		regexPattern = "^" + strings.TrimRight(pattern, "%")
	}
	regex := bson.M{"$regex": regexPattern, "$options": "i"} // 不區分大小寫
	return bson.E{Key: field, Value: regex}
}
