package tag

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/nasu/nasulog/infrastructure/dynamodb"
)

var tableName = "blog"

// SelectAll gets all tags.
func SelectAll(ctx context.Context, db *dynamodb.DB) ([]*Tag, error) {
	items, err := db.SelectAll(ctx, tableName, "tag")
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, len(items), len(items))
	for i, item := range items {
		tags[i] = NewTagWithAttributeValue(item)
	}
	return tags, nil
}

// SelectByNames gets tags with names.
func SelectByNames(ctx context.Context, db *dynamodb.DB, names []string) ([]*Tag, error) {
	items, err := db.SelectBySortKeys(ctx, tableName, "tag", names)
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, len(items), len(items))
	for i, item := range items {
		tags[i] = NewTagWithAttributeValue(item)
	}
	return tags, nil
}

// InsertMulti inserts all tags received.
func InsertMulti(ctx context.Context, db *dynamodb.DB, tags []*Tag) error {
	items := make([]map[string]types.AttributeValue, len(tags), len(tags))
	for i, t := range tags {
		items[i] = t.ToAttributeValue()
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
