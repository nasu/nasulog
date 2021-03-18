package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nasu/nasulog/domain/article"
	"github.com/nasu/nasulog/domain/tag"
	"github.com/nasu/nasulog/graph/generated"
	"github.com/nasu/nasulog/graph/model"
)

func (r *mutationResolver) CreateArticle(ctx context.Context, input model.NewArticle) (*model.Article, error) {
	//TODO: transaction
	entity, err := article.Insert(ctx, r.DB, &article.Article{
		ID:      uuid.NewString(),
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	})
	if err != nil {
		return nil, err
	}
	log.Println("insert")

	if err := tag.InsertMulti(ctx, r.DB, input.Tags); err != nil {
		return nil, err
	}

	return &model.Article{
		ID:        entity.ID,
		Title:     entity.Title,
		Content:   entity.Content,
		Tags:      entity.Tags,
		CreatedAt: entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt: entity.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (r *mutationResolver) DeleteArticle(ctx context.Context, id string) (*bool, error) {
	err := article.DeleteByPK(ctx, r.DB, id)
	res := true
	return &res, err
}

func (r *mutationResolver) DeleteTag(ctx context.Context, name string) (*bool, error) {
	err := tag.DeleteByPK(ctx, r.DB, name)
	res := true
	return &res, err
}

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	entities, err := article.SelectAll(ctx, r.DB)
	if err != nil {
		return nil, err
	}
	models := make([]*model.Article, len(entities), len(entities))
	for i := 0; i < len(entities); i++ {
		entity := entities[i]
		models[i] = &model.Article{
			ID:        entity.ID,
			Title:     entity.Title,
			Content:   entity.Content,
			Tags:      entity.Tags,
			CreatedAt: entity.CreatedAt.Format(time.RFC3339),
			UpdatedAt: entity.UpdatedAt.Format(time.RFC3339),
		}
	}
	return models, nil
}

func (r *queryResolver) Tags(ctx context.Context) ([]string, error) {
	return tag.SelectAll(ctx, r.DB)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
