package tag

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/nasu/nasulog/api/infrastructure/dynamodb"
)

type Repository struct {
	tableName    string
	partitionKey string
	ctx          context.Context
	db           *dynamodb.DB
}

func NewRepositoryWithContextAndDB(ctx context.Context, db *dynamodb.DB) *Repository {
	return &Repository{
		tableName:    "blog",
		partitionKey: "tag",
		ctx:          ctx,
		db:           db,
	}
}

// SelectAll gets all tags.
func (r *Repository) SelectAll() ([]*Tag, error) {
	items, err := r.db.SelectAll(r.ctx, r.tableName, r.partitionKey)
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
func (r *Repository) SelectByName(name string) (*Tag, error) {
	tags, err := r.SelectByNames([]string{name})
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, nil
	}
	return tags[0], nil
}

// SelectByNames gets tags with names.
func (r *Repository) SelectByNames(names []string) ([]*Tag, error) {
	items, err := r.db.SelectBySortKeys(r.ctx, r.tableName, r.partitionKey, names)
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
func (r *Repository) InsertMulti(tags []*Tag) error {
	items := make([]map[string]types.AttributeValue, len(tags), len(tags))
	for i, t := range tags {
		items[i] = t.ToAttributeValue()
	}
	return r.db.UpsertMulti(r.ctx, r.tableName, items)
}

// UpsertMulti upserts tags.
func (r *Repository) UpsertMulti(tags []*Tag) error {
	items := make([]map[string]types.AttributeValue, len(tags), len(tags))
	for i, t := range tags {
		items[i] = t.ToAttributeValue()
	}
	return r.db.UpsertMulti(r.ctx, r.tableName, items)
}

// DeleteByName deletes a tag with name.
func (r *Repository) DeleteByName(name string) error {
	key := map[string]types.AttributeValue{
		"partition_key": &types.AttributeValueMemberS{Value: r.partitionKey},
		"sort_key":      &types.AttributeValueMemberS{Value: name},
	}
	return r.db.DeleteByPK(r.ctx, r.tableName, key)
}

// DeleteMulti deletes a tag with name.
func (r *Repository) DeleteMulti(tags []*Tag) error {
	for _, t := range tags {
		if err := r.DeleteByName(t.Name); err != nil {
			return nil
		}
	}
	return nil
}
