package article

//TODO: write test

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Article is an entity of an article.
type Article struct {
	ID        string
	Title     string
	Content   string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewArticleWithAttributeValue creates an Article with dyannamodb.types.AttributeValue
func NewArticleWithAttributeValue(values map[string]types.AttributeValue) *Article {
	if len(values) == 0 {
		return nil
	}

	article := &Article{}
	if v, ok := values["sort_key"].(*types.AttributeValueMemberS); ok {
		article.ID = v.Value
	}
	if v, ok := values["title"].(*types.AttributeValueMemberS); ok {
		article.Title = v.Value
	}
	if v, ok := values["content"].(*types.AttributeValueMemberS); ok {
		article.Content = v.Value
	}
	if v, ok := values["tags"].(*types.AttributeValueMemberSS); ok {
		article.Tags = v.Value
	}
	if v, ok := values["created_at"].(*types.AttributeValueMemberS); ok {
		var err error
		article.CreatedAt, err = time.Parse(time.RFC3339, v.Value)
		if err != nil {
			log.Printf("Failed to convert created_at, %v", err)
		}
	}
	if v, ok := values["updated_at"].(*types.AttributeValueMemberS); ok {
		var err error
		article.UpdatedAt, err = time.Parse(time.RFC3339, v.Value)
		if err != nil {
			log.Printf("Failed to convert updated_at, %v", err)
		}
	}
	return article
}

// ToAttributeValue converts to attribute value.
func (article Article) ToAttributeValue() map[string]types.AttributeValue {
	item := make(map[string]types.AttributeValue)
	item["partition_key"] = &types.AttributeValueMemberS{Value: "article"}
	item["sort_key"] = &types.AttributeValueMemberS{Value: article.ID}
	item["title"] = &types.AttributeValueMemberS{Value: article.Title}
	item["content"] = &types.AttributeValueMemberS{Value: article.Content}
	if len(article.Tags) != 0 {
		item["tags"] = &types.AttributeValueMemberSS{Value: article.Tags}
	}
	item["created_at"] = &types.AttributeValueMemberS{Value: article.CreatedAt.Format(time.RFC3339)}
	item["updated_at"] = &types.AttributeValueMemberS{Value: article.UpdatedAt.Format(time.RFC3339)}
	return item
}
