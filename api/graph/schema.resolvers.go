package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
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

	tags, err := tag.SelectByNames(ctx, r.DB, input.Tags)
	if err != nil {
		return nil, err
	}
	tagsMap := make(map[string]*tag.Tag)
	for _, t := range tags {
		tagsMap[t.Name] = t
	}
	for _, it := range input.Tags {
		if et, ok := tagsMap[it]; ok {
			et.Articles = append(et.Articles, entity.ID)
		} else {
			tags = append(tags, &tag.Tag{
				Name:     it,
				Articles: []string{entity.ID},
			})
		}
	}

	if err := tag.InsertMulti(ctx, r.DB, tags); err != nil {
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
	err := article.DeleteByID(ctx, r.DB, id)
	res := true
	return &res, err
}

func (r *mutationResolver) DeleteTag(ctx context.Context, name string) (*bool, error) {
	err := tag.DeleteByName(ctx, r.DB, name)
	res := true
	return &res, err
}

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	entities, err := article.SelectAll(ctx, r.DB)
	if err != nil {
		return nil, err
	}
	models := make([]*model.Article, len(entities), len(entities))
	for i, entity := range entities {
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

func (r *queryResolver) Tags(ctx context.Context) ([]*model.Tag, error) {
	entities, err := tag.SelectAll(ctx, r.DB)
	if err != nil {
		return nil, err
	}
	models := make([]*model.Tag, len(entities), len(entities))
	for i, entity := range entities {
		models[i] = &model.Tag{
			Name:     entity.Name,
			Articles: entity.Articles,
		}
	}
	return models, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
