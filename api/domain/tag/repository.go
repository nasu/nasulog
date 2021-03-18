package tag

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/nasu/nasulog/infrastructure/dynamodb"
)

var tableName = "tags"

// SelectAll gets all tags.
func SelectAll(ctx context.Context, db *dynamodb.DB) ([]string, error) {
	items, err := db.SelectAll(ctx, tableName)
	if err != nil {
		return nil, err
	}

	tags := make([]string, len(items), len(items))
	for i, item := range items {
		if v, ok := item["name"].(*types.AttributeValueMemberS); ok {
			tags[i] = v.Value
		}
	}
	return tags, nil
}

// InsertMulti inserts all tags received.
func InsertMulti(ctx context.Context, db *dynamodb.DB, names []string) error {
	items := make([]map[string]types.AttributeValue, len(names), len(names))
	for i, name := range names {
		item := make(map[string]types.AttributeValue)
		item["name"] = &types.AttributeValueMemberS{Value: name}

		items[i] = item
	}
	return db.UpsertMulti(ctx, tableName, items)
}

// DeleteByPK deletes an article with name.
func DeleteByPK(ctx context.Context, db *dynamodb.DB, name string) error {
	key := map[string]types.AttributeValue{
		"name": &types.AttributeValueMemberS{Value: name},
	}
	return db.DeleteByPK(ctx, tableName, key)
}
