package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nasu/nasulog/api/domain/article"
	"github.com/nasu/nasulog/api/domain/tag"
	"github.com/nasu/nasulog/api/graph/generated"
	"github.com/nasu/nasulog/api/graph/model"
)

func (r *mutationResolver) CreateArticle(ctx context.Context, input model.NewArticle) (*model.Article, error) {
	repoArticle := article.NewRepositoryWithContextAndDB(ctx, r.DB)
	repoTag := tag.NewRepositoryWithContextAndDB(ctx, r.DB)

	//TODO: transaction
	entity, err := repoArticle.Insert(&article.Article{
		ID:      uuid.NewString(),
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	})
	if err != nil {
		return nil, err
	}

	tags, err := repoTag.SelectByNames(input.Tags)
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

	if err := repoTag.InsertMulti(tags); err != nil {
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
	repoArticle := article.NewRepositoryWithContextAndDB(ctx, r.DB)
	repoTag := tag.NewRepositoryWithContextAndDB(ctx, r.DB)

	entity, err := repoArticle.SelectByID(id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		res := false
		return &res, nil
	}

	// Remove an article matched to the article from tag.articles
	tags, err := repoTag.SelectByNames(entity.Tags)
	if err != nil {
		return nil, err
	}
	tagsShouldDelete := make([]*tag.Tag, 0)
	tagsShouldUpdate := make([]*tag.Tag, 0)
	for _, t := range tags {
		for i, a := range t.Articles {
			if a == id {
				t.Articles = append(t.Articles[:i], t.Articles[i+1:]...)
				break
			}
		}
		if len(t.Articles) == 0 {
			tagsShouldDelete = append(tagsShouldDelete, t)
		} else {
			tagsShouldUpdate = append(tagsShouldUpdate, t)
		}
	}

	if err := repoArticle.DeleteByID(id); err != nil {
		return nil, err
	}
	if err := repoTag.UpsertMulti(tagsShouldUpdate); err != nil {
		return nil, err
	}
	if err := repoTag.DeleteMulti(tagsShouldDelete); err != nil {
		return nil, err
	}

	res := true
	return &res, err
}

func (r *mutationResolver) DeleteTag(ctx context.Context, name string) (*bool, error) {
	repoArticle := article.NewRepositoryWithContextAndDB(ctx, r.DB)
	repoTag := tag.NewRepositoryWithContextAndDB(ctx, r.DB)

	entity, err := repoTag.SelectByName(name)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		res := false
		return &res, nil
	}

	// Remove a tag matched to the tag from article.tags
	articles, err := repoArticle.SelectByIDs(entity.Articles)
	if err != nil {
		return nil, err
	}
	for _, a := range articles {
		for i, t := range a.Tags {
			if t == name {
				a.Tags = append(a.Tags[:i], a.Tags[i+1:]...)
				break
			}
		}
	}

	// transaction
	if err := repoTag.DeleteByName(name); err != nil {
		return nil, err
	}
	if err := repoArticle.UpsertMulti(articles); err != nil {
		return nil, err
	}

	res := true
	return &res, err
}

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	repo := article.NewRepositoryWithContextAndDB(ctx, r.DB)
	entities, err := repo.SelectAll()
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
	repo := tag.NewRepositoryWithContextAndDB(ctx, r.DB)
	entities, err := repo.SelectAll()
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
