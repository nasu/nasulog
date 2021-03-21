package article

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	mydb "github.com/nasu/nasulog/infrastructure/dynamodb"
)

var tableName = "blog"
var partitionKey = "article"
var consistentRead = true
var scanLimit = int32(10)

// SelectByID gets an article with ID.
func SelectByID(ctx context.Context, db *mydb.DB, id string) (*Article, error) {
	items, err := SelectByIDs(ctx, db, []string{id})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return items[0], nil
}

// SelectByIDs gets articles with IDs.
func SelectByIDs(ctx context.Context, db *mydb.DB, ids []string) ([]*Article, error) {
	items, err := db.SelectBySortKeys(ctx, tableName, partitionKey, ids)
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, len(items), len(items))
	for i, item := range items {
		articles[i] = NewArticleWithAttributeValue(item)
	}
	return articles, nil
}

// SelectAll gets all articles.
func SelectAll(ctx context.Context, db *mydb.DB) ([]*Article, error) {
	items, err := db.SelectAll(ctx, tableName, partitionKey)
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, len(items), len(items))
	for i, item := range items {
		articles[i] = NewArticleWithAttributeValue(item)
	}
	return articles, nil
}

// Insert inserts an article to DB.
func Insert(ctx context.Context, db *mydb.DB, article *Article) (*Article, error) {
	now := time.Now()
	if article.CreatedAt.IsZero() {
		article.CreatedAt = now
	}
	if article.UpdatedAt.IsZero() {
		article.UpdatedAt = now
	}

	item := article.ToAttributeValue()
	err := db.Insert(ctx, tableName, item)
	if err != nil {
		return nil, err
	}
	return NewArticleWithAttributeValue(item), nil
}

// UpsertMulti upserts articles.
// This method doesn't automatically update as it's possible to update only tag.
func UpsertMulti(ctx context.Context, db *mydb.DB, articles []*Article) error {
	items := make([]map[string]types.AttributeValue, len(articles), len(articles))
	for i, a := range articles {
		items[i] = a.ToAttributeValue()
	}
	return db.UpsertMulti(ctx, tableName, items)
}

// DeleteByID deletes an article with id.
func DeleteByID(ctx context.Context, db *mydb.DB, id string) error {
	key := map[string]types.AttributeValue{
		"partition_key": &types.AttributeValueMemberS{Value: partitionKey},
		"sort_key":      &types.AttributeValueMemberS{Value: id},
	}
	return db.DeleteByPK(ctx, tableName, key)
}
