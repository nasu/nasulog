package article

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	mydb "github.com/nasu/nasulog/infrastructure/dynamodb"
)

var tableName = "articles"
var consistentRead = true
var scanLimit = int32(10)

// SelectByPK gets an article with ID.
func SelectByPK(ctx context.Context, db *mydb.DB, id string) (*Article, error) {
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}
	consistentRead := true
	res, err := db.SelectByPK(ctx, tableName, key, consistentRead)
	if err != nil {
		return nil, err
	}
	return NewArticleWithAttributeValue(res), nil
}

// SelectAll gets all articles.
func SelectAll(ctx context.Context, db *mydb.DB) ([]*Article, error) {
	items, err := db.SelectAll(ctx, tableName)
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

// DeleteByPK deletes an article with id.
func DeleteByPK(ctx context.Context, db *mydb.DB, id string) error {
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}
	return db.DeleteByPK(ctx, tableName, key)
}
