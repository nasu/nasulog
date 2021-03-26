package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nasu/nasulog/api/graph/generated"
	"github.com/nasu/nasulog/api/graph/model"
	"github.com/nasu/nasulog/api/graph/presenter"
)

func (r *mutationResolver) PostArticle(ctx context.Context, input model.PostArticle) (*model.Article, error) {
	return presenter.PostArticle(ctx, r.DB, &input)
}

func (r *mutationResolver) DeleteArticle(ctx context.Context, id string) (*bool, error) {
	return presenter.DeleteArticle(ctx, r.DB, id)
}

func (r *mutationResolver) DeleteTag(ctx context.Context, name string) (*bool, error) {
	return presenter.DeleteTag(ctx, r.DB, name)
}

func (r *queryResolver) Articles(ctx context.Context, cond *model.ArticleCondition) ([]*model.Article, error) {
	return presenter.GetArticles(ctx, r.DB, cond)
}

func (r *queryResolver) Article(ctx context.Context, id string) (*model.Article, error) {
	return presenter.GetArticle(ctx, r.DB, id)
}

func (r *queryResolver) Tags(ctx context.Context) ([]*model.Tag, error) {
	return presenter.GetTags(ctx, r.DB)
}

func (r *queryResolver) Tag(ctx context.Context, name string) (*model.Tag, error) {
	return presenter.GetTag(ctx, r.DB, name)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
