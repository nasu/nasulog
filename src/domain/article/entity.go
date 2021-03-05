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
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WithDynamoDbItems create an Article with dyannamodb.types.AttributeValue
//TODO: write test
func (e *Article) WithDynamoDbItems(values map[string]types.AttributeValue) *Article {
	var err error
	if v, ok := values["id"].(*types.AttributeValueMemberS); ok {
		e.ID = v.Value
	}
	if v, ok := values["title"].(*types.AttributeValueMemberS); ok {
		e.Title = v.Value
	}
	if v, ok := values["content"].(*types.AttributeValueMemberS); ok {
		e.Content = v.Value
	}
	if v, ok := values["created_at"].(*types.AttributeValueMemberS); ok {
		e.CreatedAt, err = time.Parse(time.RFC3339, v.Value)
		if err != nil {
			log.Printf("Failed to convert created_at, %v", err)
		}
	}
	if v, ok := values["updated_at"].(*types.AttributeValueMemberS); ok {
		e.UpdatedAt, err = time.Parse(time.RFC3339, v.Value)
		if err != nil {
			log.Printf("Failed to convert updated_at, %v", err)
		}
	}
	return e
}
