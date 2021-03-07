package article

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
//TODO: write test
func NewArticleWithAttributeValue(values map[string]types.AttributeValue) *Article {
	article := &Article{}
	var err error
	if v, ok := values["id"].(*types.AttributeValueMemberS); ok {
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
		article.CreatedAt, err = time.Parse(time.RFC3339, v.Value)
		if err != nil {
			log.Printf("Failed to convert created_at, %v", err)
		}
	}
	if v, ok := values["updated_at"].(*types.AttributeValueMemberS); ok {
		article.UpdatedAt, err = time.Parse(time.RFC3339, v.Value)
		if err != nil {
			log.Printf("Failed to convert updated_at, %v", err)
		}
	}
	return article
}
