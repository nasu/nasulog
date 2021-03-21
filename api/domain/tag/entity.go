package tag

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

// Tag is an entity of an tag.
type Tag struct {
	Name     string
	Articles []string
}

// NewTagWithAttributeValue creates an Tag with dyannamodb.types.AttributeValue
func NewTagWithAttributeValue(values map[string]types.AttributeValue) *Tag {
	if len(values) == 0 {
		return nil
	}

	tag := &Tag{}
	if v, ok := values["sort_key"].(*types.AttributeValueMemberS); ok {
		tag.Name = v.Value
	}
	if v, ok := values["articles"].(*types.AttributeValueMemberSS); ok {
		tag.Articles = v.Value
	}
	return tag
}

// ToAttributeValue converts to attribute value.
func (tag Tag) ToAttributeValue() map[string]types.AttributeValue {
	item := make(map[string]types.AttributeValue)
	item["partition_key"] = &types.AttributeValueMemberS{Value: "tag"}
	item["sort_key"] = &types.AttributeValueMemberS{Value: tag.Name}
	if len(tag.Articles) != 0 {
		item["articles"] = &types.AttributeValueMemberSS{Value: tag.Articles}
	}
	return item
}
