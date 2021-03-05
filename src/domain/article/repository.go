package article

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

var tableName = "articles"
var consistentRead = true
var scanLimit = int32(10)

func selectOne(ctx context.Context, client *dynamodb.Client, id string) (*Article, error) {
	consistentRead := true
	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}
	res, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      &tableName,
		Key:            key,
		ConsistentRead: &consistentRead,
	})
	if err != nil {
		return nil, err
	}
	return NewArticleWithAttributeValue(res.Item), nil
}

func SelectAll(ctx context.Context, client *dynamodb.Client) ([]*Article, error) {
	res, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName:      &tableName,
		ConsistentRead: &consistentRead,
		Limit:          &scanLimit,
	})
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, res.Count, res.Count)
	for i := 0; i < int(res.Count); i++ {
		articles[i] = NewArticleWithAttributeValue(res.Items[i])
	}
	return articles, nil
}

func Insert(ctx context.Context, client *dynamodb.Client, article *Article) (*Article, error) {
	now := time.Now()
	if article.CreatedAt.IsZero() {
		article.CreatedAt = now
	}
	if article.UpdatedAt.IsZero() {
		article.UpdatedAt = now
	}

	item := make(map[string]types.AttributeValue)
	item["id"] = &types.AttributeValueMemberS{Value: uuid.NewString()}
	item["title"] = &types.AttributeValueMemberS{Value: article.Title}
	item["content"] = &types.AttributeValueMemberS{Value: article.Content}
	item["created_at"] = &types.AttributeValueMemberS{Value: article.CreatedAt.Format(time.RFC3339)}
	item["updated_at"] = &types.AttributeValueMemberS{Value: article.UpdatedAt.Format(time.RFC3339)}

	conditionExpression := "attribute_not_exists(id)"
	params := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
		//TODO: write test checking a case of ConditionalCheckFailedException
		ConditionExpression: &conditionExpression,
	}
	// PutItem returns nothing as the first value if we don't specify ReturnValues.
	// So we ignore the value.
	_, err := client.PutItem(ctx, params)
	if err != nil {
		return nil, err
	}
	return NewArticleWithAttributeValue(item), nil
}
