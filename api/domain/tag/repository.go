package tag

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/nasu/nasulog/infrastructure/dynamodb"
)

var tableName = "blog"
var partitionKey = "tag"

// SelectAll gets all tags.
func SelectAll(ctx context.Context, db *dynamodb.DB) ([]*Tag, error) {
	items, err := db.SelectAll(ctx, tableName, partitionKey)
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, len(items), len(items))
	for i, item := range items {
		tags[i] = NewTagWithAttributeValue(item)
	}
	return tags, nil
}

// SelectByName gets tags with name.
func SelectByName(ctx context.Context, db *dynamodb.DB, name string) (*Tag, error) {
	tags, err := SelectByNames(ctx, db, []string{name})
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, nil
	}
	return tags[0], nil
}

// SelectByNames gets tags with names.
func SelectByNames(ctx context.Context, db *dynamodb.DB, names []string) ([]*Tag, error) {
	items, err := db.SelectBySortKeys(ctx, tableName, partitionKey, names)
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

// UpsertMulti upserts tags.
func UpsertMulti(ctx context.Context, db *dynamodb.DB, tags []*Tag) error {
	items := make([]map[string]types.AttributeValue, len(tags), len(tags))
	for i, t := range tags {
		items[i] = t.ToAttributeValue()
	}
	return db.UpsertMulti(ctx, tableName, items)
}

// DeleteByName deletes an article with name.
func DeleteByName(ctx context.Context, db *dynamodb.DB, name string) error {
	key := map[string]types.AttributeValue{
		"partition_key": &types.AttributeValueMemberS{Value: partitionKey},
		"sort_key":      &types.AttributeValueMemberS{Value: name},
	}
	return db.DeleteByPK(ctx, tableName, key)
}
