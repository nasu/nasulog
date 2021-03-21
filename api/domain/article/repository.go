package article

import (
	"context"
	"time"

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
		partitionKey: "article",
		ctx:          ctx,
		db:           db,
	}
}

// SelectByID gets an article with ID.
func (r *Repository) SelectByID(id string) (*Article, error) {
	items, err := r.SelectByIDs([]string{id})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	return items[0], nil
}

// SelectByIDs gets articles with IDs.
func (r *Repository) SelectByIDs(ids []string) ([]*Article, error) {
	items, err := r.db.SelectBySortKeys(r.ctx, r.tableName, r.partitionKey, ids)
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
func (r *Repository) SelectAll() ([]*Article, error) {
	items, err := r.db.SelectAll(r.ctx, r.tableName, r.partitionKey)
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
func (r *Repository) Insert(article *Article) (*Article, error) {
	now := time.Now()
	if article.CreatedAt.IsZero() {
		article.CreatedAt = now
	}
	if article.UpdatedAt.IsZero() {
		article.UpdatedAt = now
	}

	item := article.ToAttributeValue()
	err := r.db.Insert(r.ctx, r.tableName, item)
	if err != nil {
		return nil, err
	}
	return NewArticleWithAttributeValue(item), nil
}

// UpsertMulti upserts articles.
// This method doesn't automatically update as it's possible to update only tag.
func (r *Repository) UpsertMulti(articles []*Article) error {
	items := make([]map[string]types.AttributeValue, len(articles), len(articles))
	for i, a := range articles {
		items[i] = a.ToAttributeValue()
	}
	return r.db.UpsertMulti(r.ctx, r.tableName, items)
}

// DeleteByID deletes an article with id.
func (r *Repository) DeleteByID(id string) error {
	key := map[string]types.AttributeValue{
		"partition_key": &types.AttributeValueMemberS{Value: r.partitionKey},
		"sort_key":      &types.AttributeValueMemberS{Value: id},
	}
	return r.db.DeleteByPK(r.ctx, r.tableName, key)
}
